package extractor

import (
	"github.com/kyma-project/kyma/components/cms-controller-manager/pkg/apis/cms/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"github.com/pkg/errors"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/cms/pretty"
)

type DocsTopicUnstructuredExtractor struct{}

func (ext *DocsTopicUnstructuredExtractor) Single(obj interface{}) (*v1alpha1.DocsTopic, error) {
	if obj == nil {
		return nil, nil
	}

	u, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, errors.Wrapf(err, "while converting resource %s %s to unstructured", pretty.DocsTopic, obj)
	}
	if len(u) == 0 {
		return nil, nil
	}

	var docsTopic v1alpha1.DocsTopic
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(u, &docsTopic)
	if err != nil {
		return nil, errors.Wrapf(err, "while converting unstructured to resource %s %s", pretty.DocsTopic, u)
	}

	return &docsTopic, nil
}
