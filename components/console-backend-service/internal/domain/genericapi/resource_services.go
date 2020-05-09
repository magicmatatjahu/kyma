package genericapi

import (
	"fmt"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/resource"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type ResourcesServices map[string]*ResourceService

func NewResourceServices(serviceFactory *resource.ServiceFactory, schemas []schema.GroupVersionResource) ResourcesServices {
	services := make(map[string]*ResourceService, 0)

	for _, s := range schemas {
		service := NewResourceService(serviceFactory, s)
		id := prepareResourceServiceID(s.Resource, s.Group, s.Version)
		services[id] = service
	}

	return services
}

func (s ResourcesServices) Get(schema gqlschema.SchemaResourceInput) *ResourceService {
	id := prepareResourceServiceID(schema.Resource, schema.Group, schema.Version)
	return s[id]
}

func prepareResourceServiceID(resource, group, version string) string {
	if group != "" {
		return fmt.Sprintf("%s.%s/%s", resource, group, version)
	}
	return fmt.Sprintf("%s/%s", resource, version)
}
