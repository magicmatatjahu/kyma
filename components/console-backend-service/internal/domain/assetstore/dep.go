package assetstore

import (
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"
	"github.com/kyma-project/kyma/components/asset-store-controller-manager/pkg/apis/assetstore/v1alpha2"
)

//go:generate mockery -name=clusterAssetGetter -output=automock -outpkg=automock -case=underscore
//go:generate failery -name=clusterAssetGetter -case=underscore -output disabled -outpkg disabled
type clusterAssetGetter interface {
	List(groupName string) ([]*v1alpha2.ClusterAsset, error)
}

//go:generate mockery -name=gqlClusterAssetConverter -output=automock -outpkg=automock -case=underscore
type gqlClusterAssetConverter interface {
	ToGQL(in *v1alpha2.ClusterAsset) (*gqlschema.ClusterAsset, error)
	ToGQLs(in []*v1alpha2.ClusterAsset) ([]gqlschema.ClusterAsset, error)
}

//go:generate mockery -name=assetGetter -output=automock -outpkg=automock -case=underscore
//go:generate failery -name=assetGetter -case=underscore -output disabled -outpkg disabled
type assetGetter interface {
	List(namespace, groupName string) ([]*v1alpha2.Asset, error)
}

//go:generate mockery -name=gqlAssetConverter -output=automock -outpkg=automock -case=underscore
type gqlAssetConverter interface {
	ToGQL(in *v1alpha2.Asset) (*gqlschema.Asset, error)
	ToGQLs(in []*v1alpha2.Asset) ([]gqlschema.Asset, error)
}
