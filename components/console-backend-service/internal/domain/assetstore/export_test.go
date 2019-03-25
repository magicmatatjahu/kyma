package assetstore

import (
	"k8s.io/apimachinery/pkg/runtime"
	fakeDynamic "k8s.io/client-go/dynamic/fake"
	"k8s.io/client-go/tools/cache"
)

func NewClusterAssetResolver(clusterAssetSvc clusterAssetSvc) *clusterAssetResolver {
	return newClusterAssetResolver(clusterAssetSvc)
}

func (r *clusterAssetResolver) SetAssetConverter(converter gqlClusterAssetConverter) {
	r.clusterAssetConverter = converter
}

func (r *clusterAssetResolver) SetFileConverter(converter gqlFileConverter) {
	r.fileConverter = converter
}

func NewClusterAssetService(informer cache.SharedIndexInformer) (*clusterAssetService, error) {
	return newClusterAssetService(informer)
}

func NewAssetResolver(assetSvc assetSvc) *assetResolver {
	return newAssetResolver(assetSvc)
}

func (r *assetResolver) SetAssetConverter(converter gqlAssetConverter) {
	r.assetConverter = converter
}

func (r *assetResolver) SetFileConverter(converter gqlFileConverter) {
	r.fileConverter = converter
}

func NewAssetService(informer cache.SharedIndexInformer) (*assetService, error) {
	return newAssetService(informer)
}

func (r *PluggableContainer) SetFakeClient() {
	r.cfg.dynamicClient = fakeDynamic.NewSimpleDynamicClient(runtime.NewScheme())
}
