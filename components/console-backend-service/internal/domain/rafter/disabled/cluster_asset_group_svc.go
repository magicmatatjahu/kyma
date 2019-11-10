// Code generated by failery v1.0.0. DO NOT EDIT.

package disabled

import resource "github.com/kyma-project/kyma/components/console-backend-service/pkg/resource"
import v1beta1 "github.com/kyma-project/rafter/pkg/apis/rafter/v1beta1"

// clusterAssetGroupSvc is an autogenerated failing mock type for the clusterAssetGroupSvc type
type clusterAssetGroupSvc struct {
	err error
}

// NewClusterAssetGroupSvc creates a new clusterAssetGroupSvc type instance
func NewClusterAssetGroupSvc(err error) *clusterAssetGroupSvc {
	return &clusterAssetGroupSvc{err: err}
}

// Find provides a failing mock function with given fields: name
func (_m *clusterAssetGroupSvc) Find(name string) (*v1beta1.ClusterAssetGroup, error) {
	var r0 *v1beta1.ClusterAssetGroup
	var r1 error
	r1 = _m.err

	return r0, r1
}

// List provides a failing mock function with given fields: viewContext, groupName
func (_m *clusterAssetGroupSvc) List(viewContext *string, groupName *string) ([]*v1beta1.ClusterAssetGroup, error) {
	var r0 []*v1beta1.ClusterAssetGroup
	var r1 error
	r1 = _m.err

	return r0, r1
}

// Subscribe provides a failing mock function with given fields: listener
func (_m *clusterAssetGroupSvc) Subscribe(listener resource.Listener) {
}

// Unsubscribe provides a failing mock function with given fields: listener
func (_m *clusterAssetGroupSvc) Unsubscribe(listener resource.Listener) {
}
