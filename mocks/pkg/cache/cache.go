// Code generated by MockGen. DO NOT EDIT.
// Source: pkg/cache/cache.go

// Package mock_cache is a generated GoMock package.
package mock_cache

import (
	context "context"
	reflect "reflect"
	time "time"

	cache "github.com/adanyl0v/pocket-ideas/pkg/cache"
	gomock "github.com/golang/mock/gomock"
)

// MockScanIterator is a mock of ScanIterator interface.
type MockScanIterator struct {
	ctrl     *gomock.Controller
	recorder *MockScanIteratorMockRecorder
}

// MockScanIteratorMockRecorder is the mock recorder for MockScanIterator.
type MockScanIteratorMockRecorder struct {
	mock *MockScanIterator
}

// NewMockScanIterator creates a new mock instance.
func NewMockScanIterator(ctrl *gomock.Controller) *MockScanIterator {
	mock := &MockScanIterator{ctrl: ctrl}
	mock.recorder = &MockScanIteratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockScanIterator) EXPECT() *MockScanIteratorMockRecorder {
	return m.recorder
}

// Err mocks base method.
func (m *MockScanIterator) Err() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Err")
	ret0, _ := ret[0].(error)
	return ret0
}

// Err indicates an expected call of Err.
func (mr *MockScanIteratorMockRecorder) Err() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Err", reflect.TypeOf((*MockScanIterator)(nil).Err))
}

// Next mocks base method.
func (m *MockScanIterator) Next(ctx context.Context) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Next", ctx)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Next indicates an expected call of Next.
func (mr *MockScanIteratorMockRecorder) Next(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Next", reflect.TypeOf((*MockScanIterator)(nil).Next), ctx)
}

// Val mocks base method.
func (m *MockScanIterator) Val() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Val")
	ret0, _ := ret[0].(string)
	return ret0
}

// Val indicates an expected call of Val.
func (mr *MockScanIteratorMockRecorder) Val() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Val", reflect.TypeOf((*MockScanIterator)(nil).Val))
}

// MockScanner is a mock of Scanner interface.
type MockScanner struct {
	ctrl     *gomock.Controller
	recorder *MockScannerMockRecorder
}

// MockScannerMockRecorder is the mock recorder for MockScanner.
type MockScannerMockRecorder struct {
	mock *MockScanner
}

// NewMockScanner creates a new mock instance.
func NewMockScanner(ctrl *gomock.Controller) *MockScanner {
	mock := &MockScanner{ctrl: ctrl}
	mock.recorder = &MockScannerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockScanner) EXPECT() *MockScannerMockRecorder {
	return m.recorder
}

// Scan mocks base method.
func (m *MockScanner) Scan(ctx context.Context) cache.ScanIterator {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Scan", ctx)
	ret0, _ := ret[0].(cache.ScanIterator)
	return ret0
}

// Scan indicates an expected call of Scan.
func (mr *MockScannerMockRecorder) Scan(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Scan", reflect.TypeOf((*MockScanner)(nil).Scan), ctx)
}

// MockCursorScanner is a mock of CursorScanner interface.
type MockCursorScanner struct {
	ctrl     *gomock.Controller
	recorder *MockCursorScannerMockRecorder
}

// MockCursorScannerMockRecorder is the mock recorder for MockCursorScanner.
type MockCursorScannerMockRecorder struct {
	mock *MockCursorScanner
}

// NewMockCursorScanner creates a new mock instance.
func NewMockCursorScanner(ctrl *gomock.Controller) *MockCursorScanner {
	mock := &MockCursorScanner{ctrl: ctrl}
	mock.recorder = &MockCursorScannerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCursorScanner) EXPECT() *MockCursorScannerMockRecorder {
	return m.recorder
}

// Scan mocks base method.
func (m *MockCursorScanner) Scan(ctx context.Context) cache.ScanIterator {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Scan", ctx)
	ret0, _ := ret[0].(cache.ScanIterator)
	return ret0
}

// Scan indicates an expected call of Scan.
func (mr *MockCursorScannerMockRecorder) Scan(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Scan", reflect.TypeOf((*MockCursorScanner)(nil).Scan), ctx)
}

// WithArgs mocks base method.
func (m *MockCursorScanner) WithArgs(cursor uint64, match string, count int64) cache.CursorScanner {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithArgs", cursor, match, count)
	ret0, _ := ret[0].(cache.CursorScanner)
	return ret0
}

// WithArgs indicates an expected call of WithArgs.
func (mr *MockCursorScannerMockRecorder) WithArgs(cursor, match, count interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithArgs", reflect.TypeOf((*MockCursorScanner)(nil).WithArgs), cursor, match, count)
}

// MockKeyCursorScanner is a mock of KeyCursorScanner interface.
type MockKeyCursorScanner struct {
	ctrl     *gomock.Controller
	recorder *MockKeyCursorScannerMockRecorder
}

// MockKeyCursorScannerMockRecorder is the mock recorder for MockKeyCursorScanner.
type MockKeyCursorScannerMockRecorder struct {
	mock *MockKeyCursorScanner
}

// NewMockKeyCursorScanner creates a new mock instance.
func NewMockKeyCursorScanner(ctrl *gomock.Controller) *MockKeyCursorScanner {
	mock := &MockKeyCursorScanner{ctrl: ctrl}
	mock.recorder = &MockKeyCursorScannerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockKeyCursorScanner) EXPECT() *MockKeyCursorScannerMockRecorder {
	return m.recorder
}

// Scan mocks base method.
func (m *MockKeyCursorScanner) Scan(ctx context.Context) cache.ScanIterator {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Scan", ctx)
	ret0, _ := ret[0].(cache.ScanIterator)
	return ret0
}

// Scan indicates an expected call of Scan.
func (mr *MockKeyCursorScannerMockRecorder) Scan(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Scan", reflect.TypeOf((*MockKeyCursorScanner)(nil).Scan), ctx)
}

// WithArgs mocks base method.
func (m *MockKeyCursorScanner) WithArgs(key string, cursor uint64, match string, count int64) cache.KeyCursorScanner {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithArgs", key, cursor, match, count)
	ret0, _ := ret[0].(cache.KeyCursorScanner)
	return ret0
}

// WithArgs indicates an expected call of WithArgs.
func (mr *MockKeyCursorScannerMockRecorder) WithArgs(key, cursor, match, count interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithArgs", reflect.TypeOf((*MockKeyCursorScanner)(nil).WithArgs), key, cursor, match, count)
}

// MockConn is a mock of Conn interface.
type MockConn struct {
	ctrl     *gomock.Controller
	recorder *MockConnMockRecorder
}

// MockConnMockRecorder is the mock recorder for MockConn.
type MockConnMockRecorder struct {
	mock *MockConn
}

// NewMockConn creates a new mock instance.
func NewMockConn(ctrl *gomock.Controller) *MockConn {
	mock := &MockConn{ctrl: ctrl}
	mock.recorder = &MockConnMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConn) EXPECT() *MockConnMockRecorder {
	return m.recorder
}

// Begin mocks base method.
func (m *MockConn) Begin(ctx context.Context) cache.Tx {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Begin", ctx)
	ret0, _ := ret[0].(cache.Tx)
	return ret0
}

// Begin indicates an expected call of Begin.
func (mr *MockConnMockRecorder) Begin(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Begin", reflect.TypeOf((*MockConn)(nil).Begin), ctx)
}

// Delete mocks base method.
func (m *MockConn) Delete(ctx context.Context, key string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, key)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockConnMockRecorder) Delete(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockConn)(nil).Delete), ctx, key)
}

// Exists mocks base method.
func (m *MockConn) Exists(ctx context.Context, keys ...string) (int64, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range keys {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Exists", varargs...)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exists indicates an expected call of Exists.
func (mr *MockConnMockRecorder) Exists(ctx interface{}, keys ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, keys...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exists", reflect.TypeOf((*MockConn)(nil).Exists), varargs...)
}

// Get mocks base method.
func (m *MockConn) Get(ctx context.Context, key string, dest any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, key, dest)
	ret0, _ := ret[0].(error)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockConnMockRecorder) Get(ctx, key, dest interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockConn)(nil).Get), ctx, key, dest)
}

// Scan mocks base method.
func (m *MockConn) Scan(ctx context.Context, scanner cache.Scanner) cache.ScanIterator {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Scan", ctx, scanner)
	ret0, _ := ret[0].(cache.ScanIterator)
	return ret0
}

// Scan indicates an expected call of Scan.
func (mr *MockConnMockRecorder) Scan(ctx, scanner interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Scan", reflect.TypeOf((*MockConn)(nil).Scan), ctx, scanner)
}

// Set mocks base method.
func (m *MockConn) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", ctx, key, value, expiration)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockConnMockRecorder) Set(ctx, key, value, expiration interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockConn)(nil).Set), ctx, key, value, expiration)
}

// MockTx is a mock of Tx interface.
type MockTx struct {
	ctrl     *gomock.Controller
	recorder *MockTxMockRecorder
}

// MockTxMockRecorder is the mock recorder for MockTx.
type MockTxMockRecorder struct {
	mock *MockTx
}

// NewMockTx creates a new mock instance.
func NewMockTx(ctrl *gomock.Controller) *MockTx {
	mock := &MockTx{ctrl: ctrl}
	mock.recorder = &MockTxMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTx) EXPECT() *MockTxMockRecorder {
	return m.recorder
}

// Begin mocks base method.
func (m *MockTx) Begin(ctx context.Context) cache.Tx {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Begin", ctx)
	ret0, _ := ret[0].(cache.Tx)
	return ret0
}

// Begin indicates an expected call of Begin.
func (mr *MockTxMockRecorder) Begin(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Begin", reflect.TypeOf((*MockTx)(nil).Begin), ctx)
}

// Delete mocks base method.
func (m *MockTx) Delete(ctx context.Context, key string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, key)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Delete indicates an expected call of Delete.
func (mr *MockTxMockRecorder) Delete(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockTx)(nil).Delete), ctx, key)
}

// Discard mocks base method.
func (m *MockTx) Discard(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Discard", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Discard indicates an expected call of Discard.
func (mr *MockTxMockRecorder) Discard(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Discard", reflect.TypeOf((*MockTx)(nil).Discard), ctx)
}

// Exec mocks base method.
func (m *MockTx) Exec(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Exec", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Exec indicates an expected call of Exec.
func (mr *MockTxMockRecorder) Exec(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exec", reflect.TypeOf((*MockTx)(nil).Exec), ctx)
}

// Exists mocks base method.
func (m *MockTx) Exists(ctx context.Context, keys ...string) (int64, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range keys {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Exists", varargs...)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Exists indicates an expected call of Exists.
func (mr *MockTxMockRecorder) Exists(ctx interface{}, keys ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, keys...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Exists", reflect.TypeOf((*MockTx)(nil).Exists), varargs...)
}

// Get mocks base method.
func (m *MockTx) Get(ctx context.Context, key string, dest any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, key, dest)
	ret0, _ := ret[0].(error)
	return ret0
}

// Get indicates an expected call of Get.
func (mr *MockTxMockRecorder) Get(ctx, key, dest interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockTx)(nil).Get), ctx, key, dest)
}

// Scan mocks base method.
func (m *MockTx) Scan(ctx context.Context, scanner cache.Scanner) cache.ScanIterator {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Scan", ctx, scanner)
	ret0, _ := ret[0].(cache.ScanIterator)
	return ret0
}

// Scan indicates an expected call of Scan.
func (mr *MockTxMockRecorder) Scan(ctx, scanner interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Scan", reflect.TypeOf((*MockTx)(nil).Scan), ctx, scanner)
}

// Set mocks base method.
func (m *MockTx) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", ctx, key, value, expiration)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockTxMockRecorder) Set(ctx, key, value, expiration interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockTx)(nil).Set), ctx, key, value, expiration)
}
