// Automatically generated by MockGen. DO NOT EDIT!
// Source: ./systest/api_v1/api.httpgw.go

package api_v1

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
)

// Mock of APIInterface interface
type MockAPIInterface struct {
	ctrl     *gomock.Controller
	recorder *_MockAPIInterfaceRecorder
}

// Recorder for MockAPIInterface (not exported)
type _MockAPIInterfaceRecorder struct {
	mock *MockAPIInterface
}

func NewMockAPIInterface(ctrl *gomock.Controller) *MockAPIInterface {
	mock := &MockAPIInterface{ctrl: ctrl}
	mock.recorder = &_MockAPIInterfaceRecorder{mock}
	return mock
}

func (_m *MockAPIInterface) EXPECT() *_MockAPIInterfaceRecorder {
	return _m.recorder
}

func (_m *MockAPIInterface) Transfer(_param0 context.Context, _param1 *TransferRequest) (*TransferResponse, error) {
	ret := _m.ctrl.Call(_m, "Transfer", _param0, _param1)
	ret0, _ := ret[0].(*TransferResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockAPIInterfaceRecorder) Transfer(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Transfer", arg0, arg1)
}

func (_m *MockAPIInterface) UpdateSettings(_param0 context.Context, _param1 *UpdateSettingsRequest) (*UpdateSettingsResponse, error) {
	ret := _m.ctrl.Call(_m, "UpdateSettings", _param0, _param1)
	ret0, _ := ret[0].(*UpdateSettingsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockAPIInterfaceRecorder) UpdateSettings(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "UpdateSettings", arg0, arg1)
}

func (_m *MockAPIInterface) GetPrevHash(_param0 context.Context, _param1 *PrevHashRequest) (*PrevHashResponse, error) {
	ret := _m.ctrl.Call(_m, "GetPrevHash", _param0, _param1)
	ret0, _ := ret[0].(*PrevHashResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockAPIInterfaceRecorder) GetPrevHash(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetPrevHash", arg0, arg1)
}

func (_m *MockAPIInterface) GetHistory(_param0 context.Context, _param1 *HistoryRequest) (*HistoryResponse, error) {
	ret := _m.ctrl.Call(_m, "GetHistory", _param0, _param1)
	ret0, _ := ret[0].(*HistoryResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockAPIInterfaceRecorder) GetHistory(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetHistory", arg0, arg1)
}

func (_m *MockAPIInterface) GetStats(_param0 context.Context, _param1 *StatsRequest) (*StatsResponse, error) {
	ret := _m.ctrl.Call(_m, "GetStats", _param0, _param1)
	ret0, _ := ret[0].(*StatsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockAPIInterfaceRecorder) GetStats(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetStats", arg0, arg1)
}

func (_m *MockAPIInterface) GetAccounts(_param0 context.Context, _param1 *AccountsRequest) (*AccountsResponse, error) {
	ret := _m.ctrl.Call(_m, "GetAccounts", _param0, _param1)
	ret0, _ := ret[0].(*AccountsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockAPIInterfaceRecorder) GetAccounts(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetAccounts", arg0, arg1)
}

func (_m *MockAPIInterface) GetAccountSettings(_param0 context.Context, _param1 *AccountSettingsRequest) (*AccountSettingsResponse, error) {
	ret := _m.ctrl.Call(_m, "GetAccountSettings", _param0, _param1)
	ret0, _ := ret[0].(*AccountSettingsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockAPIInterfaceRecorder) GetAccountSettings(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetAccountSettings", arg0, arg1)
}
