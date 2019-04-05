// Code generated by mockery v1.0.0. DO NOT EDIT.
package automock

import mock "github.com/stretchr/testify/mock"
import shared "github.com/kyma-project/kyma/components/console-backend-service/internal/domain/shared"

// AssetStoreRetriever is an autogenerated mock type for the AssetStoreRetriever type
type AssetStoreRetriever struct {
	mock.Mock
}

// Asset provides a mock function with given fields:
func (_m *AssetStoreRetriever) Asset() shared.AssetGetter {
	ret := _m.Called()

	var r0 shared.AssetGetter
	if rf, ok := ret.Get(0).(func() shared.AssetGetter); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(shared.AssetGetter)
		}
	}

	return r0
}

// ClusterAsset provides a mock function with given fields:
func (_m *AssetStoreRetriever) ClusterAsset() shared.ClusterAssetGetter {
	ret := _m.Called()

	var r0 shared.ClusterAssetGetter
	if rf, ok := ret.Get(0).(func() shared.ClusterAssetGetter); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(shared.ClusterAssetGetter)
		}
	}

	return r0
}
