package cms

import (
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"
	"github.com/kyma-project/kyma/components/cms-controller-manager/pkg/apis/cms/v1alpha1"
)

type clusterDocsTopicConverter struct {}

func (c *clusterDocsTopicConverter) ToGQL(item *v1alpha1.ClusterDocsTopic) (*gqlschema.ClusterDocsTopic, error) {
	if item == nil {
		return nil, nil
	}

	clusterDocsTopic := gqlschema.ClusterDocsTopic{
		Name: item.Name,
		Description: item.Spec.Description,
		DisplayName: item.Spec.DisplayName,
		GroupName: item.Labels["groupName.cms.kyma-project.io"],
	}

	return &clusterDocsTopic, nil
}

func (c *clusterDocsTopicConverter) ToGQLs(in []*v1alpha1.ClusterDocsTopic) ([]gqlschema.ClusterDocsTopic, error) {
	var result []gqlschema.ClusterDocsTopic
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
