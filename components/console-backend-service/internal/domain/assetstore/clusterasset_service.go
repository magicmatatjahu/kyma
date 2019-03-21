package assetstore

import (
	"k8s.io/client-go/tools/cache"
	"fmt"
	"github.com/pkg/errors"
	"github.com/kyma-project/kyma/components/asset-store-controller-manager/pkg/apis/assetstore/v1alpha2"
	"k8s.io/apimachinery/pkg/runtime"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/assetstore/pretty"
)

type clusterAssetService struct {
	informer cache.SharedIndexInformer
}

func newClusterAssetService(informer cache.SharedIndexInformer) (*clusterAssetService, error) {
	svc := &clusterAssetService{
		informer: informer,
	}

	err := svc.informer.AddIndexers(cache.Indexers{
		"groupName": func(obj interface{}) ([]string, error) {
			entity, err := svc.extractClusterAsset(obj)
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


func (svc *clusterAssetService) List(groupName string) ([]*v1alpha2.ClusterAsset, error) {
	key := fmt.Sprintf("%s", groupName)
	items, err := svc.informer.GetIndexer().ByIndex("groupName", key)
	if err != nil {
		return nil, err
	}

	var clusterAssets []*v1alpha2.ClusterAsset
	for _, item := range items {
		clusterAsset, ok := item.(*v1alpha2.ClusterAsset)
		if !ok {
			return nil, fmt.Errorf("Incorrect item type: %T, should be: *ClusterAsset", item)
		}

		clusterAssets = append(clusterAssets, clusterAsset)
	}

	return clusterAssets, nil
}

func (svc *clusterAssetService) ListByType(typeArg *string) ([]*v1alpha2.ClusterAsset, error) {
	var items []interface{}
	var err error
	if typeArg != nil {
		key := fmt.Sprintf("%s", *typeArg)
		items, err = svc.informer.GetIndexer().ByIndex("groupName", key)
	} else {
		items = svc.informer.GetStore().List()
	}

	if err != nil {
		return nil, err
	}

	var clusterAssets []*v1alpha2.ClusterAsset
	for _, item := range items {
		clusterAsset, ok := item.(*v1alpha2.ClusterAsset)
		if !ok {
			return nil, fmt.Errorf("Incorrect item type: %T, should be: *ClusterAsset", item)
		}

		clusterAssets = append(clusterAssets, clusterAsset)
	}

	return clusterAssets, nil
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