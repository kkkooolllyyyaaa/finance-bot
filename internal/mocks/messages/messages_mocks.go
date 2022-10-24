// Code generated by MockGen. DO NOT EDIT.
// Source: internal/model/messages/msg_mgmt.go

// Package mock_messages is a generated GoMock package.
package mock_messages

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockMessageSender is a mock of MessageSender interface.
type MockMessageSender struct {
	ctrl     *gomock.Controller
	recorder *MockMessageSenderMockRecorder
}

// MockMessageSenderMockRecorder is the mock recorder for MockMessageSender.
type MockMessageSenderMockRecorder struct {
	mock *MockMessageSender
}

// NewMockMessageSender creates a new mock instance.
func NewMockMessageSender(ctrl *gomock.Controller) *MockMessageSender {
	mock := &MockMessageSender{ctrl: ctrl}
	mock.recorder = &MockMessageSenderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMessageSender) EXPECT() *MockMessageSenderMockRecorder {
	return m.recorder
}

// SendMessage mocks base method.
func (m *MockMessageSender) SendMessage(text string, userID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessage", text, userID)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMessage indicates an expected call of SendMessage.
func (mr *MockMessageSenderMockRecorder) SendMessage(text, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessage", reflect.TypeOf((*MockMessageSender)(nil).SendMessage), text, userID)
}
