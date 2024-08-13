package usecases

import (
	"context"
	"errors"
	"time"

	domain "example/go-clean-architecture/Domain"
	infrastructure "example/go-clean-architecture/Infrastructure"
	"github.com/go-playground/validator/v10"
)

type userUseCase struct {
	userRepository domain.UserRepository
	contextTimeout time.Duration
}

var _ domain.UserUseCase = &userUseCase{}

func NewUserUsecase(userRepository domain.UserRepository, timeout time.Duration) domain.UserUseCase {
	return &userUseCase{
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

var validate = validator.New()


func (ur *userUseCase) AuthenticateUser(c context.Context, userName string, password string) (domain.User, string, error) {

	ctx, close := context.WithTimeout(c, ur.contextTimeout)
	defer close()
	user, err := ur.userRepository.FindUser(ctx, userName)
	if err != nil{
		return domain.User{}, "", errors.New("wrong password")
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

// CreateAccount implements domain.UserUseCase.
func (ur *userUseCase) CreateAccount(c context.Context, user *domain.User) (domain.User, error) {
	ctx, close := context.WithTimeout(c, ur.contextTimeout)
	defer close()

	err := validate.Struct(user)
	if err != nil {
		return domain.User{}, err
	}
	//hash the password and send to database
	hashedPassword, err := infrastructure.HashPassword(user.Password)
	if err != nil {
		return domain.User{}, err
	}
	user.Password = hashedPassword
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	return ur.userRepository.CreateNewUser(ctx, *user)
}

// UpdateUserRole implements domain.UserUseCase.
func (ur *userUseCase) UpdateUserRole(c context.Context, userId string) error {
	// panic("unimplemented")
	ctx, close := context.WithTimeout(c, ur.contextTimeout)
	defer close()
	return ur.userRepository.PromoteUser(ctx, userId)
}