// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/IsuruHaupe/web-api/db/postgres (interfaces: Database)

// Package mockdb is a generated GoMock package.
package mockdb

import (
	context "context"
	reflect "reflect"

	db "github.com/IsuruHaupe/web-api/db/sqlc"
	gomock "github.com/golang/mock/gomock"
)

// MockDatabase is a mock of Database interface.
type MockDatabase struct {
	ctrl     *gomock.Controller
	recorder *MockDatabaseMockRecorder
}

// MockDatabaseMockRecorder is the mock recorder for MockDatabase.
type MockDatabaseMockRecorder struct {
	mock *MockDatabase
}

// NewMockDatabase creates a new mock instance.
func NewMockDatabase(ctrl *gomock.Controller) *MockDatabase {
	mock := &MockDatabase{ctrl: ctrl}
	mock.recorder = &MockDatabaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDatabase) EXPECT() *MockDatabaseMockRecorder {
	return m.recorder
}

// CreateContact mocks base method.
func (m *MockDatabase) CreateContact(arg0 context.Context, arg1 db.CreateContactParams) (db.Contact, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateContact", arg0, arg1)
	ret0, _ := ret[0].(db.Contact)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateContact indicates an expected call of CreateContact.
func (mr *MockDatabaseMockRecorder) CreateContact(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateContact", reflect.TypeOf((*MockDatabase)(nil).CreateContact), arg0, arg1)
}

// CreateContactHasSkill mocks base method.
func (m *MockDatabase) CreateContactHasSkill(arg0 context.Context, arg1 db.CreateContactHasSkillParams) (db.ContactHasSkill, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateContactHasSkill", arg0, arg1)
	ret0, _ := ret[0].(db.ContactHasSkill)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateContactHasSkill indicates an expected call of CreateContactHasSkill.
func (mr *MockDatabaseMockRecorder) CreateContactHasSkill(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateContactHasSkill", reflect.TypeOf((*MockDatabase)(nil).CreateContactHasSkill), arg0, arg1)
}

// CreateSkill mocks base method.
func (m *MockDatabase) CreateSkill(arg0 context.Context, arg1 db.CreateSkillParams) (db.Skill, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSkill", arg0, arg1)
	ret0, _ := ret[0].(db.Skill)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSkill indicates an expected call of CreateSkill.
func (mr *MockDatabaseMockRecorder) CreateSkill(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSkill", reflect.TypeOf((*MockDatabase)(nil).CreateSkill), arg0, arg1)
}

// CreateUser mocks base method.
func (m *MockDatabase) CreateUser(arg0 context.Context, arg1 db.CreateUserParams) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockDatabaseMockRecorder) CreateUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockDatabase)(nil).CreateUser), arg0, arg1)
}

// DeleteContact mocks base method.
func (m *MockDatabase) DeleteContact(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteContact", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteContact indicates an expected call of DeleteContact.
func (mr *MockDatabaseMockRecorder) DeleteContact(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteContact", reflect.TypeOf((*MockDatabase)(nil).DeleteContact), arg0, arg1)
}

// DeleteSkill mocks base method.
func (m *MockDatabase) DeleteSkill(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteSkill", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteSkill indicates an expected call of DeleteSkill.
func (mr *MockDatabaseMockRecorder) DeleteSkill(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteSkill", reflect.TypeOf((*MockDatabase)(nil).DeleteSkill), arg0, arg1)
}

// DeleteUser mocks base method.
func (m *MockDatabase) DeleteUser(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockDatabaseMockRecorder) DeleteUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockDatabase)(nil).DeleteUser), arg0, arg1)
}

// GetContact mocks base method.
func (m *MockDatabase) GetContact(arg0 context.Context, arg1 int64) (db.Contact, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContact", arg0, arg1)
	ret0, _ := ret[0].(db.Contact)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetContact indicates an expected call of GetContact.
func (mr *MockDatabaseMockRecorder) GetContact(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContact", reflect.TypeOf((*MockDatabase)(nil).GetContact), arg0, arg1)
}

// GetContactsWithSkill mocks base method.
func (m *MockDatabase) GetContactsWithSkill(arg0 context.Context, arg1 string) ([]db.Contact, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContactsWithSkill", arg0, arg1)
	ret0, _ := ret[0].([]db.Contact)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetContactsWithSkill indicates an expected call of GetContactsWithSkill.
func (mr *MockDatabaseMockRecorder) GetContactsWithSkill(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContactsWithSkill", reflect.TypeOf((*MockDatabase)(nil).GetContactsWithSkill), arg0, arg1)
}

// GetContactsWithSkillAndLevel mocks base method.
func (m *MockDatabase) GetContactsWithSkillAndLevel(arg0 context.Context, arg1 db.GetContactsWithSkillAndLevelParams) ([]db.Contact, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContactsWithSkillAndLevel", arg0, arg1)
	ret0, _ := ret[0].([]db.Contact)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetContactsWithSkillAndLevel indicates an expected call of GetContactsWithSkillAndLevel.
func (mr *MockDatabaseMockRecorder) GetContactsWithSkillAndLevel(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContactsWithSkillAndLevel", reflect.TypeOf((*MockDatabase)(nil).GetContactsWithSkillAndLevel), arg0, arg1)
}

// GetEmail mocks base method.
func (m *MockDatabase) GetEmail(arg0 context.Context, arg1 int64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEmail", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEmail indicates an expected call of GetEmail.
func (mr *MockDatabaseMockRecorder) GetEmail(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEmail", reflect.TypeOf((*MockDatabase)(nil).GetEmail), arg0, arg1)
}

// GetFirstname mocks base method.
func (m *MockDatabase) GetFirstname(arg0 context.Context, arg1 int64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFirstname", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFirstname indicates an expected call of GetFirstname.
func (mr *MockDatabaseMockRecorder) GetFirstname(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFirstname", reflect.TypeOf((*MockDatabase)(nil).GetFirstname), arg0, arg1)
}

// GetFullname mocks base method.
func (m *MockDatabase) GetFullname(arg0 context.Context, arg1 int64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFullname", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFullname indicates an expected call of GetFullname.
func (mr *MockDatabaseMockRecorder) GetFullname(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFullname", reflect.TypeOf((*MockDatabase)(nil).GetFullname), arg0, arg1)
}

// GetHomeAddress mocks base method.
func (m *MockDatabase) GetHomeAddress(arg0 context.Context, arg1 int64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetHomeAddress", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetHomeAddress indicates an expected call of GetHomeAddress.
func (mr *MockDatabaseMockRecorder) GetHomeAddress(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetHomeAddress", reflect.TypeOf((*MockDatabase)(nil).GetHomeAddress), arg0, arg1)
}

// GetIfExistsContactID mocks base method.
func (m *MockDatabase) GetIfExistsContactID(arg0 context.Context, arg1 int64) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIfExistsContactID", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetIfExistsContactID indicates an expected call of GetIfExistsContactID.
func (mr *MockDatabaseMockRecorder) GetIfExistsContactID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIfExistsContactID", reflect.TypeOf((*MockDatabase)(nil).GetIfExistsContactID), arg0, arg1)
}

// GetIfExistsSkillID mocks base method.
func (m *MockDatabase) GetIfExistsSkillID(arg0 context.Context, arg1 int64) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetIfExistsSkillID", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetIfExistsSkillID indicates an expected call of GetIfExistsSkillID.
func (mr *MockDatabaseMockRecorder) GetIfExistsSkillID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetIfExistsSkillID", reflect.TypeOf((*MockDatabase)(nil).GetIfExistsSkillID), arg0, arg1)
}

// GetLastname mocks base method.
func (m *MockDatabase) GetLastname(arg0 context.Context, arg1 int64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastname", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLastname indicates an expected call of GetLastname.
func (mr *MockDatabaseMockRecorder) GetLastname(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastname", reflect.TypeOf((*MockDatabase)(nil).GetLastname), arg0, arg1)
}

// GetPhoneNumber mocks base method.
func (m *MockDatabase) GetPhoneNumber(arg0 context.Context, arg1 int64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPhoneNumber", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPhoneNumber indicates an expected call of GetPhoneNumber.
func (mr *MockDatabaseMockRecorder) GetPhoneNumber(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPhoneNumber", reflect.TypeOf((*MockDatabase)(nil).GetPhoneNumber), arg0, arg1)
}

// GetSkill mocks base method.
func (m *MockDatabase) GetSkill(arg0 context.Context, arg1 int64) (db.Skill, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSkill", arg0, arg1)
	ret0, _ := ret[0].(db.Skill)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSkill indicates an expected call of GetSkill.
func (mr *MockDatabaseMockRecorder) GetSkill(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSkill", reflect.TypeOf((*MockDatabase)(nil).GetSkill), arg0, arg1)
}

// GetSkillLevel mocks base method.
func (m *MockDatabase) GetSkillLevel(arg0 context.Context, arg1 int64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSkillLevel", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSkillLevel indicates an expected call of GetSkillLevel.
func (mr *MockDatabaseMockRecorder) GetSkillLevel(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSkillLevel", reflect.TypeOf((*MockDatabase)(nil).GetSkillLevel), arg0, arg1)
}

// GetSkillName mocks base method.
func (m *MockDatabase) GetSkillName(arg0 context.Context, arg1 int64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSkillName", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSkillName indicates an expected call of GetSkillName.
func (mr *MockDatabaseMockRecorder) GetSkillName(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSkillName", reflect.TypeOf((*MockDatabase)(nil).GetSkillName), arg0, arg1)
}

// GetUser mocks base method.
func (m *MockDatabase) GetUser(arg0 context.Context, arg1 string) (db.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", arg0, arg1)
	ret0, _ := ret[0].(db.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockDatabaseMockRecorder) GetUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockDatabase)(nil).GetUser), arg0, arg1)
}

// ListContacts mocks base method.
func (m *MockDatabase) ListContacts(arg0 context.Context, arg1 db.ListContactsParams) ([]db.Contact, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListContacts", arg0, arg1)
	ret0, _ := ret[0].([]db.Contact)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListContacts indicates an expected call of ListContacts.
func (mr *MockDatabaseMockRecorder) ListContacts(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListContacts", reflect.TypeOf((*MockDatabase)(nil).ListContacts), arg0, arg1)
}

// ListSkills mocks base method.
func (m *MockDatabase) ListSkills(arg0 context.Context, arg1 db.ListSkillsParams) ([]db.Skill, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListSkills", arg0, arg1)
	ret0, _ := ret[0].([]db.Skill)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListSkills indicates an expected call of ListSkills.
func (mr *MockDatabaseMockRecorder) ListSkills(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListSkills", reflect.TypeOf((*MockDatabase)(nil).ListSkills), arg0, arg1)
}

// UpdateContact mocks base method.
func (m *MockDatabase) UpdateContact(arg0 context.Context, arg1 db.UpdateContactParams) (db.Contact, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateContact", arg0, arg1)
	ret0, _ := ret[0].(db.Contact)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateContact indicates an expected call of UpdateContact.
func (mr *MockDatabaseMockRecorder) UpdateContact(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateContact", reflect.TypeOf((*MockDatabase)(nil).UpdateContact), arg0, arg1)
}

// UpdateSkill mocks base method.
func (m *MockDatabase) UpdateSkill(arg0 context.Context, arg1 db.UpdateSkillParams) (db.Skill, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateSkill", arg0, arg1)
	ret0, _ := ret[0].(db.Skill)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateSkill indicates an expected call of UpdateSkill.
func (mr *MockDatabaseMockRecorder) UpdateSkill(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateSkill", reflect.TypeOf((*MockDatabase)(nil).UpdateSkill), arg0, arg1)
}