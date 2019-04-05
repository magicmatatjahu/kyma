package cms

import (
	"k8s.io/client-go/tools/cache"
	"fmt"
	"github.com/pkg/errors"
	"github.com/kyma-project/kyma/components/cms-controller-manager/pkg/apis/cms/v1alpha1"
)

type clusterDocsTopicService struct {
	informer cache.SharedIndexInformer
}

func newClusterDocsTopicService(informer cache.SharedIndexInformer) (*clusterDocsTopicService, error) {
	err := informer.AddIndexers(cache.Indexers{
		"groupName": func(obj interface{}) ([]string, error) {
			entity, ok := obj.(*v1alpha1.ClusterDocsTopic)
			if !ok {
				return nil, errors.New("Cannot convert item")
			}

			return []string{fmt.Sprintf("%s", entity.Name)}, nil
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "while adding indexers")
	}

	return &clusterDocsTopicService{
		informer: informer,
	}, nil
}

//func (svc *clusterDocsTopicService) FindByGroupName(groupName string) (*v1alpha1.ClusterDocsTopic, error) {
//	key := fmt.Sprintf("%s", groupName)
//	items, err := svc.informer.GetIndexer().ByIndex("groupName", key)
//	if err != nil {
//		return nil, err
//	}
//
//	if len(items) == 0 {
//		return nil, nil
//	}
//
//	if len(items) > 1 {
//		return nil, fmt.Errorf("Multiple ClusterDocsTopic resources with the same groupName %s", groupName)
//	}
//
//	item := items[0]
//	clusterDocsTopic, ok := item.(*v1alpha1.ClusterDocsTopic)
//	if !ok {
//		return nil, fmt.Errorf("Incorrect item type: %T, should be: *ClusterDocsTopic", item)
//	}
//
//	return clusterDocsTopic, nil
//}

func (svc *clusterDocsTopicService) List(groupName string) ([]*v1alpha1.ClusterDocsTopic, error) {
	key := fmt.Sprintf("%s", groupName)
	items, err := svc.informer.GetIndexer().ByIndex("groupName", key)
	if err != nil {
		return nil, err
	}

	var clusterDocsTopics []*v1alpha1.ClusterDocsTopic
	for _, item := range items {
		clusterDocsTopic, ok := item.(*v1alpha1.ClusterDocsTopic)
		if !ok {
			return nil, fmt.Errorf("Incorrect item type: %T, should be: *ClusterDocsTopic", item)
		}

		clusterDocsTopics = append(clusterDocsTopics, clusterDocsTopic)
	}

	return clusterDocsTopics, nil
}

