// Code generated by failery v1.0.0. DO NOT EDIT.

package disabled

import v1alpha1 "github.com/kyma-project/kyma/components/cms-controller-manager/pkg/apis/cms/v1alpha1"

// docsTopicGetter is an autogenerated failing mock type for the docsTopicGetter type
type docsTopicGetter struct {
	err error
}

// NewDocsTopicGetter creates a new docsTopicGetter type instance
func NewDocsTopicGetter(err error) *docsTopicGetter {
	return &docsTopicGetter{err: err}
}

// List provides a failing mock function with given fields: namespace, groupName
func (_m *docsTopicGetter) List(namespace string, groupName string) ([]*v1alpha1.DocsTopic, error) {
	var r0 []*v1alpha1.DocsTopic
	var r1 error
	r1 = _m.err

	return r0, r1
}
