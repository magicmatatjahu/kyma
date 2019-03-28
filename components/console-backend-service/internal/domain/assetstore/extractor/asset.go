package extractor

import (
	"k8s.io/apimachinery/pkg/runtime"
	"github.com/pkg/errors"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/assetstore/pretty"
	"github.com/kyma-project/kyma/components/asset-store-controller-manager/pkg/apis/assetstore/v1alpha2"
)

type AssetUnstructuredExtractor struct{}

func (ext *AssetUnstructuredExtractor) Single(obj interface{}) (*v1alpha2.Asset, error) {
	if obj == nil {
		return nil, nil
	}

	u, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, errors.Wrapf(err, "while converting resource %s %s to unstructured", pretty.Asset, obj)
	}
	if len(u) == 0 {
		return nil, nil
	}

	var asset v1alpha2.Asset
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(u, &asset)
	if err != nil {
		return nil, errors.Wrapf(err, "while converting unstructured to resource %s %s", pretty.Asset, u)
	}

	return &asset, nil
}
