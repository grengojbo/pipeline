// Code generated by mockery v1.0.0. DO NOT EDIT.

package cluster

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockStore is an autogenerated mock type for the Store type
type MockStore struct {
	mock.Mock
}

// GetCluster provides a mock function with given fields: ctx, id
func (_m *MockStore) GetCluster(ctx context.Context, id uint) (Cluster, error) {
	ret := _m.Called(ctx, id)

	var r0 Cluster
	if rf, ok := ret.Get(0).(func(context.Context, uint) Cluster); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(Cluster)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetClusterByName provides a mock function with given fields: ctx, orgID, clusterName
func (_m *MockStore) GetClusterByName(ctx context.Context, orgID uint, clusterName string) (Cluster, error) {
	ret := _m.Called(ctx, orgID, clusterName)

	var r0 Cluster
	if rf, ok := ret.Get(0).(func(context.Context, uint, string) Cluster); ok {
		r0 = rf(ctx, orgID, clusterName)
	} else {
		r0 = ret.Get(0).(Cluster)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint, string) error); ok {
		r1 = rf(ctx, orgID, clusterName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SetStatus provides a mock function with given fields: ctx, id, status, statusMessage
func (_m *MockStore) SetStatus(ctx context.Context, id uint, status string, statusMessage string) error {
	ret := _m.Called(ctx, id, status, statusMessage)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint, string, string) error); ok {
		r0 = rf(ctx, id, status, statusMessage)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
