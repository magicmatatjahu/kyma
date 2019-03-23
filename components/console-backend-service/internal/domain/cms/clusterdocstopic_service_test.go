package cms_test

import (
	"testing"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
	"github.com/kyma-project/kyma/components/cms-controller-manager/pkg/apis/cms/v1alpha1"
	"k8s.io/client-go/dynamic/fake"
	"github.com/kyma-project/kyma/components/console-backend-service/pkg/dynamic/dynamicinformer"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/cms"
	testingUtils "github.com/kyma-project/kyma/components/console-backend-service/internal/testing"
	"time"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/cms/listener"
)

func TestClusterDocsTopicService_List(t *testing.T) {
	t.Run("Success without parameters", func(t *testing.T) {
		clusterDocsTopic1 := fixUnstructuredClusterDocsTopic(map[string]interface{}{
			"name": "1",
		})
		clusterDocsTopic2 := fixUnstructuredClusterDocsTopic(map[string]interface{}{
			"name": "2",
		})
		clusterDocsTopic3 := fixUnstructuredClusterDocsTopic(map[string]interface{}{
			"name": "3",
		})

		expected := []*v1alpha1.ClusterDocsTopic{
			fixClusterDocsTopic("1", nil),
			fixClusterDocsTopic("2", nil),
			fixClusterDocsTopic("3", nil),
		}

		informer := fixClusterDocsTopicInformer(clusterDocsTopic1, clusterDocsTopic2, clusterDocsTopic3)

		svc, err := cms.NewClusterDocsTopicService(informer)
		require.NoError(t, err)

		testingUtils.WaitForInformerStartAtMost(t, time.Second, informer)

		clusterDocsTopics, err := svc.List(nil, nil)
		require.NoError(t, err)

		assert.Equal(t, expected, clusterDocsTopics)
	})

	t.Run("Success with all parameters", func(t *testing.T) {
		viewContext := "viewContext"
		groupName := "groupName"

		clusterDocsTopic1 := fixUnstructuredClusterDocsTopic(map[string]interface{}{
			"name": "1",
			"labels": map[string]interface{}{
				"viewContext.cms.kyma-project.io": viewContext,
				"groupName.cms.kyma-project.io": groupName,
			},
		})
		clusterDocsTopic2 := fixUnstructuredClusterDocsTopic(map[string]interface{}{
			"name": "2",
		})
		clusterDocsTopic3 := fixUnstructuredClusterDocsTopic(map[string]interface{}{
			"name": "3",
		})

		expected := []*v1alpha1.ClusterDocsTopic{
			fixClusterDocsTopic("1", map[string]string{
				"viewContext.cms.kyma-project.io": viewContext,
				"groupName.cms.kyma-project.io": groupName,
			}),
		}

		informer := fixClusterDocsTopicInformer(clusterDocsTopic1, clusterDocsTopic2, clusterDocsTopic3)

		svc, err := cms.NewClusterDocsTopicService(informer)
		require.NoError(t, err)

		testingUtils.WaitForInformerStartAtMost(t, time.Second, informer)

		clusterDocsTopics, err := svc.List(&viewContext, &groupName)
		require.NoError(t, err)

		assert.Equal(t, expected, clusterDocsTopics)
	})

	t.Run("Success with viewContext parameter", func(t *testing.T) {
		viewContext := "viewContext"

		clusterDocsTopic1 := fixUnstructuredClusterDocsTopic(map[string]interface{}{
			"name": "1",
			"labels": map[string]interface{}{
				"viewContext.cms.kyma-project.io": viewContext,
			},
		})
		clusterDocsTopic2 := fixUnstructuredClusterDocsTopic(map[string]interface{}{
			"name": "2",
		})
		clusterDocsTopic3 := fixUnstructuredClusterDocsTopic(map[string]interface{}{
			"name": "3",
		})

		expected := []*v1alpha1.ClusterDocsTopic{
			fixClusterDocsTopic("1", map[string]string{
				"viewContext.cms.kyma-project.io": viewContext,
			}),
		}

		informer := fixClusterDocsTopicInformer(clusterDocsTopic1, clusterDocsTopic2, clusterDocsTopic3)

		svc, err := cms.NewClusterDocsTopicService(informer)
		require.NoError(t, err)

		testingUtils.WaitForInformerStartAtMost(t, time.Second, informer)

		clusterDocsTopics, err := svc.List(&viewContext, nil)
		require.NoError(t, err)

		assert.Equal(t, expected, clusterDocsTopics)
	})

	t.Run("Success with groupName parameter", func(t *testing.T) {
		groupName := "groupName"

		clusterDocsTopic1 := fixUnstructuredClusterDocsTopic(map[string]interface{}{
			"name": "1",
			"labels": map[string]interface{}{
				"groupName.cms.kyma-project.io": groupName,
			},
		})
		clusterDocsTopic2 := fixUnstructuredClusterDocsTopic(map[string]interface{}{
			"name": "2",
		})
		clusterDocsTopic3 := fixUnstructuredClusterDocsTopic(map[string]interface{}{
			"name": "3",
		})

		expected := []*v1alpha1.ClusterDocsTopic{
			fixClusterDocsTopic("1", map[string]string{
				"groupName.cms.kyma-project.io": groupName,
			}),
		}

		informer := fixClusterDocsTopicInformer(clusterDocsTopic1, clusterDocsTopic2, clusterDocsTopic3)

		svc, err := cms.NewClusterDocsTopicService(informer)
		require.NoError(t, err)

		testingUtils.WaitForInformerStartAtMost(t, time.Second, informer)

		clusterDocsTopics, err := svc.List(nil, &groupName)
		require.NoError(t, err)

		assert.Equal(t, expected, clusterDocsTopics)
	})

	t.Run("NotFound", func(t *testing.T) {
		informer := fixClusterDocsTopicInformer()

		svc, err := cms.NewClusterDocsTopicService(informer)
		require.NoError(t, err)

		testingUtils.WaitForInformerStartAtMost(t, time.Second, informer)

		clusterDocsTopics, err := svc.List(nil, nil)
		require.NoError(t, err)
		assert.Nil(t, clusterDocsTopics)
	})
}

func TestClusterDocsTopicService_ListForServiceClass(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		clusterDocsTopic1 := fixUnstructuredClusterDocsTopic(map[string]interface{}{
			"name": "exampleClassA",
		})
		clusterDocsTopic2 := fixUnstructuredClusterDocsTopic(map[string]interface{}{
			"name": "exampleClassB",
		})
		clusterDocsTopic3 := fixUnstructuredClusterDocsTopic(map[string]interface{}{
			"name": "exampleClassC",
		})

		expected := []*v1alpha1.ClusterDocsTopic{
			fixClusterDocsTopic("exampleClassA", nil),
		}

		informer := fixClusterDocsTopicInformer(clusterDocsTopic1, clusterDocsTopic2, clusterDocsTopic3)

		svc, err := cms.NewClusterDocsTopicService(informer)
		require.NoError(t, err)

		testingUtils.WaitForInformerStartAtMost(t, time.Second, informer)

		clusterDocsTopics, err := svc.ListForServiceClass("exampleClassA")
		require.NoError(t, err)

		assert.Equal(t, expected, clusterDocsTopics)
	})

	t.Run("NotFound", func(t *testing.T) {
		informer := fixClusterDocsTopicInformer()

		svc, err := cms.NewClusterDocsTopicService(informer)
		require.NoError(t, err)

		testingUtils.WaitForInformerStartAtMost(t, time.Second, informer)

		clusterDocsTopics, err := svc.ListForServiceClass("exampleClass")
		require.NoError(t, err)
		assert.Nil(t, clusterDocsTopics)
	})
}

func TestClusterDocsTopicService_Subscribe(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		svc, err := cms.NewClusterDocsTopicService(fixClusterDocsTopicInformer())
		require.NoError(t, err)

		clusterDocsTopicListener := listener.NewClusterDocsTopic(nil, nil, nil)
		svc.Subscribe(clusterDocsTopicListener)
	})

	t.Run("Duplicated", func(t *testing.T) {
		svc, err := cms.NewClusterDocsTopicService(fixClusterDocsTopicInformer())
		require.NoError(t, err)

		clusterDocsTopicListener := listener.NewClusterDocsTopic(nil, nil, nil)
		svc.Subscribe(clusterDocsTopicListener)
		svc.Subscribe(clusterDocsTopicListener)
	})

	t.Run("Multiple", func(t *testing.T) {
		svc, err := cms.NewClusterDocsTopicService(fixClusterDocsTopicInformer())
		require.NoError(t, err)

		clusterDocsTopicListenerA := listener.NewClusterDocsTopic(nil, nil, nil)
		clusterDocsTopicListenerB := listener.NewClusterDocsTopic(nil, nil, nil)

		svc.Subscribe(clusterDocsTopicListenerA)
		svc.Subscribe(clusterDocsTopicListenerB)
	})

	t.Run("Nil", func(t *testing.T) {
		svc, err := cms.NewClusterDocsTopicService(fixClusterDocsTopicInformer())
		require.NoError(t, err)

		svc.Subscribe(nil)
	})
}

func TestClusterDocsTopicService_Unsubscribe(t *testing.T) {
	t.Run("Existing", func(t *testing.T) {
		svc, err := cms.NewClusterDocsTopicService(fixClusterDocsTopicInformer())
		require.NoError(t, err)

		clusterDocsTopicListener := listener.NewClusterDocsTopic(nil, nil, nil)

		svc.Subscribe(clusterDocsTopicListener)
		svc.Unsubscribe(clusterDocsTopicListener)
	})

	t.Run("Duplicated", func(t *testing.T) {
		svc, err := cms.NewClusterDocsTopicService(fixClusterDocsTopicInformer())
		require.NoError(t, err)

		clusterDocsTopicListener := listener.NewClusterDocsTopic(nil, nil, nil)
		svc.Subscribe(clusterDocsTopicListener)
		svc.Subscribe(clusterDocsTopicListener)

		svc.Unsubscribe(clusterDocsTopicListener)
	})

	t.Run("Multiple", func(t *testing.T) {
		svc, err := cms.NewClusterDocsTopicService(fixClusterDocsTopicInformer())
		require.NoError(t, err)

		clusterDocsTopicListenerA := listener.NewClusterDocsTopic(nil, nil, nil)
		clusterDocsTopicListenerB := listener.NewClusterDocsTopic(nil, nil, nil)

		svc.Subscribe(clusterDocsTopicListenerA)
		svc.Subscribe(clusterDocsTopicListenerB)

		svc.Unsubscribe(clusterDocsTopicListenerA)
	})

	t.Run("Nil", func(t *testing.T) {
		svc, err := cms.NewClusterDocsTopicService(fixClusterDocsTopicInformer())
		require.NoError(t, err)

		svc.Unsubscribe(nil)
	})
}

func fixUnstructuredClusterDocsTopic(metadata map[string]interface{}) *unstructured.Unstructured {
	return testingUtils.NewUnstructured(v1alpha1.SchemeGroupVersion.String(), "ClusterDocsTopic", metadata, nil, nil)
}

func fixClusterDocsTopic(name string, labels map[string]string) *v1alpha1.ClusterDocsTopic {
	return &v1alpha1.ClusterDocsTopic{
		TypeMeta: metav1.TypeMeta{
			Kind: "ClusterDocsTopic",
			APIVersion: v1alpha1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Labels:    labels,
		},
	}
}

func fixClusterDocsTopicInformer(objects ...runtime.Object) cache.SharedIndexInformer {
	fakeClient := fake.NewSimpleDynamicClient(runtime.NewScheme(), objects...)
	informerFactory := dynamicinformer.NewDynamicSharedInformerFactory(fakeClient, 0)

	informer := informerFactory.ForResource(schema.GroupVersionResource{
		Version:  v1alpha1.SchemeGroupVersion.Version,
		Group:    v1alpha1.SchemeGroupVersion.Group,
		Resource: "clusterdocstopics",
	}).Informer()

	return informer
}
