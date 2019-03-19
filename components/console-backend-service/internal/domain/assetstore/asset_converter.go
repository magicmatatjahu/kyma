package assetstore

import (
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"
	"github.com/kyma-project/kyma/components/asset-store-controller-manager/pkg/apis/assetstore/v1alpha2"
)

type assetConverter struct {}

func (c *assetConverter) ToGQL(item *v1alpha2.Asset) (*gqlschema.Asset, error) {
	if item == nil {
		return nil, nil
	}

	asset := gqlschema.Asset{
		Name: item.Name,
	}

	return &asset, nil
}

func (c *assetConverter) ToGQLs(in []*v1alpha2.Asset) ([]gqlschema.Asset, error) {
	var result []gqlschema.Asset
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
