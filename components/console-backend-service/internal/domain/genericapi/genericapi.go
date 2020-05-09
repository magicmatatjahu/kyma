package genericapi

import (
	"context"

	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/module"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/resource"
)

type PluggableContainer struct {
	*module.Pluggable

	Resolver
	serviceFactory *resource.ServiceFactory
}

func New(serviceFactory *resource.ServiceFactory) (*PluggableContainer, error) {
	resolver := &PluggableContainer{
		Pluggable: module.NewPluggable("serverless"),
		serviceFactory: serviceFactory,
	}

	err := resolver.Disable()
	if err != nil {
		return nil, err
	}

	return resolver, nil
}

func (r *PluggableContainer) Enable() error {
	r.Pluggable.EnableAndSyncDynamicInformerFactory(r.serviceFactory.InformerFactory, func() {
		r.Resolver = &domainResolver{
			ResourceResolver: NewResourceResolver(r.serviceFactory),
		}
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
	Get(ctx context.Context, schema gqlschema.SchemaResourceInput, name string, namespace *string) (*gqlschema.Resource, error)
	List(ctx context.Context, schema gqlschema.SchemaResourceInput, namespace *string) (gqlschema.ResourceListOutput, error)

	ResourceSpec(ctx context.Context, obj *gqlschema.Resource, fields []gqlschema.ResourceFieldInput, rootField *string) (gqlschema.ResourceSpecOutput, error)
	ResourceSubResource(ctx context.Context, parent *gqlschema.Resource, schema gqlschema.SchemaResourceInput, name string, namespace *string) (*gqlschema.Resource, error)
	ResourceSubResources(ctx context.Context, parent *gqlschema.Resource, schema gqlschema.SchemaResourceInput, namespace *string) (gqlschema.ResourceListOutput, error)
}

type domainResolver struct {
	*ResourceResolver
}
