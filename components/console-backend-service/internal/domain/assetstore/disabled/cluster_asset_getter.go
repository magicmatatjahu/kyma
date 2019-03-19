// Code generated by failery v1.0.0. DO NOT EDIT.

package disabled

import v1alpha2 "github.com/kyma-project/kyma/components/asset-store-controller-manager/pkg/apis/assetstore/v1alpha2"

// clusterAssetGetter is an autogenerated failing mock type for the clusterAssetGetter type
type clusterAssetGetter struct {
	err error
}

// NewClusterAssetGetter creates a new clusterAssetGetter type instance
func NewClusterAssetGetter(err error) *clusterAssetGetter {
	return &clusterAssetGetter{err: err}
}

// List provides a failing mock function with given fields: groupName
func (_m *clusterAssetGetter) List(groupName string) ([]*v1alpha2.ClusterAsset, error) {
	var r0 []*v1alpha2.ClusterAsset
	var r1 error
	r1 = _m.err

	return r0, r1
}
