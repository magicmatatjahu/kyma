// Code generated by failery v1.0.0. DO NOT EDIT.

package disabled

import resource "github.com/kyma-project/kyma/components/console-backend-service/pkg/resource"
import v1alpha1 "github.com/kyma-project/kyma/components/cms-controller-manager/pkg/apis/cms/v1alpha1"

// clusterDocsTopicSvc is an autogenerated failing mock type for the clusterDocsTopicSvc type
type clusterDocsTopicSvc struct {
	err error
}

// NewClusterDocsTopicSvc creates a new clusterDocsTopicSvc type instance
func NewClusterDocsTopicSvc(err error) *clusterDocsTopicSvc {
	return &clusterDocsTopicSvc{err: err}
}

// List provides a failing mock function with given fields: groupName
func (_m *clusterDocsTopicSvc) List(groupName string) ([]*v1alpha1.ClusterDocsTopic, error) {
	var r0 []*v1alpha1.ClusterDocsTopic
	var r1 error
	r1 = _m.err

	return r0, r1
}

// Subscribe provides a failing mock function with given fields: listener
func (_m *clusterDocsTopicSvc) Subscribe(listener resource.Listener) {
}

// Unsubscribe provides a failing mock function with given fields: listener
func (_m *clusterDocsTopicSvc) Unsubscribe(listener resource.Listener) {
}
