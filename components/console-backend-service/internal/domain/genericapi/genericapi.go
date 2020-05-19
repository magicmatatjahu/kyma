package genericapi

import (
	"context"
	"github.com/kyma-project/kyma/components/console-backend-service/pkg/dynamic/dynamicinformer"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"time"

	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/module"
)

type PluggableContainer struct {
	*module.Pluggable

	Resolver
	dynamicClient   dynamic.Interface
	informerFactory dynamicinformer.DynamicSharedInformerFactory
}

func New(restConfig *rest.Config, informerResyncPeriod time.Duration) (*PluggableContainer, error) {
	dynamicClient, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		return nil, errors.Wrap(err, "while initializing Dynamic Clientset")
	}
	informerFactory := dynamicinformer.NewDynamicSharedInformerFactory(dynamicClient, informerResyncPeriod)

	resolver := &PluggableContainer{
		dynamicClient:        dynamicClient,
		informerFactory: 	  informerFactory,
		Pluggable: module.NewPluggable("serverless"),
	}

	err = resolver.Disable()
	if err != nil {
		return nil, err
	}

	return resolver, nil
}

func (r *PluggableContainer) Enable() error {
	schemas := []schema.GroupVersionResource{
		{
			Version:  "v1beta1",
			Group:    "servicecatalog.k8s.io",
			Resource: "serviceinstances",
		},
		{
			Version:  "v1beta1",
			Group:    "servicecatalog.k8s.io",
			Resource: "clusterserviceclasses",
		},
		{
			Version:  "v1beta1",
			Group:    "servicecatalog.k8s.io",
			Resource: "servicebindings",
		},
	}
	services := NewResourceServices(r.dynamicClient, r.informerFactory, schemas)
	pager := NewResourcePager()
	converter := NewResourceConverter(pager)
	sort := NewResourceSort()
	filter := NewResourceFilter()

	r.Pluggable.EnableAndSyncDynamicInformerFactory(r.informerFactory, func() {
		r.Resolver = NewResourceResolver(services, converter, pager, sort, filter)
	})

	return nil
}

func (r *PluggableContainer) Disable() error {
	r.Pluggable.Disable(func(disabledErr error) {
		//r.Resolver = disabled.NewResolver(disabledErr)
	})

	return nil
}

//go:generate failery -name=Resolver -case=underscore -output disabled -outpkg disabled
type Resolver interface {
	Get(ctx context.Context, schema string, name string, namespace *string) (*gqlschema.Resource, error)
	List(ctx context.Context, schema string, namespace *string, pager *gqlschema.ResourcePager, options *gqlschema.ResourceListOptions) (gqlschema.ResourceListOutput, error)

	ResourceSpec(ctx context.Context, obj *gqlschema.Resource, fields []gqlschema.ResourceFieldInput, rootField *string) (gqlschema.JSON, error)
	ResourceSubResource(ctx context.Context, parent *gqlschema.Resource, schema string, name string, namespace *string) (*gqlschema.Resource, error)
	ResourceSubResources(ctx context.Context, parent *gqlschema.Resource, schema string, namespace *string, pager *gqlschema.ResourcePager, options *gqlschema.ResourceListOptions) (gqlschema.ResourceListOutput, error)

	Watch(ctx context.Context, schema string, namespace, name *string) (<-chan gqlschema.ResourceEvent, error)
}
