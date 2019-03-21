package assetstore

import (
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"
	"context"
)

type assetResolver struct {
	assetSvc assetSvc
	assetConverter gqlAssetConverter
}

func newAssetResolver(assetService *assetService) *assetResolver {
	return &assetResolver{
		assetSvc: assetService,
		assetConverter: &assetConverter{},
	}
}

func (r *assetResolver) AssetFilesField(ctx context.Context, obj *gqlschema.Asset, filterExtension *string) ([]gqlschema.File, error) {
	return nil, nil
}
