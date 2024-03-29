// Code generated by MockGen. DO NOT EDIT.
// Source: db_transaction_interface.go

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockDBTransactionInterface is a mock of DBTransactionInterface interface.
type MockDBTransactionInterface struct {
	ctrl     *gomock.Controller
	recorder *MockDBTransactionInterfaceMockRecorder
}

// MockDBTransactionInterfaceMockRecorder is the mock recorder for MockDBTransactionInterface.
type MockDBTransactionInterfaceMockRecorder struct {
	mock *MockDBTransactionInterface
}

// NewMockDBTransactionInterface creates a new mock instance.
func NewMockDBTransactionInterface(ctrl *gomock.Controller) *MockDBTransactionInterface {
	mock := &MockDBTransactionInterface{ctrl: ctrl}
	mock.recorder = &MockDBTransactionInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDBTransactionInterface) EXPECT() *MockDBTransactionInterfaceMockRecorder {
	return m.recorder
}

// Commit mocks base method.
func (m *MockDBTransactionInterface) Commit() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Commit")
	ret0, _ := ret[0].(error)
	return ret0
}

// Commit indicates an expected call of Commit.
func (mr *MockDBTransactionInterfaceMockRecorder) Commit() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Commit", reflect.TypeOf((*MockDBTransactionInterface)(nil).Commit))
}

// GetTransaction mocks base method.
func (m *MockDBTransactionInterface) GetTransaction() interface{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransaction")
	ret0, _ := ret[0].(interface{})
	return ret0
}

// GetTransaction indicates an expected call of GetTransaction.
func (mr *MockDBTransactionInterfaceMockRecorder) GetTransaction() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransaction", reflect.TypeOf((*MockDBTransactionInterface)(nil).GetTransaction))
}

// Rollback mocks base method.
func (m *MockDBTransactionInterface) Rollback() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Rollback")
	ret0, _ := ret[0].(error)
	return ret0
}

// Rollback indicates an expected call of Rollback.
func (mr *MockDBTransactionInterfaceMockRecorder) Rollback() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Rollback", reflect.TypeOf((*MockDBTransactionInterface)(nil).Rollback))
}
