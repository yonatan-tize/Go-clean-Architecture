package usecases

import (
	"context"
	"errors"
	"testing"
	"time"

	"example/go-clean-architecture/Domain"
	infrastructure "example/go-clean-architecture/Infrastructure"
	"example/go-clean-architecture/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserUseCaseSuite struct {
	suite.Suite
	mockUserRepo *mocks.UserRepository
	userUseCase  Domain.UserUseCase
}

var userName = "johndoe"
var password = "password123"
var userRole = "ADMIN"


func (suite *UserUseCaseSuite) SetupTest() {
	suite.mockUserRepo = new(mocks.UserRepository)
	suite.userUseCase = NewUserUsecase(suite.mockUserRepo, time.Second*2)
}

// TestAuthenticateUser_Success tests the successful authentication of a user.
//
// It arranges the necessary dependencies and sets up the test environment.
// - It hashes the password using the infrastructure.HashPassword function.
// - It creates a mock user with a username, role, and hashed password.
// - It mocks the FindUser method of the user repository to return the mock user.
//
// It verifies the validity of the password by comparing it with the hashed password using the infrastructure.VerifyPassword function.
//
// It calls the AuthenticateUser method of the user use case with the username and password.
//
// It asserts that there is no error returned from the AuthenticateUser method.
// It asserts that the username of the returned user matches the username of the mock user.
// It asserts that all expectations of the mock user repository are met.
func (suite *UserUseCaseSuite) TestAuthenticateUser_Success() {
	// Arrange
	hashedPassword, err := infrastructure.HashPassword(password)
	assert.Equal(suite.T(), nil, err, "failed to hash password")

	mockUser := Domain.User{
		ID: primitive.NewObjectID(),
		Username: userName,
		Role: userRole, 
	}

	mockUser.Password = hashedPassword
	suite.mockUserRepo.On("FindUser", mock.Anything, "johndoe").Return(mockUser, nil)

	
	isValid := infrastructure.VerifyPassword(password, hashedPassword)
	assert.Equal(suite.T(), true, isValid)

	user, _, err := suite.userUseCase.AuthenticateUser(context.Background(), "johndoe", password)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), mockUser.Username, user.Username)
	suite.mockUserRepo.AssertExpectations(suite.T())
}

// TestAuthenticateUser_UserNotFound tests the scenario where the user is not found during authentication.
//
// It arranges a mock user repository to return an error when trying to find the user with the given username.
// Then it calls the AuthenticateUser method of the user use case with the username and password.
// Finally, it asserts that an error occurred and the error message matches the expected "user not found".
//
// This test ensures that the user use case handles the case when the user is not found during authentication.
func (suite *UserUseCaseSuite) TestAuthenticateUser_UserNotFound() {
	// Arrange
	suite.mockUserRepo.On("FindUser", mock.Anything, "johndoe").Return(Domain.User{}, errors.New("user not found"))

	// Act
	_, _, err := suite.userUseCase.AuthenticateUser(context.Background(), "johndoe", "password")

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "user not found", err.Error())
	suite.mockUserRepo.AssertExpectations(suite.T())
}

// TestAuthenticateUser_WrongPassword tests the behavior of the AuthenticateUser method when an incorrect password is provided.
//
// It performs the following steps:
// 1. Creates a mock user with a hashed password.
// 2. Sets up the mock user repository to return the mock user when FindUser is called.
// 3. Verifies the validity of the provided password against the hashed password.
// 4. Calls the AuthenticateUser method with the username and a wrong password.
// 5. Asserts that an error occurred and the error message is "wrong password".
// 6. Asserts that all expectations on the mock user repository were met.
func (suite *UserUseCaseSuite) TestAuthenticateUser_WrongPassword() {
	// Arrange
	mockUser := Domain.User{
		Username: userName,
		Password: password,
		Role: userRole, 
	}
	hashedPassword, err := infrastructure.HashPassword(password)
	assert.Equal(suite.T(), nil, err)

	mockUser.Password = hashedPassword
	suite.mockUserRepo.On("FindUser", mock.Anything, "johndoe").Return(mockUser, nil)

	
	isValid := infrastructure.VerifyPassword(password, hashedPassword)
	assert.Equal(suite.T(), true, isValid)

	wrongPassword := "wrongpassword"
	
	_, _, err = suite.userUseCase.AuthenticateUser(context.Background(), "johndoe", wrongPassword)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), "wrong password", err.Error())
	suite.mockUserRepo.AssertExpectations(suite.T())
}

// TestCreateAccount_Success tests the successful creation of a user account.
//
// It arranges a mock user with the given username, password, and role.
// Then it mocks the CreateNewUser method of the user repository to return the mock user and no error.
//
// The test acts by calling the CreateAccount method of the user use case with the mock user.
// It expects no error to occur during the account creation.
//
// Finally, it asserts that the username of the returned user matches the given username.
// It also asserts that all expectations of the mock user repository are met.
func (suite *UserUseCaseSuite) TestCreateAccount_Success() {
	// Arrange
	mockUser := &Domain.User{
		Username: userName,
		Password: password,
		Role: userRole, 
	}

	suite.mockUserRepo.On("CreateNewUser", mock.Anything, mockUser).Return(*mockUser, nil)	

	// Act
	user, err := suite.userUseCase.CreateAccount(context.Background(), mockUser)
	assert.NoError(suite.T(), err)


	// Assert
	assert.Equal(suite.T(), userName, user.Username)
	suite.mockUserRepo.AssertExpectations(suite.T())
}

// TestUpdateUserRole_Success tests the successful update of a user's role.
//
// It sets up a dummy user ID and mocks the PromoteUser method of the user repository to return nil.
// Then, it calls the UpdateUserRole method of the user use case with the dummy user ID.
// Finally, it asserts that no error occurred during the update and verifies that all expectations on the mocked user repository were met.
func (suite *UserUseCaseSuite) TestUpdateUserRole_Success() {
	// Arrange
	userId := "dummyId"
	suite.mockUserRepo.On("PromoteUser", mock.Anything, userId).Return(nil)

	// Act
	err := suite.userUseCase.UpdateUserRole(context.Background(), userId)

	// Assert
	assert.NoError(suite.T(), err)
	suite.mockUserRepo.AssertExpectations(suite.T())
}

// TestUpdateUserRole_Failure tests the failure scenario of updating a user's role.
//
// It sets up a mock user repository to return an error when promoting a user.
// Then, it calls the UpdateUserRole method of the user use case with a user ID.
// Finally, it asserts that an error is returned.
//
// This test ensures that the user use case handles the failure scenario correctly when updating a user's role.
func (suite *UserUseCaseSuite) TestUpdateUserRole_Failure() {
	// Arrange
	userId := "some-id"
	suite.mockUserRepo.On("PromoteUser", mock.Anything, userId).Return(errors.New("some error"))

	// Act
	err := suite.userUseCase.UpdateUserRole(context.Background(), userId)

	// Assert
	assert.Error(suite.T(), err)
	suite.mockUserRepo.AssertExpectations(suite.T())
}



func TestUserUseCaseSuite(t *testing.T) {
	suite.Run(t, new(UserUseCaseSuite))
}
