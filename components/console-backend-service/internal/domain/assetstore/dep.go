package assetstore

import (
	"github.com/kyma-project/kyma/components/asset-store-controller-manager/pkg/apis/assetstore/v1alpha2"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"
	"github.com/kyma-project/kyma/components/console-backend-service/pkg/resource"
)

type File struct {
	URL      string
	Metadata map[string]interface{}
}

//go:generate mockery -name=clusterAssetSvc -output=automock -outpkg=automock -case=underscore
//go:generate failery -name=clusterAssetSvc -case=underscore -output disabled -outpkg disabled
type clusterAssetSvc interface {
	Find(name string) (*v1alpha2.ClusterAsset, error)
	ListForDocsTopicByType(docsTopicName string, types []string) ([]*v1alpha2.ClusterAsset, error)
	Subscribe(listener resource.Listener)
	Unsubscribe(listener resource.Listener)
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
	ListForDocsTopicByType(namespace, docsTopicName string, types []string) ([]*v1alpha2.Asset, error)
	Subscribe(listener resource.Listener)
	Unsubscribe(listener resource.Listener)
}

//go:generate mockery -name=gqlAssetConverter -output=automock -outpkg=automock -case=underscore
//go:generate failery -name=gqlAssetConverter -case=underscore -output disabled -outpkg disabled
type gqlAssetConverter interface {
	ToGQL(in *v1alpha2.Asset) (*gqlschema.Asset, error)
	ToGQLs(in []*v1alpha2.Asset) ([]gqlschema.Asset, error)
}

//go:generate mockery -name=fileSvc -output=automock -outpkg=automock -case=underscore
type fileSvc interface {
	FilterByExtensions(statusRef *v1alpha2.AssetStatusRef, filterExtensions []string) ([]*File, error)
}

//go:generate mockery -name=gqlFileConverter -output=automock -outpkg=automock -case=underscore
type gqlFileConverter interface {
	ToGQL(file *File) (*gqlschema.File, error)
	ToGQLs(files []*File) ([]gqlschema.File, error)
}

type notifier interface {
	AddListener(observer resource.Listener)
	DeleteListener(observer resource.Listener)
}
