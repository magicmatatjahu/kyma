package genericapi

import (
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sort"
)

type ResourceSort struct {}

func NewResourceSort() *ResourceSort {
	return &ResourceSort{}
}

func (s *ResourceSort) Sort(items []unstructured.Unstructured, sortOptions []gqlschema.ResourceSort) ([]unstructured.Unstructured, error) {
	if sortOptions == nil || len(sortOptions) == 0 {
		s.sortByUID(items)
	}

	return items, nil
}

func (s *ResourceSort) sortByUID(items []unstructured.Unstructured) {
	sort.Slice(items, func(i, j int) bool {
		return items[i].GetUID() < items[j].GetUID()
	})
}
