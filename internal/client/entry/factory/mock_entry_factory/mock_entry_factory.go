// Code generated by MockGen. DO NOT EDIT.
// Source: entry_factory_interface.go

// Package mock_entry_factory is a generated GoMock package.
package mock_entry_factory

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	command "github.com/anoriar/gophkeeper/internal/client/entry/dto/command"
	entry_ext "github.com/anoriar/gophkeeper/internal/client/entry/dto/repository/entry_ext"
	entity "github.com/anoriar/gophkeeper/internal/client/entry/entity"
)

// MockEntryFactoryInterface is a mock of EntryFactoryInterface interface.
type MockEntryFactoryInterface struct {
	ctrl     *gomock.Controller
	recorder *MockEntryFactoryInterfaceMockRecorder
}

// MockEntryFactoryInterfaceMockRecorder is the mock recorder for MockEntryFactoryInterface.
type MockEntryFactoryInterfaceMockRecorder struct {
	mock *MockEntryFactoryInterface
}

// NewMockEntryFactoryInterface creates a new mock instance.
func NewMockEntryFactoryInterface(ctrl *gomock.Controller) *MockEntryFactoryInterface {
	mock := &MockEntryFactoryInterface{ctrl: ctrl}
	mock.recorder = &MockEntryFactoryInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEntryFactoryInterface) EXPECT() *MockEntryFactoryInterfaceMockRecorder {
	return m.recorder
}

// CreateFromAddCmd mocks base method.
func (m *MockEntryFactoryInterface) CreateFromAddCmd(command command.AddEntryCommand) (entity.Entry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateFromAddCmd", command)
	ret0, _ := ret[0].(entity.Entry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateFromAddCmd indicates an expected call of CreateFromAddCmd.
func (mr *MockEntryFactoryInterfaceMockRecorder) CreateFromAddCmd(command interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFromAddCmd", reflect.TypeOf((*MockEntryFactoryInterface)(nil).CreateFromAddCmd), command)
}

// CreateFromEditCmd mocks base method.
func (m *MockEntryFactoryInterface) CreateFromEditCmd(command command.EditEntryCommand) (entity.Entry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateFromEditCmd", command)
	ret0, _ := ret[0].(entity.Entry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateFromEditCmd indicates an expected call of CreateFromEditCmd.
func (mr *MockEntryFactoryInterfaceMockRecorder) CreateFromEditCmd(command interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFromEditCmd", reflect.TypeOf((*MockEntryFactoryInterface)(nil).CreateFromEditCmd), command)
}

// CreateFromSyncResponse mocks base method.
func (m *MockEntryFactoryInterface) CreateFromSyncResponse(syncResponse entry_ext.SyncResponse) ([]entity.Entry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateFromSyncResponse", syncResponse)
	ret0, _ := ret[0].([]entity.Entry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateFromSyncResponse indicates an expected call of CreateFromSyncResponse.
func (mr *MockEntryFactoryInterfaceMockRecorder) CreateFromSyncResponse(syncResponse interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFromSyncResponse", reflect.TypeOf((*MockEntryFactoryInterface)(nil).CreateFromSyncResponse), syncResponse)
}
