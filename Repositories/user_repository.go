package Repositories

import (
	"context"
	"errors"
	domain "example/go-clean-architecture/Domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	database   mongo.Database
	collection string
}
var _ domain.UserRepository = &userRepository{}
func NewUserRepository(db mongo.Database, collection string) *userRepository {
	return &userRepository{
		database:   db,
		collection: collection,
	}
}
// FindUser(ctx context.Context, username string) User
// 	CreateNewUser(ctx context.Context, user User) error
// 	PromoteUser(ctx context.Context, userId string) error
func (ur *userRepository) FindUser(ctx context.Context, username string) (domain.User, error){
	collection := ur.database.Collection(ur.collection)
	var existingUser domain.User
	err := collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&existingUser)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return domain.User{}, errors.New("user not found")
		}
		return domain.User{}, err
	}
	return existingUser, nil
}

// CreateNewUser creates a new user in the database.
// It takes a context and a user object as input parameters.
// It returns the created user object and an error if any.
func (ur *userRepository) CreateNewUser(ctx context.Context, user *domain.User) (domain.User, error){
	var existingUser domain.User
	collection := ur.database.Collection(ur.collection)
	err := collection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&existingUser)
	if err == nil {
		return domain.User{}, errors.New("username already exists")
	}

	// count users in the database. 
	count, _ := collection.CountDocuments(context.TODO(), bson.M{})

	// Promote the first user to admin
	if count == 0 {
		user.Role = "ADMIN"
	}else {
		user.Role = "USER"
	}

	user.ID = primitive.NewObjectID()	

	_, err = collection.InsertOne(context.Background(), user)
	if err != nil{
		return domain.User{}, errors.New("failed to insert data")
	}
	return *user, nil 
}


func (ur *userRepository) PromoteUser(ctx context.Context, userId string) error{
	collection := ur.database.Collection(ur.collection)
	objID, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err // Return an error if userId is not a valid ObjectID
	}
	filter := bson.M{"_id": objID}
	update := bson.M{
		"$set": bson.M{
			"role": "ADMIN",
		},
	}
	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return mongo.ErrNoDocuments // No user found with the given ID
	}
	return nil
}


