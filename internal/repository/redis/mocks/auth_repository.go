// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/redis/auth_repository.go

// Package mock_redis is a generated GoMock package.
package mock_redis

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockJSONer is a mock of JSONer interface.
type MockJSONer struct {
	ctrl     *gomock.Controller
	recorder *MockJSONerMockRecorder
}

// MockJSONerMockRecorder is the mock recorder for MockJSONer.
type MockJSONerMockRecorder struct {
	mock *MockJSONer
}

// NewMockJSONer creates a new mock instance.
func NewMockJSONer(ctrl *gomock.Controller) *MockJSONer {
	mock := &MockJSONer{ctrl: ctrl}
	mock.recorder = &MockJSONerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockJSONer) EXPECT() *MockJSONerMockRecorder {
	return m.recorder
}

// Marshal mocks base method.
func (m *MockJSONer) Marshal(v interface{}) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Marshal", v)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Marshal indicates an expected call of Marshal.
func (mr *MockJSONerMockRecorder) Marshal(v interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Marshal", reflect.TypeOf((*MockJSONer)(nil).Marshal), v)
}

// Unmarshal mocks base method.
func (m *MockJSONer) Unmarshal(data []byte, v interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unmarshal", data, v)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unmarshal indicates an expected call of Unmarshal.
func (mr *MockJSONerMockRecorder) Unmarshal(data, v interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unmarshal", reflect.TypeOf((*MockJSONer)(nil).Unmarshal), data, v)
}
