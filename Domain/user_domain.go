package Domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Username  string             `bson:"username" json:"username" validate:"required"`
	Password  string             `bson:"password" json:"password" validate:"required"`
	Role      string             `bson:"role" json:"role" validate:"required,eq=ADMIN|eq=USER"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type UserRepository interface {
	FindUser(ctx context.Context, userId string) (User, error)
	CreateNewUser(ctx context.Context, user *User) (User, error)
	PromoteUser(ctx context.Context, userId string) error
}

type UserUseCase interface {
	CreateAccount(ctx context.Context, user *User) (User, error)
	AuthenticateUser(ctx context.Context, userName string, password string) (User, string, error)
	UpdateUserRole(ctx context.Context, id string) error
}
