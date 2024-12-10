// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/vmware/vsphere-automation-sdk-go/services/nsxt/orgs/projects/vpcs/subnets (interfaces: SubnetConnectionBindingMapsClient)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/vmware/vsphere-automation-sdk-go/services/nsxt/model"
)

// MockSubnetConnectionBindingMapsClient is a mock of SubnetConnectionBindingMapsClient interface.
type MockSubnetConnectionBindingMapsClient struct {
	ctrl     *gomock.Controller
	recorder *MockSubnetConnectionBindingMapsClientMockRecorder
}

// MockSubnetConnectionBindingMapsClientMockRecorder is the mock recorder for MockSubnetConnectionBindingMapsClient.
type MockSubnetConnectionBindingMapsClientMockRecorder struct {
	mock *MockSubnetConnectionBindingMapsClient
}

// NewMockSubnetConnectionBindingMapsClient creates a new mock instance.
func NewMockSubnetConnectionBindingMapsClient(ctrl *gomock.Controller) *MockSubnetConnectionBindingMapsClient {
	mock := &MockSubnetConnectionBindingMapsClient{ctrl: ctrl}
	mock.recorder = &MockSubnetConnectionBindingMapsClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSubnetConnectionBindingMapsClient) EXPECT() *MockSubnetConnectionBindingMapsClientMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockSubnetConnectionBindingMapsClient) Delete(arg0, arg1, arg2, arg3, arg4 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockSubnetConnectionBindingMapsClientMockRecorder) Delete(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockSubnetConnectionBindingMapsClient)(nil).Delete), arg0, arg1, arg2, arg3, arg4)
}

// Get mocks base method.
func (m *MockSubnetConnectionBindingMapsClient) Get(arg0, arg1, arg2, arg3, arg4 string) (model.SubnetConnectionBindingMap, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(model.SubnetConnectionBindingMap)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockSubnetConnectionBindingMapsClientMockRecorder) Get(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockSubnetConnectionBindingMapsClient)(nil).Get), arg0, arg1, arg2, arg3, arg4)
}

// List mocks base method.
func (m *MockSubnetConnectionBindingMapsClient) List(arg0, arg1, arg2, arg3 string, arg4 *string, arg5 *bool, arg6 *string, arg7 *int64, arg8 *bool, arg9 *string) (model.SubnetConnectionBindingMapListResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9)
	ret0, _ := ret[0].(model.SubnetConnectionBindingMapListResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockSubnetConnectionBindingMapsClientMockRecorder) List(arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockSubnetConnectionBindingMapsClient)(nil).List), arg0, arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8, arg9)
}

// Patch mocks base method.
func (m *MockSubnetConnectionBindingMapsClient) Patch(arg0, arg1, arg2, arg3, arg4 string, arg5 model.SubnetConnectionBindingMap) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Patch", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(error)
	return ret0
}

// Patch indicates an expected call of Patch.
func (mr *MockSubnetConnectionBindingMapsClientMockRecorder) Patch(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Patch", reflect.TypeOf((*MockSubnetConnectionBindingMapsClient)(nil).Patch), arg0, arg1, arg2, arg3, arg4, arg5)
}

// Update mocks base method.
func (m *MockSubnetConnectionBindingMapsClient) Update(arg0, arg1, arg2, arg3, arg4 string, arg5 model.SubnetConnectionBindingMap) (model.SubnetConnectionBindingMap, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(model.SubnetConnectionBindingMap)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockSubnetConnectionBindingMapsClientMockRecorder) Update(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockSubnetConnectionBindingMapsClient)(nil).Update), arg0, arg1, arg2, arg3, arg4, arg5)
}
