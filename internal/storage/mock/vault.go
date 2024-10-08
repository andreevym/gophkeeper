// Code generated by MockGen. DO NOT EDIT.
// Source: vault.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	storage "github.com/andreevym/gophkeeper/internal/storage"
	gomock "github.com/golang/mock/gomock"
)

// MockVaultStorage is a mock of VaultStorage interface.
type MockVaultStorage struct {
	ctrl     *gomock.Controller
	recorder *MockVaultStorageMockRecorder
}

// MockVaultStorageMockRecorder is the mock recorder for MockVaultStorage.
type MockVaultStorageMockRecorder struct {
	mock *MockVaultStorage
}

// NewMockVaultStorage creates a new mock instance.
func NewMockVaultStorage(ctrl *gomock.Controller) *MockVaultStorage {
	mock := &MockVaultStorage{ctrl: ctrl}
	mock.recorder = &MockVaultStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVaultStorage) EXPECT() *MockVaultStorageMockRecorder {
	return m.recorder
}

// CreateVault mocks base method.
func (m *MockVaultStorage) CreateVault(ctx context.Context, v storage.Vault) (storage.Vault, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateVault", ctx, v)
	ret0, _ := ret[0].(storage.Vault)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateVault indicates an expected call of CreateVault.
func (mr *MockVaultStorageMockRecorder) CreateVault(ctx, v interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateVault", reflect.TypeOf((*MockVaultStorage)(nil).CreateVault), ctx, v)
}

// DeleteVault mocks base method.
func (m *MockVaultStorage) DeleteVault(ctx context.Context, id uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteVault", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteVault indicates an expected call of DeleteVault.
func (mr *MockVaultStorageMockRecorder) DeleteVault(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteVault", reflect.TypeOf((*MockVaultStorage)(nil).DeleteVault), ctx, id)
}

// GetVault mocks base method.
func (m *MockVaultStorage) GetVault(ctx context.Context, id uint64) (storage.Vault, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVault", ctx, id)
	ret0, _ := ret[0].(storage.Vault)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVault indicates an expected call of GetVault.
func (mr *MockVaultStorageMockRecorder) GetVault(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVault", reflect.TypeOf((*MockVaultStorage)(nil).GetVault), ctx, id)
}

// UpdateVault mocks base method.
func (m *MockVaultStorage) UpdateVault(ctx context.Context, v storage.Vault) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateVault", ctx, v)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateVault indicates an expected call of UpdateVault.
func (mr *MockVaultStorageMockRecorder) UpdateVault(ctx, v interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateVault", reflect.TypeOf((*MockVaultStorage)(nil).UpdateVault), ctx, v)
}
