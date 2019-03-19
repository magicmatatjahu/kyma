package cms

import (
	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/cms/status"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"
	"github.com/kyma-project/kyma/components/cms-controller-manager/pkg/apis/cms/v1alpha1"
)

type docsTopicConverter struct {
	extractor status.DocsTopicExtractor
}

func (c *docsTopicConverter) ToGQL(item *v1alpha1.DocsTopic) (*gqlschema.DocsTopic, error) {
	if item == nil {
		return nil, nil
	}

	docsTopic := gqlschema.DocsTopic{
		Name: item.Name,
		Namespace: item.Namespace,
		Description: item.Spec.Description,
		DisplayName: item.Spec.DisplayName,
	}

	return &docsTopic, nil
}

func (c *docsTopicConverter) ToGQLs(in []*v1alpha1.DocsTopic) ([]gqlschema.DocsTopic, error) {
	var result []gqlschema.DocsTopic
	for _, u := range in {
		converted, err := c.ToGQL(u)
		if err != nil {
			return nil, err
		}

		if converted != nil {
			result = append(result, *converted)
		}
	}
	return result, nil
}
