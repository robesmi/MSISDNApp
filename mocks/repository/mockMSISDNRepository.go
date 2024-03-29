// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/robesmi/MSISDNApp/repository (interfaces: MSISDNRepository)

// Package repository is a generated GoMock package.
package repository

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/robesmi/MSISDNApp/model"
	dto "github.com/robesmi/MSISDNApp/model/dto"
)

// MockMSISDNRepository is a mock of MSISDNRepository interface.
type MockMSISDNRepository struct {
	ctrl     *gomock.Controller
	recorder *MockMSISDNRepositoryMockRecorder
}

// MockMSISDNRepositoryMockRecorder is the mock recorder for MockMSISDNRepository.
type MockMSISDNRepositoryMockRecorder struct {
	mock *MockMSISDNRepository
}

// NewMockMSISDNRepository creates a new mock instance.
func NewMockMSISDNRepository(ctrl *gomock.Controller) *MockMSISDNRepository {
	mock := &MockMSISDNRepository{ctrl: ctrl}
	mock.recorder = &MockMSISDNRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMSISDNRepository) EXPECT() *MockMSISDNRepositoryMockRecorder {
	return m.recorder
}

// AddNewCountry mocks base method.
func (m *MockMSISDNRepository) AddNewCountry(arg0, arg1, arg2 string, arg3 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddNewCountry", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddNewCountry indicates an expected call of AddNewCountry.
func (mr *MockMSISDNRepositoryMockRecorder) AddNewCountry(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddNewCountry", reflect.TypeOf((*MockMSISDNRepository)(nil).AddNewCountry), arg0, arg1, arg2, arg3)
}

// AddNewMobileOperator mocks base method.
func (m *MockMSISDNRepository) AddNewMobileOperator(arg0, arg1, arg2 string, arg3 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddNewMobileOperator", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddNewMobileOperator indicates an expected call of AddNewMobileOperator.
func (mr *MockMSISDNRepositoryMockRecorder) AddNewMobileOperator(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddNewMobileOperator", reflect.TypeOf((*MockMSISDNRepository)(nil).AddNewMobileOperator), arg0, arg1, arg2, arg3)
}

// GetAllCountries mocks base method.
func (m *MockMSISDNRepository) GetAllCountries() (*[]model.Country, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllCountries")
	ret0, _ := ret[0].(*[]model.Country)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllCountries indicates an expected call of GetAllCountries.
func (mr *MockMSISDNRepositoryMockRecorder) GetAllCountries() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllCountries", reflect.TypeOf((*MockMSISDNRepository)(nil).GetAllCountries))
}

// GetAllMobileOperators mocks base method.
func (m *MockMSISDNRepository) GetAllMobileOperators() (*[]model.MobileOperator, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllMobileOperators")
	ret0, _ := ret[0].(*[]model.MobileOperator)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllMobileOperators indicates an expected call of GetAllMobileOperators.
func (mr *MockMSISDNRepositoryMockRecorder) GetAllMobileOperators() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllMobileOperators", reflect.TypeOf((*MockMSISDNRepository)(nil).GetAllMobileOperators))
}

// LookupCountryCode mocks base method.
func (m *MockMSISDNRepository) LookupCountryCode(arg0 string) (*dto.CountryLookupResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LookupCountryCode", arg0)
	ret0, _ := ret[0].(*dto.CountryLookupResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LookupCountryCode indicates an expected call of LookupCountryCode.
func (mr *MockMSISDNRepositoryMockRecorder) LookupCountryCode(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LookupCountryCode", reflect.TypeOf((*MockMSISDNRepository)(nil).LookupCountryCode), arg0)
}

// LookupMobileOperator mocks base method.
func (m *MockMSISDNRepository) LookupMobileOperator(arg0, arg1 string) (*dto.MobileOperatorLookupResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LookupMobileOperator", arg0, arg1)
	ret0, _ := ret[0].(*dto.MobileOperatorLookupResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LookupMobileOperator indicates an expected call of LookupMobileOperator.
func (mr *MockMSISDNRepositoryMockRecorder) LookupMobileOperator(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LookupMobileOperator", reflect.TypeOf((*MockMSISDNRepository)(nil).LookupMobileOperator), arg0, arg1)
}

// RemoveCountry mocks base method.
func (m *MockMSISDNRepository) RemoveCountry(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveCountry", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveCountry indicates an expected call of RemoveCountry.
func (mr *MockMSISDNRepositoryMockRecorder) RemoveCountry(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveCountry", reflect.TypeOf((*MockMSISDNRepository)(nil).RemoveCountry), arg0)
}

// RemoveOperator mocks base method.
func (m *MockMSISDNRepository) RemoveOperator(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveOperator", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveOperator indicates an expected call of RemoveOperator.
func (mr *MockMSISDNRepositoryMockRecorder) RemoveOperator(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveOperator", reflect.TypeOf((*MockMSISDNRepository)(nil).RemoveOperator), arg0)
}
