package extractor

import (
	"github.com/kyma-project/kyma/components/cms-controller-manager/pkg/apis/cms/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"github.com/pkg/errors"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/cms/pretty"
)

type ClusterDocsTopicUnstructuredExtractor struct{}

func (ext *ClusterDocsTopicUnstructuredExtractor) Single(obj interface{}) (*v1alpha1.ClusterDocsTopic, error) {
	if obj == nil {
		return nil, nil
	}

	u, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, errors.Wrapf(err, "while converting resource %s %s to unstructured", pretty.ClusterDocsTopic, obj)
	}
	if len(u) == 0 {
		return nil, nil
	}

	var clusterDocsTopic v1alpha1.ClusterDocsTopic
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(u, &clusterDocsTopic)
	if err != nil {
		return nil, errors.Wrapf(err, "while converting unstructured to resource %s %s", pretty.ClusterDocsTopic, u)
	}

	return &clusterDocsTopic, nil
}
