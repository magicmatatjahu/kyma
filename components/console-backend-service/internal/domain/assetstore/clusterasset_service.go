package assetstore

import (
	"k8s.io/client-go/tools/cache"
	"fmt"
	"github.com/pkg/errors"
	"github.com/kyma-project/kyma/components/asset-store-controller-manager/pkg/apis/assetstore/v1alpha2"
)

type clusterAssetService struct {
	informer cache.SharedIndexInformer
}

func newClusterAssetService(informer cache.SharedIndexInformer) (*clusterAssetService, error) {
	err := informer.AddIndexers(cache.Indexers{
		"groupName": func(obj interface{}) ([]string, error) {
			entity, ok := obj.(*v1alpha2.ClusterAsset)
			if !ok {
				return nil, errors.New("Cannot convert item")
			}

			return []string{fmt.Sprintf("%s", entity.Spec.BucketRef)}, nil
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "while adding indexers")
	}

	return &clusterAssetService{
		informer: informer,
	}, nil
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
