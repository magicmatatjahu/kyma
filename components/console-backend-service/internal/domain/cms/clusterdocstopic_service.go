package cms

import (
	"fmt"

	"github.com/kyma-project/kyma/components/cms-controller-manager/pkg/apis/cms/v1alpha1"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/cms/pretty"
	"github.com/kyma-project/kyma/components/console-backend-service/pkg/resource"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/cache"
)

type clusterDocsTopicService struct {
	informer cache.SharedIndexInformer
	notifier notifier
}

func newClusterDocsTopicService(informer cache.SharedIndexInformer) (*clusterDocsTopicService, error) {
	svc := &clusterDocsTopicService{
		informer: informer,
	}

	err := svc.informer.AddIndexers(cache.Indexers{
		"viewContext/groupName": func(obj interface{}) ([]string, error) {
			entity, err := svc.extractClusterDocsTopic(obj)
			if err != nil {
				return nil, errors.New("Cannot convert item")
			}

			return []string{fmt.Sprintf("%s/%s", entity.Labels["viewContext.cms.kyma-project.io"], entity.Labels["groupName.cms.kyma-project.io"])}, nil
		},
		"viewContext": func(obj interface{}) ([]string, error) {
			entity, err := svc.extractClusterDocsTopic(obj)
			if err != nil {
				return nil, errors.New("Cannot convert item")
			}

			return []string{entity.Labels["viewContext.cms.kyma-project.io"]}, nil
		},
		"groupName": func(obj interface{}) ([]string, error) {
			entity, err := svc.extractClusterDocsTopic(obj)
			if err != nil {
				return nil, errors.New("Cannot convert item")
			}

			return []string{entity.Labels["groupName.cms.kyma-project.io"]}, nil
		},
		"serviceClassName": func(obj interface{}) ([]string, error) {
			entity, err := svc.extractClusterDocsTopic(obj)
			if err != nil {
				return nil, errors.New("Cannot convert item")
			}

			return []string{entity.Name}, nil
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

func (svc *clusterDocsTopicService) Find(name string) (*v1alpha1.ClusterDocsTopic, error) {
	item, exists, err := svc.informer.GetStore().GetByKey(name)
	if err != nil || !exists {
		return nil, err
	}

	clusterDocsTopic, err := svc.extractClusterDocsTopic(item)
	if err != nil {
		return nil, errors.Wrapf(err, "Incorrect item type: %T, should be: *ClusterDocsTopic", item)
	}

	return clusterDocsTopic, nil
}

func (svc *clusterDocsTopicService) List(viewContext *string, groupName *string) ([]*v1alpha1.ClusterDocsTopic, error) {
	var items []interface{}
	var err error
	if viewContext != nil && groupName != nil {
		items, err = svc.informer.GetIndexer().ByIndex("viewContext/groupName", fmt.Sprintf("%s/%s", *viewContext, *groupName))
	} else if viewContext != nil {
		items, err = svc.informer.GetIndexer().ByIndex("viewContext", *viewContext)
	} else if groupName != nil {
		items, err = svc.informer.GetIndexer().ByIndex("groupName", *groupName)
	} else {
		items = svc.informer.GetStore().List()
	}

	if err != nil {
		return nil, err
	}

	var clusterDocsTopics []*v1alpha1.ClusterDocsTopic
	for _, item := range items {
		clusterDocsTopic, err := svc.extractClusterDocsTopic(item)
		if err != nil {
			return nil, errors.Wrapf(err, "Incorrect item type: %T, should be: *ClusterDocsTopic", item)
		}

		clusterDocsTopics = append(clusterDocsTopics, clusterDocsTopic)
	}

	return clusterDocsTopics, nil
}

func (svc *clusterDocsTopicService) Subscribe(listener resource.Listener) {
	svc.notifier.AddListener(listener)
}

func (svc *clusterDocsTopicService) Unsubscribe(listener resource.Listener) {
	svc.notifier.DeleteListener(listener)
}

func (svc *clusterDocsTopicService) extractClusterDocsTopic(obj interface{}) (*v1alpha1.ClusterDocsTopic, error) {
	u, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, errors.Wrapf(err, "while converting resource %s %s to unstructured", pretty.ClusterDocsTopic, obj)
	}

	var clusterDocsTopic v1alpha1.ClusterDocsTopic
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(u, &clusterDocsTopic)
	if err != nil {
		return nil, errors.Wrapf(err, "while converting unstructured to resource %s %s", pretty.ClusterDocsTopic, u)
	}

	return &clusterDocsTopic, nil
}
