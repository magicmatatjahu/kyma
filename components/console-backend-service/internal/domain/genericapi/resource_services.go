package genericapi

import (
	"fmt"
	"github.com/kyma-project/kyma/components/console-backend-service/pkg/dynamic/dynamicinformer"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

type ResourcesServices map[string]*ResourceService

func NewResourceServices(dynamicClient dynamic.Interface, informerFactory dynamicinformer.DynamicSharedInformerFactory, schemas []schema.GroupVersionResource) ResourcesServices {
	services := make(map[string]*ResourceService, 0)

	for _, s := range schemas {
		service := NewResourceService(dynamicClient, informerFactory, s)
		id := prepareResourceServiceID(s.Resource, s.Group, s.Version)
		services[id] = service
	}

	return services
}

func (s ResourcesServices) Get(schema string) *ResourceService {
	return s[schema]
}

func prepareResourceServiceID(resource, group, version string) string {
	if group != "" {
		return fmt.Sprintf("%s.%s/%s", resource, group, version)
	}
	return fmt.Sprintf("%s/%s", resource, version)
}
