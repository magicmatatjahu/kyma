package shared

import (
	"github.com/kyma-project/kyma/components/asset-store-controller-manager/pkg/apis/assetstore/v1alpha2"
)

//go:generate mockery -name=AssetStoreRetriever -output=automock -outpkg=automock -case=underscore
type AssetStoreRetriever interface {
	ClusterAsset() ClusterAssetGetter
	Asset() AssetGetter
}

//go:generate mockery -name=ClusterAssetGetter -output=automock -outpkg=automock -case=underscore
type ClusterAssetGetter interface {
	List(groupName string) ([]*v1alpha2.ClusterAsset, error)
}

//go:generate mockery -name=AssetGetter -output=automock -outpkg=automock -case=underscore
type AssetGetter interface {
	List(namespace, groupName string) ([]*v1alpha2.Asset, error)
}
