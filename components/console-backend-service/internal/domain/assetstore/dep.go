package assetstore

import (
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"
	"github.com/kyma-project/kyma/components/asset-store-controller-manager/pkg/apis/assetstore/v1alpha2"
)

type File struct {
	URL string
	Metadata map[string]interface{}
}

//go:generate mockery -name=clusterAssetSvc -output=automock -outpkg=automock -case=underscore
//go:generate failery -name=clusterAssetSvc -case=underscore -output disabled -outpkg disabled
type clusterAssetSvc interface {
	Find(name string) (*v1alpha2.ClusterAsset, error)
	List(groupName string) ([]*v1alpha2.ClusterAsset, error)
	ListForDocsTopicByType(docsTopicName string, types []string) ([]*v1alpha2.ClusterAsset, error)
}

//go:generate mockery -name=gqlClusterAssetConverter -output=automock -outpkg=automock -case=underscore
//go:generate failery -name=gqlClusterAssetConverter -case=underscore -output disabled -outpkg disabled
type gqlClusterAssetConverter interface {
	ToGQL(in *v1alpha2.ClusterAsset) (*gqlschema.ClusterAsset, error)
	ToGQLs(in []*v1alpha2.ClusterAsset) ([]gqlschema.ClusterAsset, error)
}

//go:generate mockery -name=assetSvc -output=automock -outpkg=automock -case=underscore
//go:generate failery -name=assetSvc -case=underscore -output disabled -outpkg disabled
type assetSvc interface {
	Find(namespace, name string) (*v1alpha2.Asset, error)
	List(namespace, groupName string) ([]*v1alpha2.Asset, error)
	ListForDocsTopicByType(namespace, docsTopicName string, types []string) ([]*v1alpha2.Asset, error)
}

//go:generate mockery -name=gqlAssetConverter -output=automock -outpkg=automock -case=underscore
//go:generate failery -name=gqlAssetConverter -case=underscore -output disabled -outpkg disabled
type gqlAssetConverter interface {
	ToGQL(in *v1alpha2.Asset) (*gqlschema.Asset, error)
	ToGQLs(in []*v1alpha2.Asset) ([]gqlschema.Asset, error)
}


//go:generate mockery -name=gqlFileConverter -output=automock -outpkg=automock -case=underscore
//go:generate failery -name=gqlFileConverter -case=underscore -output disabled -outpkg disabled
type gqlFileConverter interface {
	ToGQL(in *v1alpha2.Asset) (*gqlschema.File, error)
	ToGQLs(in []*v1alpha2.Asset) ([]gqlschema.File, error)
}
