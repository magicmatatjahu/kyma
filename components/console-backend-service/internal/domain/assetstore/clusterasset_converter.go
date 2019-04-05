package assetstore

import (
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"
	"github.com/kyma-project/kyma/components/asset-store-controller-manager/pkg/apis/assetstore/v1alpha2"
)

type clusterAssetConverter struct {}

func (c *clusterAssetConverter) ToGQL(item *v1alpha2.ClusterAsset) (*gqlschema.ClusterAsset, error) {
	if item == nil {
		return nil, nil
	}

	clusterAsset := gqlschema.ClusterAsset{
		Name: item.Name,
	}

	return &clusterAsset, nil
}

func (c *clusterAssetConverter) ToGQLs(in []*v1alpha2.ClusterAsset) ([]gqlschema.ClusterAsset, error) {
	var result []gqlschema.ClusterAsset
	for _, u := range in {
		converted, err := c.ToGQL(u)
		if err != nil {
			return nil, err
		}

		if converted != nil {
			result = append(result, *converted)
		}
	}
	return result, nil
}
