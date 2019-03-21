package cms

import (
	"github.com/kyma-project/kyma/components/cms-controller-manager/pkg/apis/cms/v1alpha1"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"
	"github.com/kyma-project/kyma/components/console-backend-service/pkg/resource"
)

//go:generate mockery -name=clusterDocsTopicSvc -output=automock -outpkg=automock -case=underscore
//go:generate failery -name=clusterDocsTopicSvc -case=underscore -output disabled -outpkg disabled
type clusterDocsTopicSvc interface {
	List(groupName string) ([]*v1alpha1.ClusterDocsTopic, error)
	ListForServiceClass(className string) ([]*v1alpha1.ClusterDocsTopic, error)
	Subscribe(listener resource.Listener)
	Unsubscribe(listener resource.Listener)
}

//go:generate mockery -name=gqlClusterDocsTopicConverter -output=automock -outpkg=automock -case=underscore
//go:generate failery -name=gqlClusterDocsTopicConverter -case=underscore -output disabled -outpkg disabled
type gqlClusterDocsTopicConverter interface {
	ToGQL(in *v1alpha1.ClusterDocsTopic) (*gqlschema.ClusterDocsTopic, error)
	ToGQLs(in []*v1alpha1.ClusterDocsTopic) ([]gqlschema.ClusterDocsTopic, error)
}

//go:generate mockery -name=docsTopicSvc -output=automock -outpkg=automock -case=underscore
//go:generate failery -name=docsTopicSvc -case=underscore -output disabled -outpkg disabled
type docsTopicSvc interface {
	List(namespace, groupName string) ([]*v1alpha1.DocsTopic, error)
	ListForServiceClass(namespace, className string) ([]*v1alpha1.DocsTopic, error)
	Subscribe(listener resource.Listener)
	Unsubscribe(listener resource.Listener)
}

//go:generate mockery -name=gqlDocsTopicConverter -output=automock -outpkg=automock -case=underscore
//go:generate failery -name=gqlDocsTopicConverter -case=underscore -output disabled -outpkg disabled
type gqlDocsTopicConverter interface {
	ToGQL(in *v1alpha1.DocsTopic) (*gqlschema.DocsTopic, error)
	ToGQLs(in []*v1alpha1.DocsTopic) ([]gqlschema.DocsTopic, error)
}

type notifier interface {
	AddListener(observer resource.Listener)
	DeleteListener(observer resource.Listener)
}
