package cms

import (
	"k8s.io/client-go/tools/cache"
	"github.com/kyma-project/kyma/components/cms-controller-manager/pkg/apis/cms/v1alpha1"
	"fmt"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/cms/pretty"
)

type docsTopicService struct {
	informer cache.SharedIndexInformer
}

func newDocsTopicService(informer cache.SharedIndexInformer) (*docsTopicService, error) {
	err := informer.AddIndexers(cache.Indexers{
		"groupName": func(obj interface{}) ([]string, error) {
			entity, err := extractDocsTopic(obj)
			if err != nil {
				return nil, errors.New("Cannot convert item")
			}

			return []string{fmt.Sprintf("%s/%s", entity.Namespace, entity.Name)}, nil
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "while adding indexers")
	}

	return &docsTopicService{
		informer: informer,
	}, nil
}

func (svc *docsTopicService) List(namespace, groupName string) ([]*v1alpha1.DocsTopic, error) {
	key := fmt.Sprintf("%s/%s", namespace, groupName)
	items, err := svc.informer.GetIndexer().ByIndex("groupName", key)
	if err != nil {
		return nil, err
	}

	var docsTopics []*v1alpha1.DocsTopic
	for _, item := range items {
		docsTopic, ok := item.(*v1alpha1.DocsTopic)
		if !ok {
			return nil, fmt.Errorf("Incorrect item type: %T, should be: *DocsTopic", item)
		}

		docsTopics = append(docsTopics, docsTopic)
	}

	return docsTopics, nil
}

func extractDocsTopic(obj interface{}) (*v1alpha1.DocsTopic, error) {
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
