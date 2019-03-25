package cms

import (
	"k8s.io/client-go/tools/cache"
	"github.com/kyma-project/kyma/components/cms-controller-manager/pkg/apis/cms/v1alpha1"
	"fmt"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/cms/pretty"
	"github.com/kyma-project/kyma/components/console-backend-service/pkg/resource"
)

type docsTopicService struct {
	informer cache.SharedIndexInformer
	notifier notifier
}

func newDocsTopicService(informer cache.SharedIndexInformer) (*docsTopicService, error) {
	svc := &docsTopicService{
		informer: informer,
	}

	err := svc.informer.AddIndexers(cache.Indexers{
		"serviceClassName": func(obj interface{}) ([]string, error) {
			entity, err := svc.extractDocsTopic(obj)
			if err != nil {
				return nil, errors.New("Cannot convert item")
			}

			return []string{fmt.Sprintf("%s/%s", entity.Namespace, entity.Name)}, nil
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "while adding indexers")
	}

	notifier := resource.NewNotifier()
	informer.AddEventHandler(notifier)
	svc.notifier = notifier

	return svc, nil
}

func (svc *docsTopicService) Find(namespace, name string) (*v1alpha1.DocsTopic, error) {
	key := fmt.Sprintf("%s/%s", namespace, name)
	item, exists, err := svc.informer.GetStore().GetByKey(key)
	if err != nil || !exists {
		return nil, err
	}

	docsTopic, err := svc.extractDocsTopic(item)
	if err != nil {
		return nil, errors.Wrapf(err, "Incorrect item type: %T, should be: *DocsTopic", item)
	}

	return docsTopic, nil
}

func (svc *docsTopicService) Subscribe(listener resource.Listener) {
	svc.notifier.AddListener(listener)
}

func (svc *docsTopicService) Unsubscribe(listener resource.Listener) {
	svc.notifier.DeleteListener(listener)
}

func (svc *docsTopicService) extractDocsTopic(obj interface{}) (*v1alpha1.DocsTopic, error) {
	u, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, errors.Wrapf(err, "while converting resource %s %s to unstructured", pretty.DocsTopic, obj)
	}

	var docsTopic v1alpha1.DocsTopic
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(u, &docsTopic)
	if err != nil {
		return nil, errors.Wrapf(err, "while converting unstructured to resource %s %s", pretty.DocsTopic, u)
	}

	return &docsTopic, nil
}
