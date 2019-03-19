package cms

import (
	"k8s.io/client-go/tools/cache"
	"fmt"
	"github.com/pkg/errors"
	"github.com/kyma-project/kyma/components/cms-controller-manager/pkg/apis/cms/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/cms/pretty"
)

type clusterDocsTopicService struct {
	informer cache.SharedIndexInformer
}

func newClusterDocsTopicService(informer cache.SharedIndexInformer) (*clusterDocsTopicService, error) {
	err := informer.AddIndexers(cache.Indexers{
		"groupName": func(obj interface{}) ([]string, error) {
			_, err := extractClusterDocsTopic(obj)
			if err != nil {
				return nil, errors.New("Cannot convert item")
			}

			return []string{fmt.Sprintf("%s", "lol")}, nil
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "while adding indexers")
	}

	return &clusterDocsTopicService{
		informer: informer,
	}, nil
}

func (svc *clusterDocsTopicService) List(groupName string) ([]*v1alpha1.ClusterDocsTopic, error) {
	key := fmt.Sprintf("%s", groupName)
	fmt.Println(key)
	items, err := svc.informer.GetIndexer().ByIndex("groupName", key)
	if err != nil {
		return nil, err
	}

	var clusterDocsTopics []*v1alpha1.ClusterDocsTopic
	for _, item := range items {
		clusterDocsTopic, err := extractClusterDocsTopic(item)
		if err != nil {
			return nil, err
		}

		clusterDocsTopics = append(clusterDocsTopics, clusterDocsTopic)
	}

	return clusterDocsTopics, nil
}

func extractClusterDocsTopic(obj interface{}) (*v1alpha1.ClusterDocsTopic, error) {
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
