package shared

import "github.com/kyma-project/kyma/components/cms-controller-manager/pkg/apis/cms/v1alpha1"

//go:generate mockery -name=CmsRetriever -output=automock -outpkg=automock -case=underscore
type CmsRetriever interface {
	ClusterDocsTopic() ClusterDocsTopicGetter
	DocsTopic() DocsTopicGetter
}

//go:generate mockery -name=ClusterDocsTopicGetter -output=automock -outpkg=automock -case=underscore
type ClusterDocsTopicGetter interface {
	List(groupName string) ([]*v1alpha1.ClusterDocsTopic, error)
}

//go:generate mockery -name=DocsTopicGetter -output=automock -outpkg=automock -case=underscore
type DocsTopicGetter interface {
	List(namespace, groupName string) ([]*v1alpha1.DocsTopic, error)
}
