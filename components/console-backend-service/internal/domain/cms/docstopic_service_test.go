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

const (
	DocsTopicNamespace = "DocsTopicNamespace"
)

func TestDocsTopicService_ListForServiceClass(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		docsTopic1 := fixUnstructuredDocsTopic(map[string]interface{}{
			"name": "exampleClassA",
			"namespace": DocsTopicNamespace,
		})
		docsTopic2 := fixUnstructuredDocsTopic(map[string]interface{}{
			"name": "exampleClassB",
			"namespace": DocsTopicNamespace,
		})
		docsTopic3 := fixUnstructuredDocsTopic(map[string]interface{}{
			"name": "exampleClassC",
			"namespace": DocsTopicNamespace,
		})

		expected := []*v1alpha1.DocsTopic{
			fixDocsTopic("exampleClassA", nil),
		}

		informer := fixDocsTopicInformer(docsTopic1, docsTopic2, docsTopic3)

		svc, err := cms.NewDocsTopicService(informer)
		require.NoError(t, err)

		testingUtils.WaitForInformerStartAtMost(t, time.Second, informer)

		docsTopics, err := svc.ListForServiceClass(DocsTopicNamespace, "exampleClassA")
		require.NoError(t, err)

		assert.Equal(t, expected, docsTopics)
	})

	t.Run("NotFound", func(t *testing.T) {
		informer := fixDocsTopicInformer()

		svc, err := cms.NewDocsTopicService(informer)
		require.NoError(t, err)

		testingUtils.WaitForInformerStartAtMost(t, time.Second, informer)

		docsTopics, err := svc.ListForServiceClass(DocsTopicNamespace, "exampleClass")
		require.NoError(t, err)
		assert.Nil(t, docsTopics)
	})
}

func TestDocsTopicService_Subscribe(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		svc, err := cms.NewDocsTopicService(fixDocsTopicInformer())
		require.NoError(t, err)

		docsTopicListener := listener.NewDocsTopic(nil, nil, nil)
		svc.Subscribe(docsTopicListener)
	})

	t.Run("Duplicated", func(t *testing.T) {
		svc, err := cms.NewDocsTopicService(fixDocsTopicInformer())
		require.NoError(t, err)

		docsTopicListener := listener.NewDocsTopic(nil, nil, nil)
		svc.Subscribe(docsTopicListener)
		svc.Subscribe(docsTopicListener)
	})

	t.Run("Multiple", func(t *testing.T) {
		svc, err := cms.NewDocsTopicService(fixDocsTopicInformer())
		require.NoError(t, err)

		docsTopicListenerA := listener.NewDocsTopic(nil, nil, nil)
		docsTopicListenerB := listener.NewDocsTopic(nil, nil, nil)

		svc.Subscribe(docsTopicListenerA)
		svc.Subscribe(docsTopicListenerB)
	})

	t.Run("Nil", func(t *testing.T) {
		svc, err := cms.NewDocsTopicService(fixDocsTopicInformer())
		require.NoError(t, err)

		svc.Subscribe(nil)
	})
}

func TestDocsTopicService_Unsubscribe(t *testing.T) {
	t.Run("Existing", func(t *testing.T) {
		svc, err := cms.NewDocsTopicService(fixDocsTopicInformer())
		require.NoError(t, err)

		docsTopicListener := listener.NewDocsTopic(nil, nil, nil)

		svc.Subscribe(docsTopicListener)
		svc.Unsubscribe(docsTopicListener)
	})

	t.Run("Duplicated", func(t *testing.T) {
		svc, err := cms.NewDocsTopicService(fixDocsTopicInformer())
		require.NoError(t, err)

		docsTopicListener := listener.NewDocsTopic(nil, nil, nil)
		svc.Subscribe(docsTopicListener)
		svc.Subscribe(docsTopicListener)

		svc.Unsubscribe(docsTopicListener)
	})

	t.Run("Multiple", func(t *testing.T) {
		svc, err := cms.NewDocsTopicService(fixDocsTopicInformer())
		require.NoError(t, err)

		docsTopicListenerA := listener.NewDocsTopic(nil, nil, nil)
		docsTopicListenerB := listener.NewDocsTopic(nil, nil, nil)

		svc.Subscribe(docsTopicListenerA)
		svc.Subscribe(docsTopicListenerB)

		svc.Unsubscribe(docsTopicListenerA)
	})

	t.Run("Nil", func(t *testing.T) {
		svc, err := cms.NewDocsTopicService(fixDocsTopicInformer())
		require.NoError(t, err)

		svc.Unsubscribe(nil)
	})
}

func fixUnstructuredDocsTopic(metadata map[string]interface{}) *unstructured.Unstructured {
	return testingUtils.NewUnstructured(v1alpha1.SchemeGroupVersion.String(), "DocsTopic", metadata, nil, nil)
}

func fixDocsTopic(name string, labels map[string]string) *v1alpha1.DocsTopic {
	return &v1alpha1.DocsTopic{
		TypeMeta: metav1.TypeMeta{
			Kind: "DocsTopic",
			APIVersion: v1alpha1.SchemeGroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: DocsTopicNamespace,
			Labels:    labels,
		},
	}
}

func fixDocsTopicInformer(objects ...runtime.Object) cache.SharedIndexInformer {
	fakeClient := fake.NewSimpleDynamicClient(runtime.NewScheme(), objects...)
	informerFactory := dynamicinformer.NewDynamicSharedInformerFactory(fakeClient, 0)

	informer := informerFactory.ForResource(schema.GroupVersionResource{
		Version:  v1alpha1.SchemeGroupVersion.Version,
		Group:    v1alpha1.SchemeGroupVersion.Group,
		Resource: "docstopics",
	}).Informer()

	return informer
}

