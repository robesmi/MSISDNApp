// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/robesmi/MSISDNApp/service (interfaces: MSISDNService)

// Package service is a generated GoMock package.
package service

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/robesmi/MSISDNApp/model"
	dto "github.com/robesmi/MSISDNApp/model/dto"
)

// MockMSISDNService is a mock of MSISDNService interface.
type MockMSISDNService struct {
	ctrl     *gomock.Controller
	recorder *MockMSISDNServiceMockRecorder
}

// MockMSISDNServiceMockRecorder is the mock recorder for MockMSISDNService.
type MockMSISDNServiceMockRecorder struct {
	mock *MockMSISDNService
}

// NewMockMSISDNService creates a new mock instance.
func NewMockMSISDNService(ctrl *gomock.Controller) *MockMSISDNService {
	mock := &MockMSISDNService{ctrl: ctrl}
	mock.recorder = &MockMSISDNServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMSISDNService) EXPECT() *MockMSISDNServiceMockRecorder {
	return m.recorder
}

// AddNewCountry mocks base method.
func (m *MockMSISDNService) AddNewCountry(arg0 *dto.CountryRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddNewCountry", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddNewCountry indicates an expected call of AddNewCountry.
func (mr *MockMSISDNServiceMockRecorder) AddNewCountry(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddNewCountry", reflect.TypeOf((*MockMSISDNService)(nil).AddNewCountry), arg0)
}

// AddNewMobileOperator mocks base method.
func (m *MockMSISDNService) AddNewMobileOperator(arg0 *dto.OperatorRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddNewMobileOperator", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddNewMobileOperator indicates an expected call of AddNewMobileOperator.
func (mr *MockMSISDNServiceMockRecorder) AddNewMobileOperator(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddNewMobileOperator", reflect.TypeOf((*MockMSISDNService)(nil).AddNewMobileOperator), arg0)
}

// GetAllCountries mocks base method.
func (m *MockMSISDNService) GetAllCountries() (*[]model.Country, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllCountries")
	ret0, _ := ret[0].(*[]model.Country)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllCountries indicates an expected call of GetAllCountries.
func (mr *MockMSISDNServiceMockRecorder) GetAllCountries() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllCountries", reflect.TypeOf((*MockMSISDNService)(nil).GetAllCountries))
}

// GetAllMobileOperators mocks base method.
func (m *MockMSISDNService) GetAllMobileOperators() (*[]model.MobileOperator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllMobileOperators")
	ret0, _ := ret[0].(*[]model.MobileOperator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllMobileOperators indicates an expected call of GetAllMobileOperators.
func (mr *MockMSISDNServiceMockRecorder) GetAllMobileOperators() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllMobileOperators", reflect.TypeOf((*MockMSISDNService)(nil).GetAllMobileOperators))
}

// LookupMSISDN mocks base method.
func (m *MockMSISDNService) LookupMSISDN(arg0 string) (*dto.NumberLookupResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LookupMSISDN", arg0)
	ret0, _ := ret[0].(*dto.NumberLookupResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LookupMSISDN indicates an expected call of LookupMSISDN.
func (mr *MockMSISDNServiceMockRecorder) LookupMSISDN(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LookupMSISDN", reflect.TypeOf((*MockMSISDNService)(nil).LookupMSISDN), arg0)
}

// RemoveCountry mocks base method.
func (m *MockMSISDNService) RemoveCountry(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveCountry", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveCountry indicates an expected call of RemoveCountry.
func (mr *MockMSISDNServiceMockRecorder) RemoveCountry(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveCountry", reflect.TypeOf((*MockMSISDNService)(nil).RemoveCountry), arg0)
}

// RemoveOperator mocks base method.
func (m *MockMSISDNService) RemoveOperator(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveOperator", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveOperator indicates an expected call of RemoveOperator.
func (mr *MockMSISDNServiceMockRecorder) RemoveOperator(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveOperator", reflect.TypeOf((*MockMSISDNService)(nil).RemoveOperator), arg0)
}
