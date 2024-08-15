package repositories

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

func (suite *UserRepositoryTestSuite) TearDownSuite() {
	// Drop the test database after all tests have run
	err := suite.db.Drop(context.Background())
	suite.Require().NoError(err)

	// Disconnect from MongoDB
	err = suite.client.Disconnect(context.Background())
	suite.Require().NoError(err)
}


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


func (suite *UserRepositoryTestSuite) TestFindUser() {
	fetchedUser, err := suite.repo.FindUser(context.Background(), userName)
	suite.Require().NoError(err)

	// Verify that the correct number of users were fetched
	assert.Equal(suite.T(), userName, fetchedUser.Username, "user name same with the previous inserted")
	assert.Equal(suite.T(), userRole, fetchedUser.Role, "user role same with the previous inserted")
}


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