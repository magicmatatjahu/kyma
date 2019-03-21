// Code generated by failery v1.0.0. DO NOT EDIT.

package disabled

import v1alpha2 "github.com/kyma-project/kyma/components/asset-store-controller-manager/pkg/apis/assetstore/v1alpha2"

// clusterAssetSvc is an autogenerated failing mock type for the clusterAssetSvc type
type clusterAssetSvc struct {
	err error
}

// NewClusterAssetSvc creates a new clusterAssetSvc type instance
func NewClusterAssetSvc(err error) *clusterAssetSvc {
	return &clusterAssetSvc{err: err}
}

// Find provides a failing mock function with given fields: name
func (_m *clusterAssetSvc) Find(name string) (*v1alpha2.ClusterAsset, error) {
	var r0 *v1alpha2.ClusterAsset
	var r1 error
	r1 = _m.err

	return r0, r1
}

// List provides a failing mock function with given fields: groupName
func (_m *clusterAssetSvc) List(groupName string) ([]*v1alpha2.ClusterAsset, error) {
	var r0 []*v1alpha2.ClusterAsset
	var r1 error
	r1 = _m.err

	return r0, r1
}

// ListForDocsTopicByType provides a failing mock function with given fields: docsTopicName, types
func (_m *clusterAssetSvc) ListForDocsTopicByType(docsTopicName string, types []string) ([]*v1alpha2.ClusterAsset, error) {
	var r0 []*v1alpha2.ClusterAsset
	var r1 error
	r1 = _m.err

	return r0, r1
}
