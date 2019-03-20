package cms

import (
	"github.com/kyma-project/kyma/components/cms-controller-manager/pkg/apis/cms/v1alpha1"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/module"
	"k8s.io/client-go/rest"
	"github.com/pkg/errors"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"time"
	"context"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/cms/disabled"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/shared"
)

type cmsRetriever struct {
	ClusterDocsTopicGetter      shared.ClusterDocsTopicGetter
	DocsTopicGetter      		shared.DocsTopicGetter
}

func (r *cmsRetriever) ClusterDocsTopic() shared.ClusterDocsTopicGetter {
	return r.ClusterDocsTopicGetter
}

func (r *cmsRetriever) DocsTopic() shared.DocsTopicGetter {
	return r.DocsTopicGetter
}

type PluggableContainer struct {
	*module.Pluggable
	cfg *resolverConfig

	Resolver Resolver
	CmsRetriever *cmsRetriever
	informerFactory dynamicinformer.DynamicSharedInformerFactory
}

func New(restConfig *rest.Config, informerResyncPeriod time.Duration, assetStoreRetriever shared.AssetStoreRetriever) (*PluggableContainer, error) {
	dynamicClient, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return nil, errors.Wrap(err, "while initializing Dynamic Clientset")
	}

	container := &PluggableContainer{
		cfg: &resolverConfig{
			dynamicClient: dynamicClient,
			informerResyncPeriod: informerResyncPeriod,
			assetStoreRetriever: assetStoreRetriever,
		},
		Pluggable: module.NewPluggable("content"),
		CmsRetriever: &cmsRetriever{},
	}

	err = container.Disable()
	if err != nil {
		return nil, err
	}

	return container, nil
}

func (r *PluggableContainer) Enable() error {
	informerResyncPeriod := r.cfg.informerResyncPeriod
	dynamicClient := r.cfg.dynamicClient

	assetStoreRetriever := r.cfg.assetStoreRetriever

	informerFactory := dynamicinformer.NewDynamicSharedInformerFactory(dynamicClient, informerResyncPeriod)
	r.informerFactory = informerFactory

	clusterDocsTopicService, err := newClusterDocsTopicService(informerFactory.ForResource(schema.GroupVersionResource{
		Version:  v1alpha1.SchemeGroupVersion.Version,
		Group:    v1alpha1.SchemeGroupVersion.Group,
		Resource: "clusterdocstopics",
	}).Informer())
	if err != nil {
		return errors.Wrapf(err, "while creating clusterDocsTopic service")
	}

	docsTopicService, err := newDocsTopicService(informerFactory.ForResource(schema.GroupVersionResource{
		Version:  v1alpha1.SchemeGroupVersion.Version,
		Group:    v1alpha1.SchemeGroupVersion.Group,
		Resource: "docstopics",
	}).Informer())
	if err != nil {
		return errors.Wrapf(err, "while creating docsTopic service")
	}

	r.Pluggable.EnableAndSyncDynamicInformerFactory(r.informerFactory, func() {
		r.Resolver = &domainResolver{
			clusterDocsTopicResolver: newClusterDocsTopicResolver(clusterDocsTopicService, assetStoreRetriever),
			docsTopicResolver: newDocsTopicResolver(docsTopicService, assetStoreRetriever),
		}
		r.CmsRetriever.ClusterDocsTopicGetter = clusterDocsTopicService
		r.CmsRetriever.DocsTopicGetter = docsTopicService
	})

	return nil
}

func (r *PluggableContainer) Disable() error {
	r.Pluggable.Disable(func(disabledErr error) {
		r.Resolver = disabled.NewResolver(disabledErr)
		r.CmsRetriever.ClusterDocsTopicGetter = disabled.NewClusterDocsTopicGetter(disabledErr)
		r.CmsRetriever.DocsTopicGetter = disabled.NewDocsTopicGetter(disabledErr)
		r.informerFactory = nil
	})

	return nil
}

type resolverConfig struct {
	dynamicClient             dynamic.Interface
	informerResyncPeriod      time.Duration
	assetStoreRetriever       shared.AssetStoreRetriever
}

//go:generate failery -name=Resolver -case=underscore -output disabled -outpkg disabled
type Resolver interface {
	ClusterDocsTopicsQuery(ctx context.Context, viewContext *string, groupName *string) ([]gqlschema.ClusterDocsTopic, error)
	ClusterDocsTopicAssetsField(ctx context.Context, obj *gqlschema.ClusterDocsTopic, typeArg *string) ([]gqlschema.ClusterAsset, error)
	DocsTopicAssetsField(ctx context.Context, obj *gqlschema.DocsTopic, typeArg *string) ([]gqlschema.Asset, error)
}

type domainResolver struct {
	*clusterDocsTopicResolver
	*docsTopicResolver
}
