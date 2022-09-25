// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	models "github.com/todanni/auth/models"
)

// DashboardStorage is an autogenerated mock type for the DashboardStorage type
type DashboardStorage struct {
	mock.Mock
}

// Create provides a mock function with given fields: owner, invited
func (_m *DashboardStorage) Create(owner uint, invited uint) (models.Dashboard, error) {
	ret := _m.Called(owner, invited)

	var r0 models.Dashboard
	if rf, ok := ret.Get(0).(func(uint, uint) models.Dashboard); ok {
		r0 = rf(owner, invited)
	} else {
		r0 = ret.Get(0).(models.Dashboard)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint, uint) error); ok {
		r1 = rf(owner, invited)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: id
func (_m *DashboardStorage) Delete(id string) error {
	ret := _m.Called(id)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetById provides a mock function with given fields: id
func (_m *DashboardStorage) GetById(id string) (models.Dashboard, error) {
	ret := _m.Called(id)

	var r0 models.Dashboard
	if rf, ok := ret.Get(0).(func(string) models.Dashboard); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(models.Dashboard)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// List provides a mock function with given fields: userid
func (_m *DashboardStorage) List(userid uint) ([]models.Dashboard, error) {
	ret := _m.Called(userid)

	var r0 []models.Dashboard
	if rf, ok := ret.Get(0).(func(uint) []models.Dashboard); ok {
		r0 = rf(userid)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Dashboard)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(userid)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateStatus provides a mock function with given fields: id, status
func (_m *DashboardStorage) UpdateStatus(id string, status models.Status) (models.Dashboard, error) {
	ret := _m.Called(id, status)

	var r0 models.Dashboard
	if rf, ok := ret.Get(0).(func(string, models.Status) models.Dashboard); ok {
		r0 = rf(id, status)
	} else {
		r0 = ret.Get(0).(models.Dashboard)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, models.Status) error); ok {
		r1 = rf(id, status)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewDashboardStorage interface {
	mock.TestingT
	Cleanup(func())
}

// NewDashboardStorage creates a new instance of DashboardStorage. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewDashboardStorage(t mockConstructorTestingTNewDashboardStorage) *DashboardStorage {
	mock := &DashboardStorage{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}