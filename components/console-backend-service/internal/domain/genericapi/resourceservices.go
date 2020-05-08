package genericapi

import (
	"fmt"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"

	"github.com/kyma-project/kyma/components/console-backend-service/internal/resource"
	notifierResource "github.com/kyma-project/kyma/components/console-backend-service/pkg/resource"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type resourceService struct {
	*resource.Service
	notifier  notifierResource.Notifier
}

type resourcesServices map[string]*resourceService

func newServices(serviceFactory *resource.ServiceFactory, schemas []schema.GroupVersionResource) resourcesServices {
	services := make(map[string]*resourceService, 0)
	for _, s := range schemas {
		svc := &resourceService{
			Service: serviceFactory.ForResource(schema.GroupVersionResource{
				Version:  s.Version,
				Group:    s.Group,
				Resource: s.Resource,
			}),
		}
		notifier := notifierResource.NewNotifier()
		svc.Informer.AddEventHandler(notifier)
		svc.notifier = notifier

		key := prepareKey(s.Resource, s.Group, s.Version)
		services[key] = svc
	}
	return services
}

func (s resourcesServices) retrieveService(schema gqlschema.SchemaResourceInput) *resourceService {
	key := prepareKey(schema.Resource, schema.Group, schema.Version)
	return s[key]
}

func prepareKey(resource, group, version string) string {
	if group == "" {
		return fmt.Sprintf("%s/%s", resource, version)
	}
	return fmt.Sprintf("%s.%s/%s", resource, version)
}
