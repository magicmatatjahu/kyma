package assetstore

import (
	"k8s.io/client-go/tools/cache"
	"fmt"
	"github.com/pkg/errors"
	"github.com/kyma-project/kyma/components/asset-store-controller-manager/pkg/apis/assetstore/v1alpha2"
	"k8s.io/apimachinery/pkg/runtime"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/assetstore/pretty"
	"github.com/kyma-project/kyma/components/console-backend-service/pkg/resource"
)

type assetService struct {
	informer cache.SharedIndexInformer
	notifier notifier
}

func newAssetService(informer cache.SharedIndexInformer) (*assetService, error) {
	svc := &assetService{
		informer: informer,
	}

	err := svc.informer.AddIndexers(cache.Indexers{
		"docsTopicName": func(obj interface{}) ([]string, error) {
			entity, err := svc.extractAsset(obj)
			if err != nil {
				return nil, errors.New("Cannot convert item")
			}

			return []string{fmt.Sprintf("%s/%s", entity.Namespace, entity.Labels["docstopic.cms.kyma-project.io"])}, nil
		},
		"docsTopicName/type": func(obj interface{}) ([]string, error) {
			entity, err := svc.extractAsset(obj)
			if err != nil {
				return nil, errors.New("Cannot convert item")
			}

			return []string{fmt.Sprintf("%s/%s/%s", entity.Namespace, entity.Labels["docstopic.cms.kyma-project.io"], entity.Labels["type.cms.kyma-project.io"])}, nil
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "while adding indexers")
	}

	notifier := resource.NewNotifier()
	informer.AddEventHandler(notifier)
	svc.notifier = notifier

	return svc, nil
}

func (svc *assetService) Find(namespace, name string) (*v1alpha2.Asset, error) {
	key := fmt.Sprintf("%s/%s", namespace, name)

	item, exists, err := svc.informer.GetStore().GetByKey(key)
	if err != nil || !exists {
		return nil, err
	}

	asset, err := svc.extractAsset(item)
	if err != nil {
		errors.Wrapf(err, "Incorrect item type: %T, should be: *Asset", item)
	}

	return asset, nil
}

func (svc *assetService) List(namespace, groupName string) ([]*v1alpha2.Asset, error) {
	key := fmt.Sprintf("%s", groupName)
	items, err := svc.informer.GetIndexer().ByIndex("groupName", key)
	if err != nil {
		return nil, err
	}

	var assets []*v1alpha2.Asset
	for _, item := range items {
		asset, err := svc.extractAsset(item)
		if err != nil {
			return nil, errors.Wrapf(err, "Incorrect item type: %T, should be: *Asset", item)
		}

		assets = append(assets, asset)
	}

	return assets, nil
}

func (svc *assetService) ListForDocsTopicByType(namespace, docsTopicName string, types []string) ([]*v1alpha2.Asset, error) {
	var items []interface{}
	var err error
	if len(types) == 0 {
		items, err = svc.informer.GetIndexer().ByIndex("docsTopicName/type", fmt.Sprintf("%s/%s", namespace, docsTopicName))
	} else {
		for _, typeArg := range types {
			its, err := svc.informer.GetIndexer().ByIndex("docsTopicName/type", fmt.Sprintf("%s/%s/%s", namespace, docsTopicName, typeArg))
			if err != nil {
				return nil, err
			}
			items = append(items, its...)
		}
	}

	if err != nil {
		return nil, err
	}

	var assets []*v1alpha2.Asset
	for _, item := range items {
		asset, err := svc.extractAsset(item)
		if err != nil {
			return nil, errors.Wrapf(err, "Incorrect item type: %T, should be: *Asset", item)
		}

		assets = append(assets, asset)
	}

	return assets, nil
}

func (svc *assetService) Subscribe(listener resource.Listener) {
	svc.notifier.AddListener(listener)
}

func (svc *assetService) Unsubscribe(listener resource.Listener) {
	svc.notifier.DeleteListener(listener)
}

func (svc *assetService) extractAsset(obj interface{}) (*v1alpha2.Asset, error) {
	u, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, errors.Wrapf(err, "while converting resource %s %s to unstructured", pretty.Asset, obj)
	}

	var asset v1alpha2.Asset
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(u, &asset)
	if err != nil {
		return nil, errors.Wrapf(err, "while converting unstructured to resource %s %s", pretty.Asset, u)
	}

	return &asset, nil
}