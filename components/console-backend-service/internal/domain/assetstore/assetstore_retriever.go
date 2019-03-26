package assetstore

import "github.com/kyma-project/kyma/components/console-backend-service/internal/domain/shared"

type assetStoreRetriever struct {
	ClusterAssetGetter       shared.ClusterAssetGetter
	AssetGetter              shared.AssetGetter
	GqlClusterAssetConverter shared.GqlClusterAssetConverter
	GqlAssetConverter        shared.GqlAssetConverter
}

func (r *assetStoreRetriever) ClusterAsset() shared.ClusterAssetGetter {
	return r.ClusterAssetGetter
}

func (r *assetStoreRetriever) Asset() shared.AssetGetter {
	return r.AssetGetter
}

func (r *assetStoreRetriever) ClusterAssetConverter() shared.GqlClusterAssetConverter {
	return r.GqlClusterAssetConverter
}

func (r *assetStoreRetriever) AssetConverter() shared.GqlAssetConverter {
	return r.GqlAssetConverter
}
