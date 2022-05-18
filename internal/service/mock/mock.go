// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/service/serviceshell.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	models "voting-app/pkg/models"
	reflect "reflect"

	exterr "github.com/alexmolinanasaev/exterr"
	gomock "github.com/golang/mock/gomock"
)

// MockBlogApp is a mock of BlogApp interface.
type MockBlogApp struct {
	ctrl     *gomock.Controller
	recorder *MockBlogAppMockRecorder
}

// MockBlogAppMockRecorder is the mock recorder for MockBlogApp.
type MockBlogAppMockRecorder struct {
	mock *MockBlogApp
}

// NewMockBlogApp creates a new mock instance.
func NewMockBlogApp(ctrl *gomock.Controller) *MockBlogApp {
	mock := &MockBlogApp{ctrl: ctrl}
	mock.recorder = &MockBlogAppMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBlogApp) EXPECT() *MockBlogAppMockRecorder {
	return m.recorder
}

// AddFavoritPost mocks base method.
func (m *MockBlogApp) AddFavoritPost(userID, postID int) exterr.ErrExtender {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddFavoritPost", userID, postID)
	ret0, _ := ret[0].(exterr.ErrExtender)
	return ret0
}

// AddFavoritPost indicates an expected call of AddFavoritPost.
func (mr *MockBlogAppMockRecorder) AddFavoritPost(userID, postID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddFavoritPost", reflect.TypeOf((*MockBlogApp)(nil).AddFavoritPost), userID, postID)
}

// CreatePost mocks base method.
func (m *MockBlogApp) CreatePost(incomingPost *models.Post) (int, exterr.ErrExtender) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePost", incomingPost)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(exterr.ErrExtender)
	return ret0, ret1
}

// CreatePost indicates an expected call of CreatePost.
func (mr *MockBlogAppMockRecorder) CreatePost(incomingPost interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePost", reflect.TypeOf((*MockBlogApp)(nil).CreatePost), incomingPost)
}

// DeletePost mocks base method.
func (m *MockBlogApp) DeletePost(id int) exterr.ErrExtender {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePost", id)
	ret0, _ := ret[0].(exterr.ErrExtender)
	return ret0
}

// DeletePost indicates an expected call of DeletePost.
func (mr *MockBlogAppMockRecorder) DeletePost(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePost", reflect.TypeOf((*MockBlogApp)(nil).DeletePost), id)
}

// DeleteUser mocks base method.
func (m *MockBlogApp) DeleteUser(id int) exterr.ErrExtender {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteUser", id)
	ret0, _ := ret[0].(exterr.ErrExtender)
	return ret0
}

// DeleteUser indicates an expected call of DeleteUser.
func (mr *MockBlogAppMockRecorder) DeleteUser(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteUser", reflect.TypeOf((*MockBlogApp)(nil).DeleteUser), id)
}

// EditPost mocks base method.
func (m *MockBlogApp) EditPost(newPost *models.Post, admin bool) exterr.ErrExtender {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EditPost", newPost, admin)
	ret0, _ := ret[0].(exterr.ErrExtender)
	return ret0
}

// EditPost indicates an expected call of EditPost.
func (mr *MockBlogAppMockRecorder) EditPost(newPost, admin interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EditPost", reflect.TypeOf((*MockBlogApp)(nil).EditPost), newPost, admin)
}

// GetPost mocks base method.
func (m *MockBlogApp) GetPost(id int) (*models.Post, exterr.ErrExtender) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPost", id)
	ret0, _ := ret[0].(*models.Post)
	ret1, _ := ret[1].(exterr.ErrExtender)
	return ret0, ret1
}

// GetPost indicates an expected call of GetPost.
func (mr *MockBlogAppMockRecorder) GetPost(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPost", reflect.TypeOf((*MockBlogApp)(nil).GetPost), id)
}

// GetPostAuthorID mocks base method.
func (m *MockBlogApp) GetPostAuthorID(id int) (int, exterr.ErrExtender) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPostAuthorID", id)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(exterr.ErrExtender)
	return ret0, ret1
}

// GetPostAuthorID indicates an expected call of GetPostAuthorID.
func (mr *MockBlogAppMockRecorder) GetPostAuthorID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPostAuthorID", reflect.TypeOf((*MockBlogApp)(nil).GetPostAuthorID), id)
}

// GetUser mocks base method.
func (m *MockBlogApp) GetUser(id int) (*models.User, exterr.ErrExtender) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", id)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(exterr.ErrExtender)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockBlogAppMockRecorder) GetUser(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockBlogApp)(nil).GetUser), id)
}

// LikePost mocks base method.
func (m *MockBlogApp) LikePost(userID, postID int) exterr.ErrExtender {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LikePost", userID, postID)
	ret0, _ := ret[0].(exterr.ErrExtender)
	return ret0
}

// LikePost indicates an expected call of LikePost.
func (mr *MockBlogAppMockRecorder) LikePost(userID, postID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LikePost", reflect.TypeOf((*MockBlogApp)(nil).LikePost), userID, postID)
}

// SignIn mocks base method.
func (m *MockBlogApp) SignIn(email, password string) (*models.User, exterr.ErrExtender) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignIn", email, password)
	ret0, _ := ret[0].(*models.User)
	ret1, _ := ret[1].(exterr.ErrExtender)
	return ret0, ret1
}

// SignIn indicates an expected call of SignIn.
func (mr *MockBlogAppMockRecorder) SignIn(email, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignIn", reflect.TypeOf((*MockBlogApp)(nil).SignIn), email, password)
}

// SignUp mocks base method.
func (m *MockBlogApp) SignUp(user *models.User) (int, exterr.ErrExtender) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignUp", user)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(exterr.ErrExtender)
	return ret0, ret1
}

// SignUp indicates an expected call of SignUp.
func (mr *MockBlogAppMockRecorder) SignUp(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignUp", reflect.TypeOf((*MockBlogApp)(nil).SignUp), user)
}

// UnfavoritesPost mocks base method.
func (m *MockBlogApp) UnfavoritesPost(userID, postID int) exterr.ErrExtender {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnfavoritesPost", userID, postID)
	ret0, _ := ret[0].(exterr.ErrExtender)
	return ret0
}

// UnfavoritesPost indicates an expected call of UnfavoritesPost.
func (mr *MockBlogAppMockRecorder) UnfavoritesPost(userID, postID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnfavoritesPost", reflect.TypeOf((*MockBlogApp)(nil).UnfavoritesPost), userID, postID)
}

// UnlikePost mocks base method.
func (m *MockBlogApp) UnlikePost(userID, postID int) exterr.ErrExtender {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UnlikePost", userID, postID)
	ret0, _ := ret[0].(exterr.ErrExtender)
	return ret0
}

// UnlikePost indicates an expected call of UnlikePost.
func (mr *MockBlogAppMockRecorder) UnlikePost(userID, postID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UnlikePost", reflect.TypeOf((*MockBlogApp)(nil).UnlikePost), userID, postID)
}
