package assetstore

import (
	"k8s.io/client-go/tools/cache"
	"fmt"
	"github.com/pkg/errors"
	"github.com/kyma-project/kyma/components/asset-store-controller-manager/pkg/apis/assetstore/v1alpha2"
	"k8s.io/apimachinery/pkg/runtime"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/assetstore/pretty"
)

type assetService struct {
	informer cache.SharedIndexInformer
}

func newAssetService(informer cache.SharedIndexInformer) (*assetService, error) {
	svc := &assetService{
		informer: informer,
	}

	err := svc.informer.AddIndexers(cache.Indexers{
		"groupName": func(obj interface{}) ([]string, error) {
			entity, err := svc.extractAsset(obj)
			if err != nil {
				return nil, errors.New("Cannot convert item")
			}

			return []string{fmt.Sprintf("%s", entity.Spec.BucketRef)}, nil
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "while adding indexers")
	}

	return svc, nil
}

func (svc *assetService) List(namespace, groupName string) ([]*v1alpha2.Asset, error) {
	key := fmt.Sprintf("%s", groupName)
	items, err := svc.informer.GetIndexer().ByIndex("groupName", key)
	if err != nil {
		return nil, err
	}

	var assets []*v1alpha2.Asset
	for _, item := range items {
		asset, ok := item.(*v1alpha2.Asset)
		if !ok {
			return nil, fmt.Errorf("Incorrect item type: %T, should be: *Asset", item)
		}

		assets = append(assets, asset)
	}

	return assets, nil
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