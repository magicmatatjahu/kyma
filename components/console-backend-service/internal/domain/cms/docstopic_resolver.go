package cms

import (
	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/shared"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"
	"context"
)

type docsTopicResolver struct {
	docsTopicSvc docsTopicGetter
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
