package genericapi

import (
	"fmt"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"
	"github.com/kyma-project/kyma/components/console-backend-service/pkg/dynamic/dynamicinformer"
	notifierResource "github.com/kyma-project/kyma/components/console-backend-service/pkg/resource"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/cache"
)

type ResourceService struct {
	client	  dynamic.ResourceInterface
	informer  cache.SharedIndexInformer
	notifier  notifierResource.Notifier
}

func NewResourceService(dynamicClient dynamic.Interface, informerFactory dynamicinformer.DynamicSharedInformerFactory, s schema.GroupVersionResource) *ResourceService {
	resourceSchema := schema.GroupVersionResource{
		Version:  s.Version,
		Group:    s.Group,
		Resource: s.Resource,
	}

	service := &ResourceService{
		client:   dynamicClient.Resource(resourceSchema),
		informer: informerFactory.ForResource(resourceSchema).Informer(),
	}
	notifier := notifierResource.NewNotifier()
	service.informer.AddEventHandler(notifier)
	service.notifier = notifier

	return service
}

func (s *ResourceService) Get(namespace *string, name string) (interface{}, error) {
	key := name
	if namespace != nil && *namespace != "" {
		key = fmt.Sprintf("%s/%s", *namespace, name)
	}

	item, exists, err := s.informer.GetIndexer().GetByKey(key)
	if err != nil || !exists {
		return nil, err
	}

	return item, nil
}

func (s *ResourceService) List(namespace *string, options *gqlschema.ResourceListOptions) ([]interface{}, error) {
	if namespace != nil && *namespace != "" {
		return s.listInNamespace(*namespace)
	}

	if options != nil {
		if options.AllNamespaces != nil && *options.AllNamespaces {
			return s.informer.GetStore().List(), nil
		}
		if options.InNamespaces != nil {
			return s.listInNamespaces(options.InNamespaces)
		}
	}

	return s.informer.GetStore().List(), nil
}

func (s *ResourceService) Subscribe(listener notifierResource.Listener) {
	s.notifier.AddListener(listener)
}

func (s *ResourceService) Unsubscribe(listener notifierResource.Listener) {
	s.notifier.DeleteListener(listener)
}

func (s *ResourceService) listInNamespace(namespace string) ([]interface{}, error) {
	return s.informer.GetIndexer().ByIndex("namespace", namespace)
}

func (s *ResourceService) listInNamespaces(namespace []string) ([]interface{}, error) {
	out := make([]interface{}, 0)
	for _, namespace := range namespace {
		items, err := s.informer.GetIndexer().ByIndex("namespace", namespace)
		if err != nil {
			return nil, err
		}
		out = append(out, items...)
	}
	return out, nil
}