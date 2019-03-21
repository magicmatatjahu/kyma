package assetstore

import (
	"github.com/kyma-project/kyma/components/console-backend-service/internal/module"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/dynamic"
	"time"
	"k8s.io/client-go/rest"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"github.com/kyma-project/kyma/components/asset-store-controller-manager/pkg/apis/assetstore/v1alpha2"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/assetstore/disabled"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"
	"context"
)

type PluggableContainer struct {
	*module.Pluggable
	cfg *resolverConfig

	Resolver Resolver
	AssetStoreRetriever *assetStoreRetriever
	informerFactory dynamicinformer.DynamicSharedInformerFactory
}

func New(restConfig *rest.Config, informerResyncPeriod time.Duration) (*PluggableContainer, error) {
	dynamicClient, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return nil, errors.Wrap(err, "while initializing Dynamic Clientset")
	}

	container := &PluggableContainer{
		cfg: &resolverConfig{
			dynamicClient: dynamicClient,
			informerResyncPeriod: informerResyncPeriod,
		},
		Pluggable: module.NewPluggable("content"),
		AssetStoreRetriever: &assetStoreRetriever{},
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

	informerFactory := dynamicinformer.NewDynamicSharedInformerFactory(dynamicClient, informerResyncPeriod)
	r.informerFactory = informerFactory

	clusterAssetService, err := newClusterAssetService(informerFactory.ForResource(schema.GroupVersionResource{
		Version:  v1alpha2.SchemeGroupVersion.Version,
		Group:    v1alpha2.SchemeGroupVersion.Group,
		Resource: "clusterassets",
	}).Informer())
	if err != nil {
		return errors.Wrapf(err, "while creating clusterAsset service")
	}

	assetService, err := newAssetService(informerFactory.ForResource(schema.GroupVersionResource{
		Version:  v1alpha2.SchemeGroupVersion.Version,
		Group:    v1alpha2.SchemeGroupVersion.Group,
		Resource: "assets",
	}).Informer())
	if err != nil {
		return errors.Wrapf(err, "while creating asset service")
	}

	r.Pluggable.EnableAndSyncDynamicInformerFactory(r.informerFactory, func() {
		r.Resolver = &domainResolver{
			clusterAssetResolver: newClusterAssetResolver(clusterAssetService),
			assetResolver: newAssetResolver(assetService),
		}
		r.AssetStoreRetriever.ClusterAssetGetter = clusterAssetService
		r.AssetStoreRetriever.AssetGetter = assetService
		r.AssetStoreRetriever.GqlClusterAssetConverter = &clusterAssetConverter{}
		r.AssetStoreRetriever.GqlAssetConverter = &assetConverter{}
	})

	return nil
}

func (r *PluggableContainer) Disable() error {
	r.Pluggable.Disable(func(disabledErr error) {
		r.Resolver = disabled.NewResolver(disabledErr)
		r.AssetStoreRetriever.ClusterAssetGetter = disabled.NewClusterAssetSvc(disabledErr)
		r.AssetStoreRetriever.AssetGetter = disabled.NewAssetSvc(disabledErr)
		r.AssetStoreRetriever.GqlClusterAssetConverter = disabled.NewGqlClusterAssetConverter(disabledErr)
		r.AssetStoreRetriever.GqlAssetConverter = disabled.NewGqlAssetConverter(disabledErr)
		r.informerFactory = nil
	})

	return nil
}

//go:generate failery -name=Resolver -case=underscore -output disabled -outpkg disabled
type Resolver interface {
	ClusterAssetFilesField(ctx context.Context, obj *gqlschema.ClusterAsset, filterExtensions []string) ([]gqlschema.File, error)
	AssetFilesField(ctx context.Context, obj *gqlschema.Asset, filterExtensions []string) ([]gqlschema.File, error)
	ClusterAssetEventSubscription(ctx context.Context) (<-chan gqlschema.ClusterAssetEvent, error)
	AssetEventSubscription(ctx context.Context) (<-chan gqlschema.AssetEvent, error)
}

type resolverConfig struct {
	dynamicClient             dynamic.Interface
	informerResyncPeriod      time.Duration
}

type domainResolver struct {
	*clusterAssetResolver
	*assetResolver
}
