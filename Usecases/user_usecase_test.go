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
