// Code generated by mockery v1.0.0. DO NOT EDIT.
package automock

import mock "github.com/stretchr/testify/mock"
import v1alpha1 "github.com/kyma-project/kyma/components/cms-controller-manager/pkg/apis/cms/v1alpha1"

// docsTopicGetter is an autogenerated mock type for the docsTopicGetter type
type docsTopicGetter struct {
	mock.Mock
}

// List provides a mock function with given fields: namespace, groupName
func (_m *docsTopicGetter) List(namespace string, groupName string) ([]*v1alpha1.DocsTopic, error) {
	ret := _m.Called(namespace, groupName)

	var r0 []*v1alpha1.DocsTopic
	if rf, ok := ret.Get(0).(func(string, string) []*v1alpha1.DocsTopic); ok {
		r0 = rf(namespace, groupName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*v1alpha1.DocsTopic)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(namespace, groupName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
