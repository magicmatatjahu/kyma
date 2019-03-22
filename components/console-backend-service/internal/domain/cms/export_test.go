package cms

import (
	"github.com/kyma-project/kyma/components/console-backend-service/internal/domain/shared"
	"k8s.io/client-go/tools/cache"
	"k8s.io/apimachinery/pkg/runtime"
	fakeDynamic "k8s.io/client-go/dynamic/fake"
)

func NewClusterDocsTopicResolver(clusterDocsTopicService *clusterDocsTopicService, assetStoreRetriever shared.AssetStoreRetriever) *clusterDocsTopicResolver {
	return newClusterDocsTopicResolver(clusterDocsTopicService, assetStoreRetriever)
}

func NewClusterDocsTopicService(informer cache.SharedIndexInformer) (*clusterDocsTopicService, error) {
	return newClusterDocsTopicService(informer)
}

func NewDocsTopicResolver(docsTopicService *docsTopicService, assetStoreRetriever shared.AssetStoreRetriever) *docsTopicResolver {
	return newDocsTopicResolver(docsTopicService, assetStoreRetriever)
}

func NewDocsTopicService(informer cache.SharedIndexInformer) (*docsTopicService, error) {
	return newDocsTopicService(informer)
}

func (r *PluggableContainer) SetFakeClient() {
	r.cfg.dynamicClient = fakeDynamic.NewSimpleDynamicClient(runtime.NewScheme())
}
