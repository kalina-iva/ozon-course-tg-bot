// Code generated by MockGen. DO NOT EDIT.
// Source: internal/model/messages/incoming_msg.go

// Package mock_messages is a generated GoMock package.
package mock_messages

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	entity "gitlab.ozon.dev/mary.kalina/telegram-bot/internal/model/messages/entity"
)

// MockexchangeRateRepository is a mock of exchangeRateRepository interface.
type MockexchangeRateRepository struct {
	ctrl     *gomock.Controller
	recorder *MockexchangeRateRepositoryMockRecorder
}

// MockexchangeRateRepositoryMockRecorder is the mock recorder for MockexchangeRateRepository.
type MockexchangeRateRepositoryMockRecorder struct {
	mock *MockexchangeRateRepository
}

// NewMockexchangeRateRepository creates a new mock instance.
func NewMockexchangeRateRepository(ctrl *gomock.Controller) *MockexchangeRateRepository {
	mock := &MockexchangeRateRepository{ctrl: ctrl}
	mock.recorder = &MockexchangeRateRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockexchangeRateRepository) EXPECT() *MockexchangeRateRepositoryMockRecorder {
	return m.recorder
}

// GetRate mocks base method.
func (m *MockexchangeRateRepository) GetRate(code string) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRate", code)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRate indicates an expected call of GetRate.
func (mr *MockexchangeRateRepositoryMockRecorder) GetRate(code interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRate", reflect.TypeOf((*MockexchangeRateRepository)(nil).GetRate), code)
}

// MockexpenseRepository is a mock of expenseRepository interface.
type MockexpenseRepository struct {
	ctrl     *gomock.Controller
	recorder *MockexpenseRepositoryMockRecorder
}

// MockexpenseRepositoryMockRecorder is the mock recorder for MockexpenseRepository.
type MockexpenseRepositoryMockRecorder struct {
	mock *MockexpenseRepository
}

// NewMockexpenseRepository creates a new mock instance.
func NewMockexpenseRepository(ctrl *gomock.Controller) *MockexpenseRepository {
	mock := &MockexpenseRepository{ctrl: ctrl}
	mock.recorder = &MockexpenseRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockexpenseRepository) EXPECT() *MockexpenseRepositoryMockRecorder {
	return m.recorder
}

// GetExpenses mocks base method.
func (m *MockexpenseRepository) GetExpenses(userID int64, period time.Time) []*entity.Expense {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetExpenses", userID, period)
	ret0, _ := ret[0].([]*entity.Expense)
	return ret0
}

// GetExpenses indicates an expected call of GetExpenses.
func (mr *MockexpenseRepositoryMockRecorder) GetExpenses(userID, period interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetExpenses", reflect.TypeOf((*MockexpenseRepository)(nil).GetExpenses), userID, period)
}

// New mocks base method.
func (m *MockexpenseRepository) New(userID int64, category string, amount uint64, date time.Time) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "New", userID, category, amount, date)
}

// New indicates an expected call of New.
func (mr *MockexpenseRepositoryMockRecorder) New(userID, category, amount, date interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "New", reflect.TypeOf((*MockexpenseRepository)(nil).New), userID, category, amount, date)
}

// MockuserRepository is a mock of userRepository interface.
type MockuserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockuserRepositoryMockRecorder
}

// MockuserRepositoryMockRecorder is the mock recorder for MockuserRepository.
type MockuserRepositoryMockRecorder struct {
	mock *MockuserRepository
}

// NewMockuserRepository creates a new mock instance.
func NewMockuserRepository(ctrl *gomock.Controller) *MockuserRepository {
	mock := &MockuserRepository{ctrl: ctrl}
	mock.recorder = &MockuserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockuserRepository) EXPECT() *MockuserRepositoryMockRecorder {
	return m.recorder
}

// DelLimit mocks base method.
func (m *MockuserRepository) DelLimit(userID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DelLimit", userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DelLimit indicates an expected call of DelLimit.
func (mr *MockuserRepositoryMockRecorder) DelLimit(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DelLimit", reflect.TypeOf((*MockuserRepository)(nil).DelLimit), userID)
}

// GetCurrency mocks base method.
func (m *MockuserRepository) GetCurrency(userID int64) *string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurrency", userID)
	ret0, _ := ret[0].(*string)
	return ret0
}

// GetCurrency indicates an expected call of GetCurrency.
func (mr *MockuserRepositoryMockRecorder) GetCurrency(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurrency", reflect.TypeOf((*MockuserRepository)(nil).GetCurrency), userID)
}

// GetLimit mocks base method.
func (m *MockuserRepository) GetLimit(userID int64) *uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLimit", userID)
	ret0, _ := ret[0].(*uint64)
	return ret0
}

// GetLimit indicates an expected call of GetLimit.
func (mr *MockuserRepositoryMockRecorder) GetLimit(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLimit", reflect.TypeOf((*MockuserRepository)(nil).GetLimit), userID)
}

// SetCurrency mocks base method.
func (m *MockuserRepository) SetCurrency(userID int64, currency string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetCurrency", userID, currency)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetCurrency indicates an expected call of SetCurrency.
func (mr *MockuserRepositoryMockRecorder) SetCurrency(userID, currency interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetCurrency", reflect.TypeOf((*MockuserRepository)(nil).SetCurrency), userID, currency)
}

// SetLimit mocks base method.
func (m *MockuserRepository) SetLimit(userID int64, limit uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetLimit", userID, limit)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetLimit indicates an expected call of SetLimit.
func (mr *MockuserRepositoryMockRecorder) SetLimit(userID, limit interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetLimit", reflect.TypeOf((*MockuserRepository)(nil).SetLimit), userID, limit)
}

// MockmessageSender is a mock of messageSender interface.
type MockmessageSender struct {
	ctrl     *gomock.Controller
	recorder *MockmessageSenderMockRecorder
}

// MockmessageSenderMockRecorder is the mock recorder for MockmessageSender.
type MockmessageSenderMockRecorder struct {
	mock *MockmessageSender
}

// NewMockmessageSender creates a new mock instance.
func NewMockmessageSender(ctrl *gomock.Controller) *MockmessageSender {
	mock := &MockmessageSender{ctrl: ctrl}
	mock.recorder = &MockmessageSenderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockmessageSender) EXPECT() *MockmessageSenderMockRecorder {
	return m.recorder
}

// SendMessage mocks base method.
func (m *MockmessageSender) SendMessage(text string, cases []string, userID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessage", text, cases, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMessage indicates an expected call of SendMessage.
func (mr *MockmessageSenderMockRecorder) SendMessage(text, cases, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessage", reflect.TypeOf((*MockmessageSender)(nil).SendMessage), text, cases, userID)
}
