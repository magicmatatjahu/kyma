package genericapi

import (
	"context"
	"github.com/kyma-project/kyma/components/function-controller/pkg/apis/serverless/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime/schema"

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
	schemas := []schema.GroupVersionResource{
		{
			Version:  v1alpha1.GroupVersion.Version,
			Group:    v1alpha1.GroupVersion.Group,
			Resource: "functions",
		},
		{
			Version:  "v1",
			Group:    "",
			Resource: "pods",
		},
	}
	services := newServices(r.serviceFactory, schemas)
	converter := newResourceConverter()

	r.Pluggable.EnableAndSyncDynamicInformerFactory(r.serviceFactory.InformerFactory, func() {
		r.Resolver = &domainResolver{
			resourceQueryResolver: newResourceQueryResolver(services, converter),
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

	ResourceFields(ctx context.Context, obj *gqlschema.Resource, fields []gqlschema.ResourceFieldInput) (gqlschema.JSON, error)
	ResourceSubResources(ctx context.Context, obj *gqlschema.Resource, resources []gqlschema.SubResourceInput) ([]gqlschema.SubResourceOutput, error)
}

type domainResolver struct {
	*resourceQueryResolver
}
