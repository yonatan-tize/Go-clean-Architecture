package usecases

import (
	"context"
	"errors"
	"time"

	domain "example/go-clean-architecture/Domain"
	infrastructure "example/go-clean-architecture/Infrastructure"
)

// userUseCase represents the use case for managing user entities.
type userUseCase struct {
	userRepository domain.UserRepository
	contextTimeout time.Duration
}

var _ domain.UserUseCase = &userUseCase{}

// NewUserUsecase creates a new instance of the UserUseCase interface.
// It takes a userRepository of type domain.UserRepository and a timeout of type time.Duration as parameters.
// It returns a pointer to a userUseCase struct that implements the UserUseCase interface.
func NewUserUsecase(userRepository domain.UserRepository, timeout time.Duration) domain.UserUseCase {
	return &userUseCase{
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

// AuthenticateUser authenticates a user by verifying their username and password.
// It takes a context.Context, userName string, and password string as input parameters.
// It returns a domain.User, a token string, and an error.
// The domain.User represents the authenticated user.
// The token string is a generated token for the authenticated user.
// The error is returned if there is an issue with the authentication process.
func (ur *userUseCase) AuthenticateUser(c context.Context, userName string, password string) (domain.User, string, error) {
	ctx, close := context.WithTimeout(c, ur.contextTimeout)
	defer close()
	
	user, err := ur.userRepository.FindUser(ctx, userName)
	if err != nil{
		return domain.User{}, "", errors.New("user not found")
	}

	//verify the password
	isValidPassword := infrastructure.VerifyPassword(password, user.Password)
	if !isValidPassword {
		return domain.User{}, "", errors.New("wrong password")
	}
	//generate the token
	token, err := infrastructure.GenerateToken(user)
	if err != nil {
		return domain.User{}, "", err
	}
	return user, token, nil
}

// CreateAccount creates a new user account.
// It takes a context.Context and a *domain.User as input parameters.
// The function hashes the user's password and sends it to the database.
// It returns the created domain.User and an error if any.
func (ur *userUseCase) CreateAccount(c context.Context, user *domain.User) (domain.User, error) {
	ctx, close := context.WithTimeout(c, ur.contextTimeout)
	defer close()

	//hash the password and send to database
	hashedPassword, err := infrastructure.HashPassword(user.Password)
	if err != nil {
		return domain.User{}, err
	}
	user.Password = hashedPassword
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return ur.userRepository.CreateNewUser(ctx, user)
}

// UpdateUserRole updates the role of a user identified by the given userId.
// It takes a context.Context as the first argument and the userId as the second argument.
// It returns an error if the operation fails.
func (ur *userUseCase) UpdateUserRole(c context.Context, userId string) error {
	ctx, close := context.WithTimeout(c, ur.contextTimeout)
	defer close()
	return ur.userRepository.PromoteUser(ctx, userId)
}