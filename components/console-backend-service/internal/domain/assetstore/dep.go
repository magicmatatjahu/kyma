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
