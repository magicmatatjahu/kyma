package cms

import (
	"github.com/kyma-project/kyma/components/cms-controller-manager/pkg/apis/cms/v1alpha1"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"
)

//go:generate mockery -name=clusterDocsTopicGetter -output=automock -outpkg=automock -case=underscore
//go:generate failery -name=clusterDocsTopicGetter -case=underscore -output disabled -outpkg disabled
type clusterDocsTopicGetter interface {
	List(groupName string) ([]*v1alpha1.ClusterDocsTopic, error)
}

//go:generate mockery -name=gqlClusterDocsTopicConverter -output=automock -outpkg=automock -case=underscore
//go:generate failery -name=gqlClusterDocsTopicConverter -case=underscore -output disabled -outpkg disabled
type gqlClusterDocsTopicConverter interface {
	ToGQL(in *v1alpha1.ClusterDocsTopic) (*gqlschema.ClusterDocsTopic, error)
	ToGQLs(in []*v1alpha1.ClusterDocsTopic) ([]gqlschema.ClusterDocsTopic, error)
}

//go:generate mockery -name=docsTopicGetter -output=automock -outpkg=automock -case=underscore
//go:generate failery -name=docsTopicGetter -case=underscore -output disabled -outpkg disabled
type docsTopicGetter interface {
	List(namespace, groupName string) ([]*v1alpha1.DocsTopic, error)
}

//go:generate mockery -name=gqlDocsTopicConverter -output=automock -outpkg=automock -case=underscore
//go:generate failery -name=gqlDocsTopicConverter -case=underscore -output disabled -outpkg disabled
type gqlDocsTopicConverter interface {
	ToGQL(in *v1alpha1.DocsTopic) (*gqlschema.DocsTopic, error)
	ToGQLs(in []*v1alpha1.DocsTopic) ([]gqlschema.DocsTopic, error)
}