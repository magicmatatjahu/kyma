package genericapi

import (
	"fmt"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/resource"
	notifierResource "github.com/kyma-project/kyma/components/console-backend-service/pkg/resource"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type ResourceService struct {
	*resource.Service
	notifier  notifierResource.Notifier
}

func NewResourceService(serviceFactory *resource.ServiceFactory, s schema.GroupVersionResource) *ResourceService {
	service := &ResourceService{
		Service: serviceFactory.ForResource(schema.GroupVersionResource{
			Version:  s.Version,
			Group:    s.Group,
			Resource: s.Resource,
		}),
	}
	notifier := notifierResource.NewNotifier()
	service.Informer.AddEventHandler(notifier)
	service.notifier = notifier

	return service
}

func (s *ResourceService) Get(namespace *string, name string) (interface{}, error) {
	key := name
	if namespace != nil {
		key = fmt.Sprintf("%s/%s", *namespace, name)
	}

	item, exists, err := s.Informer.GetIndexer().GetByKey(key)
	if err != nil || !exists {
		return nil, err
	}

	return item, nil
}

func (s *ResourceService) List(namespace *string) ([]interface{}, error) {
	if namespace != nil {
		return s.Informer.GetIndexer().ByIndex("namespace", *namespace)
	}
	return s.Informer.GetStore().List(), nil
}
