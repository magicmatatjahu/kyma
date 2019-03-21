package assetstore

import (
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"
	"context"
)

type clusterAssetResolver struct {
	clusterAssetSvc clusterAssetSvc
	clusterAssetConverter gqlClusterAssetConverter
}

func newClusterAssetResolver(clusterAssetService *clusterAssetService) *clusterAssetResolver {
	return &clusterAssetResolver{
		clusterAssetSvc: clusterAssetService,
		clusterAssetConverter: &clusterAssetConverter{},
	}
}

func (r *clusterAssetResolver) ClusterAssetFilesField(ctx context.Context, obj *gqlschema.ClusterAsset, filterExtension *string) ([]gqlschema.File, error) {
	return nil, nil
}
