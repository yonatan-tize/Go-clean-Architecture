package controllers

import (
	"bytes"
	"encoding/json"
	"example/go-clean-architecture/Domain"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	// "github.com/stretchr/testify/suite"
)

// TestSuite struct that holds all the mocks and controllers
// type TestSuite struct {
// 	suite.Suite
// 	mockUserUseCase *mocks.UserUseCase
// 	userController  UserController
// }

// SetupTest initializes the test suite before each test
// func (suite *TestSuite) SetupTest() {
// 	suite.mockUserUseCase = new(mocks.UserUseCase)
// 	suite.userController = UserController{
// 		UserUseCase: suite.mockUserUseCase,
// 	}
// }

// TestCreateAccount tests the CreateAccount method
func (suite *TestSuite) TestCreateAccount() {
	// Mock data
	newUser := Domain.User{
		Username: "testuser",
		Password: "123456789",
		Role: "USER",
	}

	// Mock the CreateAccount method
	suite.mockUserUseCase.On("CreateAccount", mock.Anything, &newUser).Return(newUser, nil)

	// Create a new gin context
	gin.SetMode(gin.TestMode)
	jsonValue, _ := json.Marshal(newUser)
	req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	suite.userController.CreateAccount(c)

	// Assert the status code and response
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.mockUserUseCase.AssertExpectations(suite.T())
}

// TestLogin tests the Login method
func (suite *TestSuite) TestLogin() {
	// Mock data
	user := Domain.User{
		Username: "testuser",
		Password: "password",
	}
	mockToken := "mockToken"

	// Mock the AuthenticateUser method
	suite.mockUserUseCase.On("AuthenticateUser", mock.Anything, user.Username, user.Password).Return(user, mockToken, nil)

	// Create a new gin context
	gin.SetMode(gin.TestMode)
	jsonValue, _ := json.Marshal(user)
	req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	suite.userController.Login(c)

	// Assert the status code and response
	assert.Equal(suite.T(), http.StatusAccepted, w.Code)
	suite.mockUserUseCase.AssertExpectations(suite.T())
}

// TestPromoteUser tests the PromoteUser method
func (suite *TestSuite) TestPromoteUser() {
	// Mock the UpdateUserRole method
	suite.mockUserUseCase.On("UpdateUserRole", mock.Anything, "1").Return(nil)

	// Create a new gin context
	gin.SetMode(gin.TestMode)
	req, _ := http.NewRequest(http.MethodPost, "/promote/1", nil)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: "1"}}

	suite.userController.PromoteUser(c)

	// Assert the status code and response
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.mockUserUseCase.AssertExpectations(suite.T())
}

// Run the test suite
// func TestControllerSuite(t *testing.T) {
// 	suite.Run(t, new(TestSuite))
// }
