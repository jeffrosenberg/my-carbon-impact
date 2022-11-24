// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/jeffrosenberg/my-carbon-impact/internal/db (interfaces: Client,PutItemInputGenerator)

// Package mock_aws is a generated GoMock package.
package mock_aws

import (
	context "context"
	reflect "reflect"

	dynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	gomock "github.com/golang/mock/gomock"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// GetItem mocks base method.
func (m *MockClient) GetItem(arg0 context.Context, arg1 *dynamodb.GetItemInput, arg2 ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetItem", varargs...)
	ret0, _ := ret[0].(*dynamodb.GetItemOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetItem indicates an expected call of GetItem.
func (mr *MockClientMockRecorder) GetItem(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetItem", reflect.TypeOf((*MockClient)(nil).GetItem), varargs...)
}

// PutItem mocks base method.
func (m *MockClient) PutItem(arg0 context.Context, arg1 *dynamodb.PutItemInput, arg2 ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "PutItem", varargs...)
	ret0, _ := ret[0].(*dynamodb.PutItemOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PutItem indicates an expected call of PutItem.
func (mr *MockClientMockRecorder) PutItem(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutItem", reflect.TypeOf((*MockClient)(nil).PutItem), varargs...)
}

// Query mocks base method.
func (m *MockClient) Query(arg0 context.Context, arg1 *dynamodb.QueryInput, arg2 ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Query", varargs...)
	ret0, _ := ret[0].(*dynamodb.QueryOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Query indicates an expected call of Query.
func (mr *MockClientMockRecorder) Query(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockClient)(nil).Query), varargs...)
}

// Scan mocks base method.
func (m *MockClient) Scan(arg0 context.Context, arg1 *dynamodb.ScanInput, arg2 ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Scan", varargs...)
	ret0, _ := ret[0].(*dynamodb.ScanOutput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Scan indicates an expected call of Scan.
func (mr *MockClientMockRecorder) Scan(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Scan", reflect.TypeOf((*MockClient)(nil).Scan), varargs...)
}

// MockPutItemInputGenerator is a mock of PutItemInputGenerator interface.
type MockPutItemInputGenerator struct {
	ctrl     *gomock.Controller
	recorder *MockPutItemInputGeneratorMockRecorder
}

// MockPutItemInputGeneratorMockRecorder is the mock recorder for MockPutItemInputGenerator.
type MockPutItemInputGeneratorMockRecorder struct {
	mock *MockPutItemInputGenerator
}

// NewMockPutItemInputGenerator creates a new mock instance.
func NewMockPutItemInputGenerator(ctrl *gomock.Controller) *MockPutItemInputGenerator {
	mock := &MockPutItemInputGenerator{ctrl: ctrl}
	mock.recorder = &MockPutItemInputGeneratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPutItemInputGenerator) EXPECT() *MockPutItemInputGeneratorMockRecorder {
	return m.recorder
}

// GeneratePutItemInput mocks base method.
func (m *MockPutItemInputGenerator) GeneratePutItemInput() (*dynamodb.PutItemInput, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GeneratePutItemInput")
	ret0, _ := ret[0].(*dynamodb.PutItemInput)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GeneratePutItemInput indicates an expected call of GeneratePutItemInput.
func (mr *MockPutItemInputGeneratorMockRecorder) GeneratePutItemInput() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GeneratePutItemInput", reflect.TypeOf((*MockPutItemInputGenerator)(nil).GeneratePutItemInput))
}