package service_test

import (
	"net/http"
	"testing"

	mock_repository "X-Blog/internal/repository/mock"
	"X-Blog/internal/service"
	"X-Blog/pkg/models"

	"github.com/alexmolinanasaev/exterr"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

// E-mails
var testEmailInvalid string = "[invalidEmail]"

// Passwords
var testPasswordStrong string = "Test123!_test890#"
var testPasswordWeak string = "weakpsswd"

// Users
var testUserEmpty = models.User{}

var testUserUser = models.User{
	FirstName: "Test",
	LastName:  "Testov",
	Email:     "test@test.test",
	Password:  testPasswordStrong,
}

var testUserUserRegistered = models.User{
	FirstName:   testUserUser.FirstName,
	LastName:    testUserUser.LastName,
	Email:       testUserUser.Email,
	Password:    service.GeneratePasswordHash(testUserUser.Password),
	AccessLevel: models.UserAccess,
}

var testUserAdmin = models.User{
	FirstName: "Admin",
	LastName:  "Adminov",
	Email:     "admin@admin.admin",
	Password:  testPasswordStrong,
}

// Posts

var testPost = models.Post{
	Title:   "Test Title",
	Content: "blablablablablabla",
	Author:  1,
}

func TestService_SignUp(t *testing.T) {
	type mockBehaviorGetUserByEmail func(m *mock_repository.MockBlogApp, email string)
	type mockBehaviorSignUp func(m *mock_repository.MockBlogApp, user *models.User)

	testTable := []struct {
		name                       string                     // Название теста
		user                       models.User                // Тестовая модель пользователя
		userRegistered             models.User                // Тестовая модель зарегистрированного пользователя
		mockBehaviorGetUserByEmail mockBehaviorGetUserByEmail // Замоканная функция получения пользователя по Email
		mockBehaviorSignUp         mockBehaviorSignUp         // Замоконная функция авторизации
		expectedID                 int                        // ожидаемый ID поста, который возвращается после регистрации поста
		expectedError              exterr.ErrExtender         // ожидамая ошибка, которая возвращается
	}{
		// -----------------------------------------------START--------------------------------------------
		{ // -------------------------------------------Case #1-------------------------------------------
			name:           "[CONDITION]:OK;[RESULT]:OK",
			user:           testUserUser,
			userRegistered: testUserUserRegistered,
			mockBehaviorGetUserByEmail: func(m *mock_repository.MockBlogApp, email string) {
				m.EXPECT().
					GetUserByEmail(email).
					Return(nil, exterr.New("User not found"))
			},
			mockBehaviorSignUp: func(m *mock_repository.MockBlogApp, user *models.User) {
				m.EXPECT().
					SignUp(user).
					Return(1, nil)
			},
			expectedID:    1,
			expectedError: nil,
		}, { // ----------------------------------------Case #2-------------------------------------------
			name:                       "[CONDITION]:Empty user info;[RESULT]:Empty user info ERROR",
			user:                       testUserEmpty,
			mockBehaviorGetUserByEmail: func(m *mock_repository.MockBlogApp, email string) {},
			mockBehaviorSignUp:         func(m *mock_repository.MockBlogApp, user *models.User) {},
			expectedID:                 0,
			expectedError: exterr.New("Empty user info").
				SetErrCode(http.StatusBadRequest).
				SetAltMsg("Empty user info (FirstName, LastName, Email or Password)"),
		}, { // ----------------------------------------Case #3-------------------------------------------
			name: "[CONDITION]:OK;[RESULT]:E-mail already taken ERROR",
			user: testUserUser,
			mockBehaviorGetUserByEmail: func(m *mock_repository.MockBlogApp, email string) {
				m.EXPECT().
					GetUserByEmail(email).
					Return(&testUserUser, nil)
			},
			mockBehaviorSignUp: func(m *mock_repository.MockBlogApp, user *models.User) {},
			expectedID:         0,
			expectedError: exterr.New("Service: user found, can not create new user").
				SetErrCode(http.StatusBadRequest).
				SetAltMsg("Email already taken"),
		}, { // ----------------------------------------Case #4-------------------------------------------
			name:           "[CONDITION]:OK;[RESULT]:Repository SignUp ERROR",
			user:           testUserUser,
			userRegistered: testUserUserRegistered,
			mockBehaviorGetUserByEmail: func(m *mock_repository.MockBlogApp, email string) {
				m.EXPECT().
					GetUserByEmail(email).
					Return(nil, exterr.New("User not found"))
			},
			mockBehaviorSignUp: func(m *mock_repository.MockBlogApp, user *models.User) {
				m.EXPECT().
					SignUp(user).
					Return(0, exterr.New("Repository error"))
			},
			expectedID: 0,
			expectedError: exterr.NewWithExtErr("Service can not add new user in repository", exterr.New("Repository error")).
				SetErrCode(http.StatusInternalServerError).
				SetAltMsg(http.StatusText(http.StatusInternalServerError)),
		},
	} // -------------------------------------------------END---------------------------------------------

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Arrange
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			mockRepository := mock_repository.NewMockBlogApp(mockController)
			testCase.mockBehaviorGetUserByEmail(mockRepository, testCase.user.Email)
			testCase.mockBehaviorSignUp(mockRepository, &testCase.userRegistered)
			s := service.NewGetService(mockRepository)

			// Act
			actualID, actualError := s.SignUp(&testCase.user)

			// Assert
			assert.Equal(t, testCase.expectedID, actualID)
			if actualError != nil { // Если ошибка присутствует, то проверяем всё, кроме Trace.
				assert.Equal(t, testCase.expectedError.Error(), actualError.Error())
				assert.Equal(t, testCase.expectedError.GetErrCode(), actualError.GetErrCode())
				assert.Equal(t, testCase.expectedError.GetAltMsg(), actualError.GetAltMsg())
			} else {
				assert.Equal(t, testCase.expectedError, actualError)
			}
		})
	}
}

func TestService_SignIn(t *testing.T) {
	type mockBehaviorGetUserByEmail func(m *mock_repository.MockBlogApp, email string)

	testTable := []struct {
		name                       string
		email                      string
		password                   string
		mockBehaviorGetUserByEmail mockBehaviorGetUserByEmail
		expectedUser               *models.User
		expectedError              exterr.ErrExtender
	}{ // -----------------------------------------------START--------------------------------------------
		{ // -------------------------------------------Case #1-------------------------------------------
			name:     "[CONDITION]:OK;[RESULT]:OK",
			email:    testUserUser.Email,
			password: testUserUser.Password,
			mockBehaviorGetUserByEmail: func(m *mock_repository.MockBlogApp, email string) {
				m.EXPECT().
					GetUserByEmail(email).
					Return(&testUserUserRegistered, nil)
			},
			expectedUser:  &testUserUserRegistered,
			expectedError: nil,
		}, { // -----------------------------------------Case #2------------------------------------------
			name:     "[CONDITION]:OK;[RESULT]:GetUserByEmail user not found ERROR",
			email:    testUserUser.Email,
			password: testUserUser.Password,
			mockBehaviorGetUserByEmail: func(m *mock_repository.MockBlogApp, email string) {
				m.EXPECT().
					GetUserByEmail(email).
					Return(nil, exterr.New("User not found"))
			},
			expectedUser: nil,
			expectedError: exterr.NewWithExtErr("Service: user not found ", exterr.New("User not found")).
				SetErrCode(http.StatusNotFound).
				SetAltMsg("User not found"),
		}, { // -----------------------------------------Case #3------------------------------------------
			name:     "[CONDITION]:Wrong password;[RESULT]:Wrong password ERROR",
			email:    testUserUser.Email,
			password: "wrong password",
			mockBehaviorGetUserByEmail: func(m *mock_repository.MockBlogApp, email string) {
				m.EXPECT().
					GetUserByEmail(email).
					Return(&testUserUserRegistered, nil)
			},
			expectedUser: nil,
			expectedError: exterr.New("Wrong password").
				SetErrCode(http.StatusUnauthorized).
				SetAltMsg("Wrong password"),
		},
	} // -------------------------------------------------END---------------------------------------------
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Arrange
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			mockRepository := mock_repository.NewMockBlogApp(mockController)
			testCase.mockBehaviorGetUserByEmail(mockRepository, testCase.email)
			s := service.NewGetService(mockRepository)

			// Act

			actualUser, actualError := s.SignIn(testCase.email, testCase.password)

			// Assert
			assert.Equal(t, testCase.expectedUser, actualUser)

			if actualError != nil { // Если ошибка присутствует, то проверяем всё, кроме Trace.
				assert.Equal(t, testCase.expectedError.Error(), actualError.Error())
				assert.Equal(t, testCase.expectedError.GetErrCode(), actualError.GetErrCode())
				assert.Equal(t, testCase.expectedError.GetAltMsg(), actualError.GetAltMsg())
			} else {
				assert.Equal(t, testCase.expectedError, actualError)
			}
		})
	}
}

func TestService_GetUser(t *testing.T) {
	type mockBehaviorGetUserByID func(m *mock_repository.MockBlogApp, id int)

	testTable := []struct {
		name                    string
		id                      int
		mockBehaviorGetUserByID mockBehaviorGetUserByID
		expectedUser            *models.User
		expectedError           exterr.ErrExtender
	}{ // -----------------------------------------------START--------------------------------------------
		{ // -------------------------------------------Case #1-------------------------------------------
			name: "[CONDITION]:OK;[RESULT]:OK",
			id:   1,
			mockBehaviorGetUserByID: func(m *mock_repository.MockBlogApp, id int) {
				m.EXPECT().
					GetUserByID(id).
					Return(&testUserUser, nil)
			},
			expectedUser:  &testUserUser,
			expectedError: nil,
		}, { // ----------------------------------------Case #2-------------------------------------------
			name: "[CONDITION]:OK;[RESULT]:GetUserByID user not found ERROR",
			id:   1,
			mockBehaviorGetUserByID: func(m *mock_repository.MockBlogApp, id int) {
				m.EXPECT().
					GetUserByID(id).
					Return(nil, exterr.New("User not found"))
			},
			expectedUser: nil,
			expectedError: exterr.NewWithExtErr("Service: GetUser error", exterr.New("User not found")).
				SetErrCode(http.StatusNotFound).
				SetAltMsg("User not found"),
		},
	} // -------------------------------------------------END---------------------------------------------
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Arrange
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			mockRepository := mock_repository.NewMockBlogApp(mockController)
			testCase.mockBehaviorGetUserByID(mockRepository, testCase.id)
			s := service.NewGetService(mockRepository)

			// Act
			actualUser, actualError := s.GetUser(testCase.id)
			// Assert
			assert.Equal(t, testCase.expectedUser, actualUser)

			if actualError != nil { // Если ошибка присутствует, то проверяем всё, кроме Trace.
				assert.Equal(t, testCase.expectedError.Error(), actualError.Error())
				assert.Equal(t, testCase.expectedError.GetErrCode(), actualError.GetErrCode())
				assert.Equal(t, testCase.expectedError.GetAltMsg(), actualError.GetAltMsg())
			} else {
				assert.Equal(t, testCase.expectedError, actualError)
			}
		})
	}
}

func TestService_DeleteUser(t *testing.T) {
	type mockBehaviorGetUserByID func(m *mock_repository.MockBlogApp, id int)
	type mockBehaviorDeleteUser func(m *mock_repository.MockBlogApp, id int)

	testTable := []struct {
		name                    string
		id                      int
		mockBehaviorGetUserByID mockBehaviorGetUserByID
		mockBehaviorDeleteUser  mockBehaviorDeleteUser
		expectedError           exterr.ErrExtender
	}{ // -----------------------------------------------START--------------------------------------------
		{ // -------------------------------------------Case #1-------------------------------------------
			name: "[CONDITION]:OK;[RESULT]:OK",
			id:   1,
			mockBehaviorGetUserByID: func(m *mock_repository.MockBlogApp, id int) {
				m.EXPECT().
					GetUserByID(id).
					Return(&testUserUser, nil)
			},
			mockBehaviorDeleteUser: func(m *mock_repository.MockBlogApp, id int) {
				m.EXPECT().
					DeleteUser(id).
					Return(nil)
			},
			expectedError: nil,
		}, { // ----------------------------------------Case #2-------------------------------------------
			name: "[CONDITION]:OK;[RESULT]:GetUserByID user not found ERROR",
			id:   1,
			mockBehaviorGetUserByID: func(m *mock_repository.MockBlogApp, id int) {
				m.EXPECT().
					GetUserByID(id).
					Return(nil, exterr.New("User not found"))
			},
			mockBehaviorDeleteUser: func(m *mock_repository.MockBlogApp, id int) {},
			expectedError: exterr.NewWithExtErr("Service: DeleteUser error", exterr.New("User not found")).
				SetErrCode(http.StatusNotFound).
				SetAltMsg("User not found"),
		}, { // ----------------------------------------Case #3-------------------------------------------
			name: "[CONDITION]:OK;[RESULT]:Repository DeleteUser ERROR",
			id:   1,
			mockBehaviorGetUserByID: func(m *mock_repository.MockBlogApp, id int) {
				m.EXPECT().
					GetUserByID(id).
					Return(&testUserUser, nil)
			},
			mockBehaviorDeleteUser: func(m *mock_repository.MockBlogApp, id int) {
				m.EXPECT().
					DeleteUser(id).
					Return(exterr.New("Repository error"))
			},
			expectedError: exterr.NewWithExtErr("Service: DeleteUser error", exterr.New("Repository error")).
				SetErrCode(http.StatusInternalServerError).
				SetAltMsg(http.StatusText(http.StatusInternalServerError)),
		},
	} // -------------------------------------------------END---------------------------------------------
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Arrange
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			mockRepository := mock_repository.NewMockBlogApp(mockController)
			testCase.mockBehaviorGetUserByID(mockRepository, testCase.id)
			testCase.mockBehaviorDeleteUser(mockRepository, testCase.id)
			s := service.NewGetService(mockRepository)

			// Act
			actualError := s.DeleteUser(testCase.id)

			// Assert
			if actualError != nil { // Если ошибка присутствует, то проверяем всё, кроме Trace.
				assert.Equal(t, testCase.expectedError.Error(), actualError.Error())
				assert.Equal(t, testCase.expectedError.GetErrCode(), actualError.GetErrCode())
				assert.Equal(t, testCase.expectedError.GetAltMsg(), actualError.GetAltMsg())
			} else {
				assert.Equal(t, testCase.expectedError, actualError)
			}
		})
	}
}

func TestService_CreatePost(t *testing.T) {
	type mockBehaviorCreatePost func(m *mock_repository.MockBlogApp, post *models.Post)

	testTable := []struct {
		name                   string
		incomingPost           *models.Post
		mockBehaviorCreatePost mockBehaviorCreatePost
		expectedID             int
		expectedError          exterr.ErrExtender
	}{ // -----------------------------------------------START--------------------------------------------
		{ // -------------------------------------------Case #1-------------------------------------------
			name:         "[CONDITION]:OK;[RESULT]:OK",
			incomingPost: &testPost,
			mockBehaviorCreatePost: func(m *mock_repository.MockBlogApp, post *models.Post) {
				m.EXPECT().
					CreatePost(post).
					Return(99, nil)
			},
			expectedID:    99,
			expectedError: nil,
		}, { // ----------------------------------------Case #2-------------------------------------------
			name: "[CONDITION]:Empty title;[RESULT]:Empty title ERROR",
			incomingPost: &models.Post{
				Title:   "",
				Content: testPost.Content,
				Author:  testPost.Author,
			},
			mockBehaviorCreatePost: func(m *mock_repository.MockBlogApp, post *models.Post) {},
			expectedID:             -1,
			expectedError: exterr.New("Service: CreatePost empty title").
				SetErrCode(http.StatusBadRequest).
				SetAltMsg("Empty title"),
		}, { // ----------------------------------------Case #3-------------------------------------------
			name: "[CONDITION]:Empty author;[RESULT]:Empty author ERROR",
			incomingPost: &models.Post{
				Title:   testPost.Title,
				Content: testPost.Content,
				Author:  0,
			},
			mockBehaviorCreatePost: func(m *mock_repository.MockBlogApp, post *models.Post) {},
			expectedID:             -1,
			expectedError: exterr.New("Service: CreatePost empty author").
				SetErrCode(http.StatusBadRequest).
				SetAltMsg("Empty author"),
		}, { // ----------------------------------------Case #4-------------------------------------------
			name:         "[CONDITION]:OK;[RESULT]:Repository CreatePost ERROR",
			incomingPost: &testPost,
			mockBehaviorCreatePost: func(m *mock_repository.MockBlogApp, post *models.Post) {
				m.EXPECT().
					CreatePost(post).
					Return(0, exterr.New("Repository error"))
			},
			expectedID: -1,
			expectedError: exterr.NewWithExtErr("Service: CreatePost error", exterr.New("Repository error")).
				SetErrCode(http.StatusInternalServerError).
				SetAltMsg(http.StatusText(http.StatusInternalServerError)),
		},
	} // -------------------------------------------------END---------------------------------------------
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Arrange
			postReadyToRecord := &models.Post{
				Title:           testCase.incomingPost.Title,
				Author:          testCase.incomingPost.Author,
				Content:         testCase.incomingPost.Content,
				AmountLikes:     0,
				AmountFavorites: 0,
				AccessLevel:     models.UserAccess,
			}

			mockController := gomock.NewController(t)
			defer mockController.Finish()

			mockRepository := mock_repository.NewMockBlogApp(mockController)
			testCase.mockBehaviorCreatePost(mockRepository, postReadyToRecord)
			s := service.NewGetService(mockRepository)

			// Act
			actualID, actualError := s.CreatePost(testCase.incomingPost)

			// Assert
			assert.Equal(t, testCase.expectedID, actualID)

			if actualError != nil { // Если ошибка присутствует, то проверяем всё, кроме Trace.
				assert.Equal(t, testCase.expectedError.Error(), actualError.Error())
				assert.Equal(t, testCase.expectedError.GetErrCode(), actualError.GetErrCode())
				assert.Equal(t, testCase.expectedError.GetAltMsg(), actualError.GetAltMsg())
			} else {
				assert.Equal(t, testCase.expectedError, actualError)
			}
		})
	}
}

func TestService_GetPost(t *testing.T) {
	type mockBehaviorGetPost func(m *mock_repository.MockBlogApp, id int)

	testTable := []struct {
		name                string
		postID              int
		mockBehaviorGetPost mockBehaviorGetPost
		expectedPost        *models.Post
		expectedError       exterr.ErrExtender
	}{ // -----------------------------------------------START--------------------------------------------
		{ // -------------------------------------------Case #1-------------------------------------------
			name:   "[CONDITION]:OK;[RESULT]:OK",
			postID: 1,
			mockBehaviorGetPost: func(m *mock_repository.MockBlogApp, id int) {
				m.EXPECT().
					GetPost(id).
					Return(&testPost, nil)
			},
			expectedPost:  &testPost,
			expectedError: nil,
		}, { // ----------------------------------------Case #2-------------------------------------------
			name:                "Empty PostID",
			postID:              0,
			mockBehaviorGetPost: func(m *mock_repository.MockBlogApp, id int) {},
			expectedPost:        nil,
			expectedError: exterr.New("Service: GetPost empty id").
				SetErrCode(http.StatusBadRequest).
				SetAltMsg("Empty post id"),
		}, { // ----------------------------------------Case #3-------------------------------------------
			name:   "[CONDITION]:OK;[RESULT]:Repository GetPost ERROR",
			postID: 1,
			mockBehaviorGetPost: func(m *mock_repository.MockBlogApp, id int) {
				m.EXPECT().
					GetPost(id).
					Return(nil, exterr.New("Repository error"))
			},
			expectedPost: nil,
			expectedError: exterr.NewWithExtErr("Service: GetPost error", exterr.New("Repository error")).
				SetErrCode(http.StatusInternalServerError).
				SetAltMsg(http.StatusText(http.StatusInternalServerError)),
		},
	} // -------------------------------------------------END---------------------------------------------
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Arrange
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			mockRepository := mock_repository.NewMockBlogApp(mockController)
			testCase.mockBehaviorGetPost(mockRepository, testCase.postID)
			s := service.NewGetService(mockRepository)

			// Act
			actualPost, actualError := s.GetPost(testCase.postID)

			// Assert
			assert.Equal(t, testCase.expectedPost, actualPost)

			if actualError != nil { // Если ошибка присутствует, то проверяем всё, кроме Trace.
				assert.Equal(t, testCase.expectedError.Error(), actualError.Error())
				assert.Equal(t, testCase.expectedError.GetErrCode(), actualError.GetErrCode())
				assert.Equal(t, testCase.expectedError.GetAltMsg(), actualError.GetAltMsg())
			} else {
				assert.Equal(t, testCase.expectedError, actualError)
			}
		})
	}
}

func TestService_GetPostAuthorID(t *testing.T) {
	type mockBehaviorGetPostAuthorID func(m *mock_repository.MockBlogApp, postID int)

	testTable := []struct {
		name                        string
		postID                      int
		mockBehaviorGetPostAuthorID mockBehaviorGetPostAuthorID
		expectedAuthorID            int
		expectedError               exterr.ErrExtender
	}{ // -----------------------------------------------START--------------------------------------------
		{ // -------------------------------------------Case #1-------------------------------------------
			name:   "[CONDITION]:OK;[RESULT]:OK",
			postID: 1,
			mockBehaviorGetPostAuthorID: func(m *mock_repository.MockBlogApp, postID int) {
				m.EXPECT().
					GetPostAuthorID(postID).
					Return(123, nil)
			},
			expectedAuthorID: 123,
			expectedError:    nil,
		}, { // ----------------------------------------Case #2-------------------------------------------
			name:                        "[CONDITION]:Empty PostID;[RESULT]:Empty PostID ERROR",
			postID:                      0,
			mockBehaviorGetPostAuthorID: func(m *mock_repository.MockBlogApp, postID int) {},
			expectedAuthorID:            0,
			expectedError: exterr.New("Service: GetPostAuthor empty id").
				SetErrCode(http.StatusBadRequest).
				SetAltMsg("Empty post id"),
		}, { // ----------------------------------------Case #3-------------------------------------------
			name:   "[CONDITION]:OK;[RESULT]:Repository GetPostAuthorID ERROR",
			postID: 1,
			mockBehaviorGetPostAuthorID: func(m *mock_repository.MockBlogApp, postID int) {
				m.EXPECT().
					GetPostAuthorID(postID).
					Return(0, exterr.New("Repository error"))
			},
			expectedAuthorID: 0,
			expectedError: exterr.NewWithExtErr("Service: GetPostAuthor error", exterr.New("Repository error")).
				SetErrCode(http.StatusInternalServerError).
				SetAltMsg(http.StatusText(http.StatusInternalServerError)),
		},
	} // -------------------------------------------------END---------------------------------------------
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Arrange
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			mockRepository := mock_repository.NewMockBlogApp(mockController)
			testCase.mockBehaviorGetPostAuthorID(mockRepository, testCase.postID)
			s := service.NewGetService(mockRepository)

			// Act
			actualAuthorID, actualError := s.GetPostAuthorID(testCase.postID)

			// Assert
			assert.Equal(t, testCase.expectedAuthorID, actualAuthorID)

			if actualError != nil { // Если ошибка присутствует, то проверяем всё, кроме Trace.
				assert.Equal(t, testCase.expectedError.Error(), actualError.Error())
				assert.Equal(t, testCase.expectedError.GetErrCode(), actualError.GetErrCode())
				assert.Equal(t, testCase.expectedError.GetAltMsg(), actualError.GetAltMsg())
			} else {
				assert.Equal(t, testCase.expectedError, actualError)
			}
		})
	}
}

func TestService_EditPost(t *testing.T) {
	type mockBehaviorGetPost func(m *mock_repository.MockBlogApp, postID int)
	type mockBehaviorUpdatePost func(m *mock_repository.MockBlogApp, newPost *models.Post)

	testTable := []struct {
		name                   string
		newPost                *models.Post
		admin                  bool
		mockBehaviorGetPost    mockBehaviorGetPost
		readyToRecordPost      *models.Post
		mockBehaviorUpdatePost mockBehaviorUpdatePost
		expectedError          exterr.ErrExtender
	}{ // -----------------------------------------------START--------------------------------------------
		{ // -------------------------------------------Case #1-------------------------------------------
			name:    "[CONDITION]:USER,OK;[RESULT]:OK",
			newPost: &testPost,
			admin:   false,
			mockBehaviorGetPost: func(m *mock_repository.MockBlogApp, postID int) {
				m.EXPECT().
					GetPost(postID).
					Return(&testPost, nil)
			},
			readyToRecordPost: &testPost,
			mockBehaviorUpdatePost: func(m *mock_repository.MockBlogApp, newPost *models.Post) {
				m.EXPECT().
					UpdatePost(&testPost).
					Return(nil)
			},
			expectedError: nil,
		}, { // ----------------------------------------Case #2-------------------------------------------
			name:    "[CONDITION]:USER,OK;[RESULT]:Repository GetPost ERROR",
			newPost: &testPost,
			admin:   false,
			mockBehaviorGetPost: func(m *mock_repository.MockBlogApp, postID int) {
				m.EXPECT().
					GetPost(postID).
					Return(nil, exterr.New("Repository error"))
			},
			readyToRecordPost:      nil,
			mockBehaviorUpdatePost: func(m *mock_repository.MockBlogApp, newPost *models.Post) {},
			expectedError: exterr.NewWithExtErr("Service: EditPost post not found", exterr.New("Repository error")).
				SetErrCode(http.StatusNotFound).
				SetAltMsg("Post not found"),
		}, { // ----------------------------------------Case #3-------------------------------------------
			name: "[CONDITION]:USER,Empty title;[RESULT]:OK",
			newPost: &models.Post{
				Title:   "",
				Content: testPost.Content,
				Author:  testPost.Author,
			},
			admin: false,
			mockBehaviorGetPost: func(m *mock_repository.MockBlogApp, postID int) {
				m.EXPECT().
					GetPost(postID).
					Return(&testPost, nil)
			},
			readyToRecordPost: &testPost,
			mockBehaviorUpdatePost: func(m *mock_repository.MockBlogApp, newPost *models.Post) {
				m.EXPECT().
					UpdatePost(&testPost).
					Return(nil)
			},
			expectedError: nil,
		}, { // ----------------------------------------Case #4-------------------------------------------
			name: "[CONDITION]:USER,Empty content;[RESULT]:OK",
			newPost: &models.Post{
				Title:   testPost.Title,
				Content: "",
				Author:  testPost.Author,
			},
			admin: false,
			mockBehaviorGetPost: func(m *mock_repository.MockBlogApp, postID int) {
				m.EXPECT().
					GetPost(postID).
					Return(&testPost, nil)
			},
			readyToRecordPost: &testPost,
			mockBehaviorUpdatePost: func(m *mock_repository.MockBlogApp, newPost *models.Post) {
				m.EXPECT().
					UpdatePost(&testPost).
					Return(nil)
			},
			expectedError: nil,
		}, { // ----------------------------------------Case #5-------------------------------------------
			name: "[CONDITION]:USER,Change access level;[RESULT]:Access denied ERROR",
			newPost: &models.Post{
				Title:       testPost.Title,
				Content:     testPost.Content,
				Author:      testPost.Author,
				AccessLevel: 2,
			},
			admin: false,
			mockBehaviorGetPost: func(m *mock_repository.MockBlogApp, postID int) {
				m.EXPECT().
					GetPost(postID).
					Return(&testPost, nil)
			},
			readyToRecordPost:      &testPost,
			mockBehaviorUpdatePost: func(m *mock_repository.MockBlogApp, newPost *models.Post) {},
			expectedError: exterr.New("Service: access denied").
				SetErrCode(http.StatusForbidden).
				SetAltMsg("Only anministrator allowed to change post access level"),
		}, { // ----------------------------------------Case #6-------------------------------------------
			name: "[CONDITION]:ADMIN,Change access level;[RESULT]:OK",
			newPost: &models.Post{
				Title:       testPost.Title,
				Content:     testPost.Content,
				Author:      testPost.Author,
				AccessLevel: 2,
			},
			admin: true,
			mockBehaviorGetPost: func(m *mock_repository.MockBlogApp, postID int) {
				m.EXPECT().
					GetPost(postID).
					Return(&testPost, nil)
			},
			readyToRecordPost: &models.Post{
				Title:       testPost.Title,
				Content:     testPost.Content,
				Author:      testPost.Author,
				AccessLevel: 2,
			},
			mockBehaviorUpdatePost: func(m *mock_repository.MockBlogApp, newPost *models.Post) {
				m.EXPECT().
					UpdatePost(newPost).
					Return(nil)
			},
			expectedError: nil,
		}, { // ----------------------------------------Case #7-------------------------------------------
			name:    "[CONDITION]:USER,OK;[RESULT]:Repository UpdatePost ERROR",
			newPost: &testPost,
			admin:   false,
			mockBehaviorGetPost: func(m *mock_repository.MockBlogApp, postID int) {
				m.EXPECT().
					GetPost(postID).
					Return(&testPost, nil)
			},
			readyToRecordPost: &testPost,
			mockBehaviorUpdatePost: func(m *mock_repository.MockBlogApp, newPost *models.Post) {
				m.EXPECT().
					UpdatePost(newPost).
					Return(exterr.New("Repository error"))
			},
			expectedError: exterr.NewWithExtErr("Service: EditPost error", exterr.New("Repository error")).
				SetErrCode(http.StatusInternalServerError).
				SetAltMsg(http.StatusText(http.StatusInternalServerError)),
		},
	} // -------------------------------------------------END---------------------------------------------
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Arrange
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			mockRepository := mock_repository.NewMockBlogApp(mockController)
			testCase.mockBehaviorGetPost(mockRepository, testCase.newPost.ID)
			testCase.mockBehaviorUpdatePost(mockRepository, testCase.readyToRecordPost)
			s := service.NewGetService(mockRepository)

			// Act
			actualError := s.EditPost(testCase.newPost, testCase.admin)

			// Assert

			if actualError != nil { // Если ошибка присутствует, то проверяем всё, кроме Trace.
				assert.Equal(t, testCase.expectedError.Error(), actualError.Error())
				assert.Equal(t, testCase.expectedError.GetErrCode(), actualError.GetErrCode())
				assert.Equal(t, testCase.expectedError.GetAltMsg(), actualError.GetAltMsg())
			} else {
				assert.Equal(t, testCase.expectedError, actualError)
			}
		})
	}
}

func TestService_DeletePost(t *testing.T) {
	type mockBehaviorDeletePost func(m *mock_repository.MockBlogApp, postID int)

	testTable := []struct {
		name                   string
		postID                 int
		mockBehaviorDeletePost mockBehaviorDeletePost
		expectedError          exterr.ErrExtender
	}{ // -----------------------------------------------START--------------------------------------------
		{ // -------------------------------------------Case #1-------------------------------------------
			name:   "[CONDITION]:OK;[RESULT]:OK",
			postID: 1,
			mockBehaviorDeletePost: func(m *mock_repository.MockBlogApp, postID int) {
				m.EXPECT().
					DeletePost(postID).
					Return(nil)
			},
			expectedError: nil,
		}, { // ----------------------------------------Case #2-------------------------------------------
			name:   "[CONDITION]:OK;[RESULT]:Repository DeletePost ERROR",
			postID: 1,
			mockBehaviorDeletePost: func(m *mock_repository.MockBlogApp, postID int) {
				m.EXPECT().
					DeletePost(postID).
					Return(exterr.New("Repository error"))
			},
			expectedError: exterr.NewWithExtErr("Service: DeletePost error", exterr.New("Repository error")).
				SetErrCode(http.StatusInternalServerError).
				SetAltMsg(http.StatusText(http.StatusInternalServerError)),
		},
	} // -------------------------------------------------END---------------------------------------------
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			// Arrange
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			mockRepository := mock_repository.NewMockBlogApp(mockController)
			testCase.mockBehaviorDeletePost(mockRepository, testCase.postID)
			s := service.NewGetService(mockRepository)

			// Act
			actualError := s.DeletePost(testCase.postID)

			// Assert

			if actualError != nil { // Если ошибка присутствует, то проверяем всё, кроме Trace.
				assert.Equal(t, testCase.expectedError.Error(), actualError.Error())
				assert.Equal(t, testCase.expectedError.GetErrCode(), actualError.GetErrCode())
				assert.Equal(t, testCase.expectedError.GetAltMsg(), actualError.GetAltMsg())
			} else {
				assert.Equal(t, testCase.expectedError, actualError)
			}
		})
	}
}

func TestService_LikePost(T *testing.T) {
	type mockBehaviorGetCheckLikePost func(m *mock_repository.MockBlogApp, userID, postID int)
	type mockBehaviorGetAmountLikePost func(m *mock_repository.MockBlogApp, postID int)
	type mockBehaviorLikePost func(m *mock_repository.MockBlogApp, userID, postID, amountLike int)

	testTable := []struct {
		name                          string                        // Название теста
		userID                        int                           // ID пользователя
		postID                        int                           // ID статьи(поста)
		amountLike                    int                           // Количество лайков у поста
		mockBehaviorGetCheckLikePost  mockBehaviorGetCheckLikePost  // mock-функция проверки лайка пользователя у поста
		mockBehaviorGetAmountLikePost mockBehaviorGetAmountLikePost // mock-функция получения количества лайков у поста
		mockBehaviorLikePost          mockBehaviorLikePost          // mock-функция занесения в БД лайка поста от пользователя
		expectedError                 exterr.ErrExtender            // возвращаемая ошибка
	}{ // -----------------------------------------------START--------------------------------------------
		{ // -------------------------------------------Case #1-------------------------------------------
			name: "[CONDITION]:OK;[RESULT]:OK",
			// Input data
			userID:     45,
			postID:     150,
			amountLike: 10,
			// Mock behavior
			mockBehaviorGetCheckLikePost: func(m *mock_repository.MockBlogApp, userID, postID int) {
				m.EXPECT().
					GetCheckLikePost(userID, postID).
					Return(exterr.New("there are no matches in the database"))
			},
			mockBehaviorGetAmountLikePost: func(m *mock_repository.MockBlogApp, postID int) {
				m.EXPECT().
					GetAmountLikePost(postID).
					Return(10, nil)
			},
			mockBehaviorLikePost: func(m *mock_repository.MockBlogApp, userID, postID, amountLike int) {
				m.EXPECT().
					LikePost(userID, postID, amountLike).
					Return(nil)
			},
			// Expected output
			expectedError: nil,
		}, { // ----------------------------------------Case #2-------------------------------------------
			name: "[CONDITION]:OK;[RESULT]:GetCheckLikePost post already liked ERROR",
			// Input data
			userID:     45,
			postID:     150,
			amountLike: 10,
			// Mock behavior
			mockBehaviorGetCheckLikePost: func(m *mock_repository.MockBlogApp, userID, postID int) {
				m.EXPECT().
					GetCheckLikePost(userID, postID).
					Return(nil)
			},
			mockBehaviorGetAmountLikePost: func(m *mock_repository.MockBlogApp, postID int) {},
			mockBehaviorLikePost:          func(m *mock_repository.MockBlogApp, userID, postID, amountLike int) {},
			// Expected output
			expectedError: exterr.New("error for checkLikePost").
				SetErrCode(http.StatusBadRequest).
				SetAltMsg("User liked has already put a like"),
		}, { // ----------------------------------------Case #3-------------------------------------------
			name: "[CONDITION]:OK;[RESULT]:GetAmountLike repository ERROR",
			// Input data
			userID:     45,
			postID:     150,
			amountLike: 10,
			// Mock behavior
			mockBehaviorGetCheckLikePost: func(m *mock_repository.MockBlogApp, userID, postID int) {
				m.EXPECT().
					GetCheckLikePost(userID, postID).
					Return(exterr.New("there are no matches in the database"))
			},
			mockBehaviorGetAmountLikePost: func(m *mock_repository.MockBlogApp, postID int) {
				m.EXPECT().
					GetAmountLikePost(postID).
					Return(0, exterr.New("Repository: GetAmountLikePost post not found"))
			},
			mockBehaviorLikePost: func(m *mock_repository.MockBlogApp, userID, postID, amountLike int) {},
			// Expected output
			expectedError: exterr.NewWithExtErr("get error amount likes post", exterr.New("Repository: GetAmountLikePost post not found")).
				SetErrCode(http.StatusBadRequest).
				SetAltMsg("get error amount likes post"),
		}, { // ----------------------------------------Case #4-------------------------------------------
			name: "[CONDITION]:OK;[RESULT]:LikePost repository ERROR",
			// Input data
			userID:     45,
			postID:     150,
			amountLike: 10,
			// Mock behavior
			mockBehaviorGetCheckLikePost: func(m *mock_repository.MockBlogApp, userID, postID int) {
				m.EXPECT().
					GetCheckLikePost(userID, postID).
					Return(exterr.New("there are no matches in the database"))
			},
			mockBehaviorGetAmountLikePost: func(m *mock_repository.MockBlogApp, postID int) {
				m.EXPECT().
					GetAmountLikePost(postID).
					Return(10, nil)
			},
			mockBehaviorLikePost: func(m *mock_repository.MockBlogApp, userID, postID, amountLike int) {
				m.EXPECT().
					LikePost(userID, postID, amountLike).
					Return(exterr.New("Repository: LikePost INSERT INTO liked_post table error"))
			},
			// Expected output
			expectedError: exterr.NewWithExtErr("get error for like post", exterr.New("Repository: LikePost INSERT INTO liked_post table error")).
				SetErrCode(http.StatusBadRequest).
				SetAltMsg("get error for like post"),
		},
	} // -------------------------------------------------END---------------------------------------------
	for _, testCase := range testTable {
		T.Run(testCase.name, func(t *testing.T) {
			// Arrange
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			mockRepository := mock_repository.NewMockBlogApp(mockController)

			testCase.mockBehaviorGetCheckLikePost(mockRepository, testCase.userID, testCase.postID)
			testCase.mockBehaviorGetAmountLikePost(mockRepository, testCase.postID)
			testCase.mockBehaviorLikePost(mockRepository, testCase.userID, testCase.postID, testCase.amountLike)

			s := service.NewGetService(mockRepository)

			// Act
			actualError := s.LikePost(testCase.userID, testCase.postID)

			// Assert
			if actualError != nil { // Если ошибка присутствует, то проверяем всё, кроме Trace.
				assert.Equal(t, testCase.expectedError.Error(), actualError.Error())
				assert.Equal(t, testCase.expectedError.GetErrCode(), actualError.GetErrCode())
				assert.Equal(t, testCase.expectedError.GetAltMsg(), actualError.GetAltMsg())
			} else {
				assert.Equal(t, testCase.expectedError, actualError)
			}
		})
	}
}

func TestService_UnlikePost(T *testing.T) {
	type mockBehaviorGetCheckLikePost func(m *mock_repository.MockBlogApp, userID, postID int)
	type mockBehaviorGetAmountLikePost func(m *mock_repository.MockBlogApp, postID int)
	type mockBehaviorUnlikePost func(m *mock_repository.MockBlogApp, userID, postID, amountLike int)

	testTable := []struct {
		name                          string                        // Название теста
		userID                        int                           // ID пользователя
		postID                        int                           // ID статьи(поста)
		amountLike                    int                           // Количество лайков у поста
		mockBehaviorGetCheckLikePost  mockBehaviorGetCheckLikePost  // mock-функция проверки лайка пользователя у поста
		mockBehaviorGetAmountLikePost mockBehaviorGetAmountLikePost // mock-функция получения количества лайков у поста
		mockBehaviorUnlikePost        mockBehaviorUnlikePost        // mock-функция занесения в БД лайка поста от пользователя
		expectedError                 exterr.ErrExtender            // возвращаемая ошибка
	}{ // -----------------------------------------------START--------------------------------------------
		{ // -------------------------------------------Case #1-------------------------------------------
			name: "[CONDITION]:OK;[RESULT]:OK",
			// Input data
			userID:     45,
			postID:     150,
			amountLike: 10,
			// Mock behavior
			mockBehaviorGetCheckLikePost: func(m *mock_repository.MockBlogApp, userID, postID int) {
				m.EXPECT().
					GetCheckLikePost(userID, postID).
					Return(nil)
			},
			mockBehaviorGetAmountLikePost: func(m *mock_repository.MockBlogApp, postID int) {
				m.EXPECT().
					GetAmountLikePost(postID).
					Return(10, nil)
			},
			mockBehaviorUnlikePost: func(m *mock_repository.MockBlogApp, userID, postID, amountLike int) {
				m.EXPECT().
					UnlikePost(userID, postID, amountLike).
					Return(nil)
			},
			// Expected output
			expectedError: nil,
		}, { // ----------------------------------------Case #2-------------------------------------------
			name: "[CONDITION]:OK;[RESULT]:GetCheckLikePost post was not liked ERROR",
			// Input data
			userID:     45,
			postID:     150,
			amountLike: 10,
			// Mock behavior
			mockBehaviorGetCheckLikePost: func(m *mock_repository.MockBlogApp, userID, postID int) {
				m.EXPECT().
					GetCheckLikePost(userID, postID).
					Return(exterr.New("there are no matches in the database"))
			},
			mockBehaviorGetAmountLikePost: func(m *mock_repository.MockBlogApp, postID int) {},
			mockBehaviorUnlikePost:        func(m *mock_repository.MockBlogApp, userID, postID, amountLike int) {},
			// Expected output
			expectedError: exterr.NewWithErr("error for checkLikePost", exterr.New("there are no matches in the database")).
				SetErrCode(http.StatusBadRequest).
				SetAltMsg("the user did not like the post"),
		}, { // ----------------------------------------Case #3-------------------------------------------
			name: "[CONDITION]:OK;[RESULT]:GetAmountLike repository ERROR",
			// Input data
			userID:     45,
			postID:     150,
			amountLike: 10,
			// Mock behavior
			mockBehaviorGetCheckLikePost: func(m *mock_repository.MockBlogApp, userID, postID int) {
				m.EXPECT().
					GetCheckLikePost(userID, postID).
					Return(nil)
			},
			mockBehaviorGetAmountLikePost: func(m *mock_repository.MockBlogApp, postID int) {
				m.EXPECT().
					GetAmountLikePost(postID).
					Return(0, exterr.New("Repository: GetAmountLikePost post not found"))
			},
			mockBehaviorUnlikePost: func(m *mock_repository.MockBlogApp, userID, postID, amountLike int) {},
			// Expected output
			expectedError: exterr.NewWithExtErr("get error amount likes post", exterr.New("Repository: GetAmountLikePost post not found")).
				SetErrCode(http.StatusBadRequest).
				SetAltMsg("get error amount likes post"),
		}, { // ----------------------------------------Case #4-------------------------------------------
			name: "[CONDITION]:OK;[RESULT]:LikePost repository ERROR",
			// Input data
			userID:     45,
			postID:     150,
			amountLike: 10,
			// Mock behavior
			mockBehaviorGetCheckLikePost: func(m *mock_repository.MockBlogApp, userID, postID int) {
				m.EXPECT().
					GetCheckLikePost(userID, postID).
					Return(nil)
			},
			mockBehaviorGetAmountLikePost: func(m *mock_repository.MockBlogApp, postID int) {
				m.EXPECT().
					GetAmountLikePost(postID).
					Return(10, nil)
			},
			mockBehaviorUnlikePost: func(m *mock_repository.MockBlogApp, userID, postID, amountLike int) {
				m.EXPECT().
					UnlikePost(userID, postID, amountLike).
					Return(exterr.New("update error like in post"))
			},
			// Expected output
			expectedError: exterr.NewWithExtErr("get error for like post", exterr.New("update error like in post")).
				SetErrCode(http.StatusBadRequest).
				SetAltMsg("get error for like post"),
		},
	} // -------------------------------------------------END---------------------------------------------
	for _, testCase := range testTable {
		T.Run(testCase.name, func(t *testing.T) {
			// Arrange
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			mockRepository := mock_repository.NewMockBlogApp(mockController)

			testCase.mockBehaviorGetCheckLikePost(mockRepository, testCase.userID, testCase.postID)
			testCase.mockBehaviorGetAmountLikePost(mockRepository, testCase.postID)
			testCase.mockBehaviorUnlikePost(mockRepository, testCase.userID, testCase.postID, testCase.amountLike)

			s := service.NewGetService(mockRepository)

			// Act
			actualError := s.UnlikePost(testCase.userID, testCase.postID)

			// Assert
			if actualError != nil { // Если ошибка присутствует, то проверяем всё, кроме Trace.
				assert.Equal(t, testCase.expectedError.Error(), actualError.Error())
				assert.Equal(t, testCase.expectedError.GetErrCode(), actualError.GetErrCode())
				assert.Equal(t, testCase.expectedError.GetAltMsg(), actualError.GetAltMsg())
			} else {
				assert.Equal(t, testCase.expectedError, actualError)
			}
		})
	}
}

func TestService_AddFavoritPost(T *testing.T) {
	type mockBehaviorGetCheckFavoritesPost func(m *mock_repository.MockBlogApp, userID, postID int)
	type mockBehaviorGetAmountFavoritesPost func(m *mock_repository.MockBlogApp, postID int)
	type mockBehaviorFavoritesPost func(m *mock_repository.MockBlogApp, userID, postID, amountFavorites int)

	testTable := []struct {
		name                               string                             // Название теста
		userID                             int                                // ID пользователя
		postID                             int                                // ID статьи(поста)
		amountFavorites                    int                                // Количество добавленного в избранное у поста
		mockBehaviorGetCheckFavoritesPost  mockBehaviorGetCheckFavoritesPost  // mock-функция проверки добавления в избранное пользователем поста
		mockBehaviorGetAmountFavoritesPost mockBehaviorGetAmountFavoritesPost // mock-функция получения количества добавленного в избранное у поста
		mockBehaviorFavoritesPost          mockBehaviorFavoritesPost          // mock-функция занесения в БД добавления в избранное поста от пользователя
		expectedError                      exterr.ErrExtender                 // возвращаемая ошибка
	}{ // -----------------------------------------------START--------------------------------------------
		{ // -------------------------------------------Case #1-------------------------------------------
			name: "[CONDITION]:OK;[RESULT]:OK",
			// Input data
			userID:          45,
			postID:          150,
			amountFavorites: 10,
			// Mock behavior
			mockBehaviorGetCheckFavoritesPost: func(m *mock_repository.MockBlogApp, userID, postID int) {
				m.EXPECT().
					GetCheckFavoritesPost(userID, postID).
					Return(exterr.New("there are no matches in the database"))
			},
			mockBehaviorGetAmountFavoritesPost: func(m *mock_repository.MockBlogApp, postID int) {
				m.EXPECT().
					GetAmountFavoritesPost(postID).
					Return(10, nil)
			},
			mockBehaviorFavoritesPost: func(m *mock_repository.MockBlogApp, userID, postID, amountFavorites int) {
				m.EXPECT().
					FavoritesPost(userID, postID, amountFavorites).
					Return(nil)
			},
			// Expected output
			expectedError: nil,
		},
		{ // ----------------------------------------Case #2-------------------------------------------
			name: "[CONDITION]:OK;[RESULT]:GetCheckFavoritesPost already in favorites ERROR",
			// Input data
			userID:          45,
			postID:          150,
			amountFavorites: 10,
			// Mock behavior
			mockBehaviorGetCheckFavoritesPost: func(m *mock_repository.MockBlogApp, userID, postID int) {
				m.EXPECT().
					GetCheckFavoritesPost(userID, postID).
					Return(nil)
			},
			mockBehaviorGetAmountFavoritesPost: func(m *mock_repository.MockBlogApp, postID int) {},
			mockBehaviorFavoritesPost:          func(m *mock_repository.MockBlogApp, userID, postID, amountFavorites int) {},
			// Expected output
			expectedError: exterr.New("error for FavoritesLikePost").
				SetErrCode(http.StatusBadRequest).
				SetAltMsg("user favorit has already put a like"),
		}, { // ----------------------------------------Case #3-------------------------------------------
			name: "[CONDITION]:OK;[RESULT]:GetAmountFavoritesPost reository ERROR",
			// Input data
			userID:          45,
			postID:          150,
			amountFavorites: 10,
			// Mock behavior
			mockBehaviorGetCheckFavoritesPost: func(m *mock_repository.MockBlogApp, userID, postID int) {
				m.EXPECT().
					GetCheckFavoritesPost(userID, postID).
					Return(exterr.New("there are no matches in the database"))
			},
			mockBehaviorGetAmountFavoritesPost: func(m *mock_repository.MockBlogApp, postID int) {
				m.EXPECT().
					GetAmountFavoritesPost(postID).
					Return(0, exterr.New("Repository: GetAmountFavoritesPost post not found"))
			},
			mockBehaviorFavoritesPost: func(m *mock_repository.MockBlogApp, userID, postID, amountFavorites int) {},
			// Expected output
			expectedError: exterr.NewWithExtErr("get error amount favorites post", exterr.New("Repository: GetAmountFavoritesPost post not found")).
				SetErrCode(http.StatusBadRequest).
				SetAltMsg("get error amount favorites post"),
		}, { // ----------------------------------------Case #4-------------------------------------------
			name: "[CONDITION]:OK;[RESULT]:FavoritesPost reository ERROR",
			// Input data
			userID:          45,
			postID:          150,
			amountFavorites: 10,
			// Mock behavior
			mockBehaviorGetCheckFavoritesPost: func(m *mock_repository.MockBlogApp, userID, postID int) {
				m.EXPECT().
					GetCheckFavoritesPost(userID, postID).
					Return(exterr.New("there are no matches in the database"))
			},
			mockBehaviorGetAmountFavoritesPost: func(m *mock_repository.MockBlogApp, postID int) {
				m.EXPECT().
					GetAmountFavoritesPost(postID).
					Return(10, nil)
			},
			mockBehaviorFavoritesPost: func(m *mock_repository.MockBlogApp, userID, postID, amountFavorites int) {
				m.EXPECT().
					FavoritesPost(userID, postID, amountFavorites).
					Return(exterr.New("Repository: FavoritesPost UPDATE post table error"))
			},
			// Expected output
			expectedError: exterr.NewWithExtErr("get error for favorit post", exterr.New("Repository: FavoritesPost UPDATE post table error")).
				SetErrCode(http.StatusBadRequest).
				SetAltMsg("get error for favorit post"),
		},
	} // -------------------------------------------------END---------------------------------------------
	for _, testCase := range testTable {
		T.Run(testCase.name, func(t *testing.T) {
			// Arrange
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			mockRepository := mock_repository.NewMockBlogApp(mockController)

			testCase.mockBehaviorGetCheckFavoritesPost(mockRepository, testCase.userID, testCase.postID)
			testCase.mockBehaviorGetAmountFavoritesPost(mockRepository, testCase.postID)
			testCase.mockBehaviorFavoritesPost(mockRepository, testCase.userID, testCase.postID, testCase.amountFavorites)

			s := service.NewGetService(mockRepository)

			// Act
			actualError := s.AddFavoritPost(testCase.userID, testCase.postID)

			// Assert
			if actualError != nil { // Если ошибка присутствует, то проверяем всё, кроме Trace.
				assert.Equal(t, testCase.expectedError.Error(), actualError.Error())
				assert.Equal(t, testCase.expectedError.GetErrCode(), actualError.GetErrCode())
				assert.Equal(t, testCase.expectedError.GetAltMsg(), actualError.GetAltMsg())
			} else {
				assert.Equal(t, testCase.expectedError, actualError)
			}
		})
	}
}

func TestService_UnfavoritesPost(T *testing.T) {
	type mockBehaviorGetCheckFavoritesPost func(m *mock_repository.MockBlogApp, userID, postID int)
	type mockBehaviorGetAmountFavoritesPost func(m *mock_repository.MockBlogApp, postID int)
	type mockBehaviorUnfavoritesPost func(m *mock_repository.MockBlogApp, userID, postID, amountFavorites int)

	testTable := []struct {
		name                               string                             // Название теста
		userID                             int                                // ID пользователя
		postID                             int                                // ID статьи(поста)
		amountFavorites                    int                                // Количество добавленного в избранное у поста
		mockBehaviorGetCheckFavoritesPost  mockBehaviorGetCheckFavoritesPost  // mock-функция проверки добавления в избранное пользователем поста
		mockBehaviorGetAmountFavoritesPost mockBehaviorGetAmountFavoritesPost // mock-функция получения количества добавленного в избранное у поста
		mockBehaviorUnfavoritesPost        mockBehaviorUnfavoritesPost        // mock-функция занесения в БД добавления в избранное поста от пользователя
		expectedError                      exterr.ErrExtender                 // возвращаемая ошибка
	}{ // -----------------------------------------------START--------------------------------------------
		{ // -------------------------------------------Case #1-------------------------------------------
			name: "[CONDITION]:OK;[RESULT]:OK",
			// Input data
			userID:          45,
			postID:          150,
			amountFavorites: 10,
			// Mock behavior
			mockBehaviorGetCheckFavoritesPost: func(m *mock_repository.MockBlogApp, userID, postID int) {
				m.EXPECT().
					GetCheckFavoritesPost(userID, postID).
					Return(nil)
			},
			mockBehaviorGetAmountFavoritesPost: func(m *mock_repository.MockBlogApp, postID int) {
				m.EXPECT().
					GetAmountFavoritesPost(postID).
					Return(10, nil)
			},
			mockBehaviorUnfavoritesPost: func(m *mock_repository.MockBlogApp, userID, postID, amountFavorites int) {
				m.EXPECT().
					UnfavoritesPost(userID, postID, amountFavorites).
					Return(nil)
			},
			// Expected output
			expectedError: nil,
		},
		{ // ----------------------------------------Case #2-------------------------------------------
			name: "[CONDITION]:OK;[RESULT]:GetCheckFavoritesPost post not added to favorites ERROR",
			// Input data
			userID:          45,
			postID:          150,
			amountFavorites: 10,
			// Mock behavior
			mockBehaviorGetCheckFavoritesPost: func(m *mock_repository.MockBlogApp, userID, postID int) {
				m.EXPECT().
					GetCheckFavoritesPost(userID, postID).
					Return(exterr.New("there are no matches in the database"))
			},
			mockBehaviorGetAmountFavoritesPost: func(m *mock_repository.MockBlogApp, postID int) {},
			mockBehaviorUnfavoritesPost:        func(m *mock_repository.MockBlogApp, userID, postID, amountFavorites int) {},
			// Expected output
			expectedError: exterr.NewWithExtErr("error for checkFavoritesPost", exterr.New("there are no matches in the database")).
				SetErrCode(http.StatusBadRequest).
				SetAltMsg("the user did not Favorites the post"),
		}, { // ----------------------------------------Case #3-------------------------------------------
			name: "[CONDITION]:OK;[RESULT]:GetAmountFavoritesPost repository ERROR",
			// Input data
			userID:          45,
			postID:          150,
			amountFavorites: 10,
			// Mock behavior
			mockBehaviorGetCheckFavoritesPost: func(m *mock_repository.MockBlogApp, userID, postID int) {
				m.EXPECT().
					GetCheckFavoritesPost(userID, postID).
					Return(nil)
			},
			mockBehaviorGetAmountFavoritesPost: func(m *mock_repository.MockBlogApp, postID int) {
				m.EXPECT().
					GetAmountFavoritesPost(postID).
					Return(0, exterr.New("Repository: GetAmountFavoritesPost post not found"))
			},
			mockBehaviorUnfavoritesPost: func(m *mock_repository.MockBlogApp, userID, postID, amountFavorites int) {},
			// Expected output
			expectedError: exterr.NewWithExtErr("get error amount Favorites post", exterr.New("Repository: GetAmountFavoritesPost post not found")).
				SetErrCode(http.StatusBadRequest).
				SetAltMsg("get error amount Favorites post"),
		}, { // ----------------------------------------Case #4-------------------------------------------
			name: "[CONDITION]:OK;[RESULT]:UnfavoritesPost repository ERROR",
			// Input data
			userID:          45,
			postID:          150,
			amountFavorites: 10,
			// Mock behavior
			mockBehaviorGetCheckFavoritesPost: func(m *mock_repository.MockBlogApp, userID, postID int) {
				m.EXPECT().
					GetCheckFavoritesPost(userID, postID).
					Return(nil)
			},
			mockBehaviorGetAmountFavoritesPost: func(m *mock_repository.MockBlogApp, postID int) {
				m.EXPECT().
					GetAmountFavoritesPost(postID).
					Return(10, nil)
			},
			mockBehaviorUnfavoritesPost: func(m *mock_repository.MockBlogApp, userID, postID, amountFavorites int) {
				m.EXPECT().
					UnfavoritesPost(userID, postID, amountFavorites).
					Return(exterr.New("update error unfavorites post"))
			},
			// Expected output
			expectedError: exterr.NewWithExtErr("get error for Favorites post", exterr.New("update error unfavorites post")).
				SetErrCode(http.StatusBadRequest).
				SetAltMsg("get error for Favorites post"),
		},
	} // -------------------------------------------------END---------------------------------------------
	for _, testCase := range testTable {
		T.Run(testCase.name, func(t *testing.T) {
			// Arrange
			mockController := gomock.NewController(t)
			defer mockController.Finish()

			mockRepository := mock_repository.NewMockBlogApp(mockController)

			testCase.mockBehaviorGetCheckFavoritesPost(mockRepository, testCase.userID, testCase.postID)
			testCase.mockBehaviorGetAmountFavoritesPost(mockRepository, testCase.postID)
			testCase.mockBehaviorUnfavoritesPost(mockRepository, testCase.userID, testCase.postID, testCase.amountFavorites)

			s := service.NewGetService(mockRepository)

			// Act
			actualError := s.UnfavoritesPost(testCase.userID, testCase.postID)

			// Assert
			if actualError != nil { // Если ошибка присутствует, то проверяем всё, кроме Trace.
				assert.Equal(t, testCase.expectedError.Error(), actualError.Error())
				assert.Equal(t, testCase.expectedError.GetErrCode(), actualError.GetErrCode())
				assert.Equal(t, testCase.expectedError.GetAltMsg(), actualError.GetAltMsg())
			} else {
				assert.Equal(t, testCase.expectedError, actualError)
			}
		})
	}
}

// func TestService_(t *testing.T) {
// 	type mockBehavior func(m *mock_repository.MockBlogApp)
// 	testTable := []struct {
// 		name                       string
// 		expectedError                        exterr.ErrExtender
// 	}{ // -----------------------------------------------START--------------------------------------------
// 		{ // -------------------------------------------Case #1-------------------------------------------
// 			name: "Happy path (OK)",
// 			expectedError: nil,
// 		}, { // ----------------------------------------Case #2-------------------------------------------
// 			name: "",
// 			expectedError: nil,
// 		},
// 	} // -------------------------------------------------END---------------------------------------------
// 	for _, testCase := range testTable {
// 		t.Run(testCase.name, func(t *testing.T) {
// 			// Arrange
// 			mockController := gomock.NewController(t)
// 			defer mockController.Finish()
// 			mockRepository := mock_repository.NewMockBlogApp(mockController)
// 			s := service.NewGetService(mockRepository)
// 			// Act
//			actualError := nil
// 			// Assert
// 			if actualError != nil { // Если ошибка присутствует, то проверяем всё, кроме Trace.
// 				assert.Equal(t, testCase.expectedError.Error(), actualError.Error())
// 				assert.Equal(t, testCase.expectedError.GetErrCode(), actualError.GetErrCode())
// 				assert.Equal(t, testCase.expectedError.GetAltMsg(), actualError.GetAltMsg())
// 			} else {
// 				assert.Equal(t, testCase.expectedError, actualError)
// 			}
// 		})
// 	}
// }
