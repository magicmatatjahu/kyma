package cms

import (
	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/shared"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"
	"context"
	"github.com/kyma-project/kyma/components/cms-controller-manager/pkg/apis/cms/v1alpha1"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/cms/listener"
	"github.com/golang/glog"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/cms/pretty"
	assetstorePretty "github.com/kyma-project/kyma/components/console-backend-service/internal/domain/assetstore/pretty"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlerror"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/module"
	"github.com/pkg/errors"
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

func (r *docsTopicResolver) DocsTopicAssetsField(ctx context.Context, obj *gqlschema.DocsTopic, types []string) ([]gqlschema.Asset, error) {
	if obj == nil {
		glog.Error(errors.New("%s cannot be empty in order to resolve `assets` field"), pretty.DocsTopic)
		return nil, gqlerror.NewInternal()
	}

	items, err := r.assetStoreRetriever.Asset().ListForDocsTopicByType(obj.Namespace, obj.Name, types)
	if err != nil {
		if module.IsDisabledModuleError(err) {
			return nil, err
		}
		glog.Error(errors.Wrapf(err, "while gathering %s for %s %s", assetstorePretty.Assets, pretty.DocsTopic, obj.Name))
		return nil, gqlerror.New(err, assetstorePretty.Assets)
	}

	assets, err := r.assetStoreRetriever.AssetConverter().ToGQLs(items)
	if err != nil {
		glog.Error(errors.Wrapf(err, "while converting %s", assetstorePretty.Assets))
		return nil, gqlerror.New(err, assetstorePretty.Assets)
	}

	return assets, nil
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
