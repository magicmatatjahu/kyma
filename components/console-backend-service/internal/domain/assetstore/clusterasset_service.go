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

type clusterAssetService struct {
	informer cache.SharedIndexInformer
	notifier notifier
}

func newClusterAssetService(informer cache.SharedIndexInformer) (*clusterAssetService, error) {
	svc := &clusterAssetService{
		informer: informer,
	}

	err := svc.informer.AddIndexers(cache.Indexers{
		"docsTopicName": func(obj interface{}) ([]string, error) {
			entity, err := svc.extractClusterAsset(obj)
			if err != nil {
				return nil, errors.New("Cannot convert item")
			}

			return []string{entity.Labels["docstopic.cms.kyma-project.io"]}, nil
		},
		"docsTopicName/type": func(obj interface{}) ([]string, error) {
			entity, err := svc.extractClusterAsset(obj)
			if err != nil {
				return nil, errors.New("Cannot convert item")
			}

			return []string{fmt.Sprintf("%s/%s", entity.Labels["docstopic.cms.kyma-project.io"], entity.Labels["type.cms.kyma-project.io"])}, nil
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

func (svc *clusterAssetService) Find(name string) (*v1alpha2.ClusterAsset, error) {
	item, exists, err := svc.informer.GetStore().GetByKey(name)
	if err != nil || !exists {
		return nil, err
	}

	clusterAsset, err := svc.extractClusterAsset(item)
	if err != nil {
		errors.Wrapf(err, "Incorrect item type: %T, should be: *ClusterAsset", item)
	}

	return clusterAsset, nil
}

func (svc *clusterAssetService) ListForDocsTopicByType(docsTopicName string, types []string) ([]*v1alpha2.ClusterAsset, error) {
	var items []interface{}
	var err error
	if len(types) == 0 {
		items, err = svc.informer.GetIndexer().ByIndex("docsTopicName", docsTopicName)
	} else {
		for _, typeArg := range types {
			its, err := svc.informer.GetIndexer().ByIndex("docsTopicName/type", fmt.Sprintf("%s/%s", docsTopicName, typeArg))
			if err != nil {
				return nil, err
			}
			items = append(items, its...)
		}
	}

	if err != nil {
		return nil, err
	}

	var clusterAssets []*v1alpha2.ClusterAsset
	for _, item := range items {
		clusterAsset, err := svc.extractClusterAsset(item)
		if err != nil {
			return nil, errors.Wrapf(err, "Incorrect item type: %T, should be: *ClusterAsset", item)
		}

		clusterAssets = append(clusterAssets, clusterAsset)
	}

	return clusterAssets, nil
}

func (svc *clusterAssetService) Subscribe(listener resource.Listener) {
	svc.notifier.AddListener(listener)
}

func (svc *clusterAssetService) Unsubscribe(listener resource.Listener) {
	svc.notifier.DeleteListener(listener)
}

func (svc *clusterAssetService) extractClusterAsset(obj interface{}) (*v1alpha2.ClusterAsset, error) {
	u, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, errors.Wrapf(err, "while converting resource %s %s to unstructured", pretty.ClusterAsset, obj)
	}

	var clusterAsset v1alpha2.ClusterAsset
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(u, &clusterAsset)
	if err != nil {
		return nil, errors.Wrapf(err, "while converting unstructured to resource %s %s", pretty.ClusterAsset, u)
	}

	return &clusterAsset, nil
}