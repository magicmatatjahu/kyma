// Code generated by failery v1.0.0. DO NOT EDIT.

package disabled

import gqlschema "github.com/kyma-project/kyma/components/console-backend-service/internal/gqlschema"

import v1beta1 "github.com/kyma-project/rafter/pkg/apis/rafter/v1beta1"

// gqlAssetConverter is an autogenerated failing mock type for the gqlAssetConverter type
type gqlAssetConverter struct {
	err error
}

// NewGqlAssetConverter creates a new gqlAssetConverter type instance
func NewGqlAssetConverter(err error) *gqlAssetConverter {
	return &gqlAssetConverter{err: err}
}

// ToGQL provides a failing mock function with given fields: in
func (_m *gqlAssetConverter) ToGQL(in *v1beta1.Asset) (*gqlschema.Asset, error) {
	var r0 *gqlschema.Asset
	var r1 error
	r1 = _m.err

	return r0, r1
}

// ToGQLs provides a failing mock function with given fields: in
func (_m *gqlAssetConverter) ToGQLs(in []*v1beta1.Asset) ([]gqlschema.Asset, error) {
	var r0 []gqlschema.Asset
	var r1 error
	r1 = _m.err

	return r0, r1
}
