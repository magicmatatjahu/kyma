package cms

import (
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"
	"context"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/cms/pretty"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlerror"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/shared"
)

type clusterDocsTopicResolver struct {
	clusterDocsTopicSvc clusterDocsTopicGetter
	assetStoreRetriever shared.AssetStoreRetriever
	clusterDocsTopicConverter gqlClusterDocsTopicConverter
}

func newClusterDocsTopicResolver(clusterDocsTopicService *clusterDocsTopicService, assetStoreRetriever shared.AssetStoreRetriever) *clusterDocsTopicResolver {
	return &clusterDocsTopicResolver{
		clusterDocsTopicSvc: clusterDocsTopicService,
		assetStoreRetriever: assetStoreRetriever,
		clusterDocsTopicConverter: &clusterDocsTopicConverter{},
	}
}

func (r *clusterDocsTopicResolver) ClusterDocsTopicsQuery(ctx context.Context, viewContext *string, groupName string) ([]gqlschema.ClusterDocsTopic, error) {
	items, err := r.clusterDocsTopicSvc.List(groupName)
	if err != nil {
		glog.Error(errors.Wrapf(err, "while listing %s", pretty.ClusterDocsTopics))
		return nil, gqlerror.New(err, pretty.ClusterDocsTopics)
	}

	clusterDocsTopics, err := r.clusterDocsTopicConverter.ToGQLs(items)
	if err != nil {
		glog.Error(errors.Wrapf(err, "while converting %s", pretty.ClusterDocsTopics))
		return nil, gqlerror.New(err, pretty.ClusterDocsTopics)
	}

	return clusterDocsTopics, nil
}
