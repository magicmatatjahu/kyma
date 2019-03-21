package cms

import (
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"
	"context"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/cms/listener"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/cms/pretty"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlerror"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/shared"
	"fmt"
	"github.com/kyma-project/kyma/components/cms-controller-manager/pkg/apis/cms/v1alpha1"
)

type clusterDocsTopicResolver struct {
	clusterDocsTopicSvc clusterDocsTopicSvc
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

func (r *clusterDocsTopicResolver) ClusterDocsTopicsQuery(ctx context.Context, viewContext *string, groupName *string) ([]gqlschema.ClusterDocsTopic, error) {
	items, err := r.clusterDocsTopicSvc.List(*groupName)
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

func (r *clusterDocsTopicResolver) ClusterDocsTopicAssetsField(ctx context.Context, obj *gqlschema.ClusterDocsTopic, typeArg *string) ([]gqlschema.ClusterAsset, error) {
	fmt.Println(typeArg)
	return nil, nil
}

func (r *clusterDocsTopicResolver) ClusterDocsTopicEventSubscription(ctx context.Context) (<-chan gqlschema.ClusterDocsTopicEvent, error) {
	channel := make(chan gqlschema.ClusterDocsTopicEvent, 1)
	filter := func(entity *v1alpha1.ClusterDocsTopic) bool {
		return true
	}

	clusterDocsTopicListener := listener.NewClusterDocsTopic(channel, filter, r.clusterDocsTopicConverter)

	r.clusterDocsTopicSvc.Subscribe(clusterDocsTopicListener)
	go func() {
		defer close(channel)
		defer r.clusterDocsTopicSvc.Unsubscribe(clusterDocsTopicListener)
		<-ctx.Done()
	}()

	return channel, nil
}
