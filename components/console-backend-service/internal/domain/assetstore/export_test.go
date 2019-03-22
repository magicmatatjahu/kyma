package assetstore

import (
	"k8s.io/apimachinery/pkg/runtime"
	fakeDynamic "k8s.io/client-go/dynamic/fake"
)

func (r *PluggableContainer) SetFakeClient() {
	r.cfg.dynamicClient = fakeDynamic.NewSimpleDynamicClient(runtime.NewScheme())
}
