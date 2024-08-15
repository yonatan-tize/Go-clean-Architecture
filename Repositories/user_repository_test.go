package Repositories

import (
	"context"
	"fmt"
	"testing"

	domain "example/go-clean-architecture/Domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userName = "johndoe"
var password = "password123"
var userRole = "ADMIN"

type UserRepositoryTestSuite struct {
	suite.Suite
	client     *mongo.Client
	db         *mongo.Database
	collection string
	repo       domain.UserRepository
}

// SetupSuite is a setup function for the User Repository test suite.
// It connects to the MongoDB test instance, pings the MongoDB instance to ensure connection is established,
// and initializes the necessary variables for the test suite.
func (suite *UserRepositoryTestSuite) SetupSuite() {
	// Connect to the MongoDB test instance
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	suite.Require().NoError(err)

	// Ping the MongoDB instance to ensure connection is established
	err = client.Ping(context.Background(), nil)
	suite.Require().NoError(err)

	suite.client = client
	suite.db = client.Database("testdb")
	suite.collection = "users"
	suite.repo = NewUserRepository(*suite.db, suite.collection)
}

// TearDownSuite is a function that is called after all tests have run to clean up resources.
// It drops the test database and disconnects from MongoDB.
func (suite *UserRepositoryTestSuite) TearDownSuite() {
	// Drop the test database after all tests have run
	err := suite.db.Drop(context.Background())
	suite.Require().NoError(err)

	// Disconnect from MongoDB
	err = suite.client.Disconnect(context.Background())
	suite.Require().NoError(err)
}


// TestCreateUser tests the CreateUser method of the UserRepository.
// It creates a new user with the given username, password, and role.
// The user is then inserted into the database and verified.
// The test ensures that the user is successfully created and can be found in the database.
// If any error occurs during the test, it will be reported.
func (suite *UserRepositoryTestSuite) TestCreateUser() {

	newUser := domain.User{
		Username:   userName,
		Password:   password,
		Role:  		userRole,
	}

	user, err := suite.repo.CreateNewUser(context.Background(), &newUser)
	suite.Require().NoError(err)

	// Verify user was inserted
	collection := suite.db.Collection(suite.collection)
	var result domain.User

	err = collection.FindOne(context.Background(), bson.M{"username": user.Username}).Decode(&result)
	suite.Require().NoError(err)
	// fmt.Println(user)
	// fmt.Println(result)


	assert.Equal(suite.T(), user.Username, result.Username)
	suite.NoError(err, "no error when creating and finding the same user")
}


// TestFindUser tests the FindUser method of the UserRepository.
// It fetches a user from the repository and verifies that the fetched user's username and role match the previously inserted values.
func (suite *UserRepositoryTestSuite) TestFindUser() {
	fetchedUser, err := suite.repo.FindUser(context.Background(), userName)
	suite.Require().NoError(err)

	// Verify that the correct number of users were fetched
	assert.Equal(suite.T(), userName, fetchedUser.Username, "user name same with the previous inserted")
	assert.Equal(suite.T(), userRole, fetchedUser.Role, "user role same with the previous inserted")
}


// TestPromoteUser tests the functionality of promoting a user.
//
// It creates a new user with the given username, password, and role.
// Then it inserts the user into the database collection.
// It retrieves the inserted user from the database and verifies that it was inserted successfully.
// Next, it calls the PromoteUser method of the repository with the user's ID.
// Finally, it retrieves the updated user from the database and verifies that the role has changed.
// If any error occurs during the test, it fails the test.
func (suite *UserRepositoryTestSuite) TestPromoteUser() {
	newUser := domain.User{
		Username: "sister",
		Password: "password123",
		Role : "USER",
	}

	// Insert the user
	collection := suite.db.Collection(suite.collection)
    _, err := collection.InsertOne(context.Background(), newUser)
    suite.Require().NoError(err)

	var user domain.User
	err = collection.FindOne(context.Background(), bson.M{"username": newUser.Username}).Decode(&user)
    suite.Require().NoError(err)
	fmt.Println(user.ID.Hex())


	err = suite.repo.PromoteUser(context.Background(), user.ID.Hex())
	suite.Require().NoError(err)


	var foundUser domain.User
	err = collection.FindOne(context.Background(), bson.M{"username": newUser.Username}).Decode(&foundUser)
	suite.Require().NoError(err)

	assert.NotEqual(suite.T(), user.Role, foundUser.Role, "role changed")
}

func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}