package genericapi

import (
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type ResourceFilter struct {}

func NewResourceFilter() *ResourceFilter {
	return &ResourceFilter{}
}

func (s *ResourceFilter) Filter(items []unstructured.Unstructured, filterOptions []gqlschema.ResourceFilters) ([]unstructured.Unstructured, error) {
	if filterOptions == nil || len(filterOptions) == 0 {
		return items, nil
	}

	return items, nil
}
