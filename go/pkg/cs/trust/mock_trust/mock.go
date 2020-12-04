// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/scionproto/scion/go/pkg/cs/trust (interfaces: CACertProvider,PolicyGen,SignerGen)

// Package mock_trust is a generated GoMock package.
package mock_trust

import (
	context "context"
	x509 "crypto/x509"
	gomock "github.com/golang/mock/gomock"
	cppki "github.com/scionproto/scion/go/lib/scrypto/cppki"
	trust "github.com/scionproto/scion/go/pkg/trust"
	reflect "reflect"
)

// MockCACertProvider is a mock of CACertProvider interface
type MockCACertProvider struct {
	ctrl     *gomock.Controller
	recorder *MockCACertProviderMockRecorder
}

// MockCACertProviderMockRecorder is the mock recorder for MockCACertProvider
type MockCACertProviderMockRecorder struct {
	mock *MockCACertProvider
}

// NewMockCACertProvider creates a new mock instance
func NewMockCACertProvider(ctrl *gomock.Controller) *MockCACertProvider {
	mock := &MockCACertProvider{ctrl: ctrl}
	mock.recorder = &MockCACertProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCACertProvider) EXPECT() *MockCACertProviderMockRecorder {
	return m.recorder
}

// CACerts mocks base method
func (m *MockCACertProvider) CACerts(arg0 context.Context) ([]*x509.Certificate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CACerts", arg0)
	ret0, _ := ret[0].([]*x509.Certificate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CACerts indicates an expected call of CACerts
func (mr *MockCACertProviderMockRecorder) CACerts(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CACerts", reflect.TypeOf((*MockCACertProvider)(nil).CACerts), arg0)
}

// MockPolicyGen is a mock of PolicyGen interface
type MockPolicyGen struct {
	ctrl     *gomock.Controller
	recorder *MockPolicyGenMockRecorder
}

// MockPolicyGenMockRecorder is the mock recorder for MockPolicyGen
type MockPolicyGenMockRecorder struct {
	mock *MockPolicyGen
}

// NewMockPolicyGen creates a new mock instance
func NewMockPolicyGen(ctrl *gomock.Controller) *MockPolicyGen {
	mock := &MockPolicyGen{ctrl: ctrl}
	mock.recorder = &MockPolicyGenMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPolicyGen) EXPECT() *MockPolicyGenMockRecorder {
	return m.recorder
}

// Generate mocks base method
func (m *MockPolicyGen) Generate(arg0 context.Context) (cppki.CAPolicy, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Generate", arg0)
	ret0, _ := ret[0].(cppki.CAPolicy)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Generate indicates an expected call of Generate
func (mr *MockPolicyGenMockRecorder) Generate(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Generate", reflect.TypeOf((*MockPolicyGen)(nil).Generate), arg0)
}

// MockSignerGen is a mock of SignerGen interface
type MockSignerGen struct {
	ctrl     *gomock.Controller
	recorder *MockSignerGenMockRecorder
}

// MockSignerGenMockRecorder is the mock recorder for MockSignerGen
type MockSignerGenMockRecorder struct {
	mock *MockSignerGen
}

// NewMockSignerGen creates a new mock instance
func NewMockSignerGen(ctrl *gomock.Controller) *MockSignerGen {
	mock := &MockSignerGen{ctrl: ctrl}
	mock.recorder = &MockSignerGenMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockSignerGen) EXPECT() *MockSignerGenMockRecorder {
	return m.recorder
}

// Generate mocks base method
func (m *MockSignerGen) Generate(arg0 context.Context) (trust.Signer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Generate", arg0)
	ret0, _ := ret[0].(trust.Signer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Generate indicates an expected call of Generate
func (mr *MockSignerGenMockRecorder) Generate(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Generate", reflect.TypeOf((*MockSignerGen)(nil).Generate), arg0)
}
