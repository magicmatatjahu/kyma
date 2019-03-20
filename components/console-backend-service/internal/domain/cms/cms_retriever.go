package cms

import "github.com/kyma-project/kyma/components/console-backend-service/internal/domain/shared"

type cmsRetriever struct {
	ClusterDocsTopicGetter      		shared.ClusterDocsTopicGetter
	DocsTopicGetter      				shared.DocsTopicGetter
	GqlClusterDocsTopicConverter      	shared.GqlClusterDocsTopicConverter
	GqlDocsTopicConverter      			shared.GqlDocsTopicConverter
}

func (r *cmsRetriever) ClusterDocsTopic() shared.ClusterDocsTopicGetter {
	return r.ClusterDocsTopicGetter
}

func (r *cmsRetriever) DocsTopic() shared.DocsTopicGetter {
	return r.DocsTopicGetter
}

func (r *cmsRetriever) ClusterDocsTopicConverter() shared.GqlClusterDocsTopicConverter {
	return r.GqlClusterDocsTopicConverter
}

func (r *cmsRetriever) DocsTopicConverter() shared.GqlDocsTopicConverter {
	return r.GqlDocsTopicConverter
}
