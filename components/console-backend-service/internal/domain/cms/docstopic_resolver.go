package cms

import (
	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/shared"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"
	"context"
	"github.com/kyma-project/kyma/components/cms-controller-manager/pkg/apis/cms/v1alpha1"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/cms/listener"
)

type docsTopicResolver struct {
	docsTopicSvc docsTopicSvc
	assetStoreRetriever shared.AssetStoreRetriever
	docsTopicConverter gqlDocsTopicConverter
}

func newDocsTopicResolver(docsTopicService *docsTopicService, assetStoreRetriever shared.AssetStoreRetriever) *docsTopicResolver {
	return &docsTopicResolver{
		docsTopicSvc: docsTopicService,
		assetStoreRetriever: assetStoreRetriever,
		docsTopicConverter: &docsTopicConverter{},
	}
}

func (r *docsTopicResolver) DocsTopicAssetsField(ctx context.Context, obj *gqlschema.DocsTopic, typeArg *string) ([]gqlschema.Asset, error) {
	return nil, nil
}

func (r *docsTopicResolver) DocsTopicEventSubscription(ctx context.Context) (<-chan gqlschema.DocsTopicEvent, error) {
	channel := make(chan gqlschema.DocsTopicEvent, 1)
	filter := func(entity *v1alpha1.DocsTopic) bool {
		return true
	}

	docsTopicListener := listener.NewDocsTopic(channel, filter, r.docsTopicConverter)

	r.docsTopicSvc.Subscribe(docsTopicListener)
	go func() {
		defer close(channel)
		defer r.docsTopicSvc.Unsubscribe(docsTopicListener)
		<-ctx.Done()
	}()

	return channel, nil
}
