package repositories

import (
	"context"
	"testing"
	"time"

	domain "example/go-clean-architecture/Domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var title = "New Task"
var description = "This is a new task"
var status = "Started"
var due_date = time.Now().UTC().Truncate(24 * time.Hour)

type TaskRepositoryTestSuite struct {
	suite.Suite
	client     *mongo.Client
	db         *mongo.Database
	collection string
	repo       domain.TaskRepository
}

func (suite *TaskRepositoryTestSuite) SetupSuite() {
	// Connect to the MongoDB test instance
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.Background(), clientOptions)
	suite.Require().NoError(err)

	// Ping the MongoDB instance to ensure connection is established
	err = client.Ping(context.Background(), nil)
	suite.Require().NoError(err)

	suite.client = client
	suite.db = client.Database("testTaskManager")
	suite.collection = "tasks"
	suite.repo = NewTaskRepository(*suite.db, suite.collection)
}

func (suite *TaskRepositoryTestSuite) TearDownSuite() {
	// Drop the test database after all tests have run
	err := suite.db.Drop(context.Background())
	suite.Require().NoError(err)

	// Disconnect from MongoDB
	err = suite.client.Disconnect(context.Background())
	suite.Require().NoError(err)
}

func (suite *TaskRepositoryTestSuite) TestCreateTask() {
	task := domain.Task{
		Title:       title,
		Description: description,
		Status:      status,
		DueDate:     due_date,
	}

	newtask, err := suite.repo.CreateTask(context.Background(), task)
	suite.Require().NoError(err)

	// Verify task wask inserted
	collection := suite.db.Collection(suite.collection)
	var result domain.Task

	err = collection.FindOne(context.Background(), bson.M{"_id": newtask.ID}).Decode(&result)
	suite.Require().NoError(err)
	assert.Equal(suite.T(), task.Title, result.Title)
	suite.NoError(err, "no error when creating and finding the same task")
}

func (suite *TaskRepositoryTestSuite) TestFindAlltasks() {
	fetchedTasks, err := suite.repo.FindAlltasks(context.Background())
	suite.Require().NoError(err)

	// Verify the correct number of tasks were fetched
	tasks := fetchedTasks[0]
	// assert.Equal(suite.T(), 1, len(fetchedTasks))

	
	assert.Equal(suite.T(), title, tasks.Title, "Same task title  with the previously created task")
	assert.Equal(suite.T(), description, tasks.Description, "Same task description with the previously created task")
	assert.Equal(suite.T(), status, tasks.Status, "Same task status with the previously created task")
	assert.Equal(suite.T(), due_date, tasks.DueDate, "Same task due_date with the previously created task")
	

}

func (suite *TaskRepositoryTestSuite) TestFindTaskById() {
	task := domain.Task{
		Title: "New Task",
		Description: "This is a new task",
		Status: "Started",
		DueDate: time.Now().UTC().Truncate(24 * time.Hour),
	}
	

	newTask, err := suite.repo.CreateTask(context.Background(), task)
	suite.Require().NoError(err)

	foundTask, err := suite.repo.FindTaskById(context.Background(), newTask.ID.Hex())
	suite.Require().NoError(err)

	assert.Equal(suite.T(), task.Title, foundTask.Title, "Same task title with the previously created task")
	assert.Equal(suite.T(), task.Description, foundTask.Description, "Same task description with the previously created task")
	assert.Equal(suite.T(), task.Status, foundTask.Status, "Same task status with the previously created task")
	assert.Equal(suite.T(), task.DueDate, foundTask.DueDate, "Same task due_date with the previously created task")
}


func (suite *TaskRepositoryTestSuite) TestUpdateTaskById() {
	task := domain.Task{
		Title: "New Task",
		Description: "This is a new task",
		Status: "Started",
		DueDate: time.Now().UTC().Truncate(24 * time.Hour),
	}
	

	newTask, err := suite.repo.CreateTask(context.Background(), task)
	suite.Require().NoError(err)

	newTitle := "Updated Task"
	newDescription := "This is an updated task"
	newStatus := "Completed"
	newDueDate := time.Now().UTC().Truncate(24 * time.Hour)

	taskUpdate := domain.Task{
		Title:       newTitle,
		Description: newDescription,
		Status:      newStatus,
		DueDate:    newDueDate,
	}

	task, err = suite.repo.UpdateTaskById(context.Background(), taskUpdate, newTask.ID.Hex())
	suite.Require().NoError(err)

	// find out if the task is updated
	collection := suite.db.Collection(suite.collection)
	var result domain.Task

	err = collection.FindOne(context.Background(), bson.M{"_id": task.ID}).Decode(&result)
	suite.Require().NoError(err)

	// Check if its the same task
	taskFound, err := suite.repo.FindTaskById(context.Background(), task.ID.Hex())
	suite.Require().NoError(err)

	// check if the fields are updated
	assert.Equal(suite.T(), newTitle, taskFound.Title, "Same updated task title with the previously updated task")
	assert.Equal(suite.T(), newDescription, taskFound.Description, "Same updated task description with the previously updated task")
	assert.Equal(suite.T(), newStatus, taskFound.Status, "Same updated task status with the previously updated task")
	assert.Equal(suite.T(), newDueDate, taskFound.DueDate, "Same updated task due_date with the previously updated task")
}

func (suite *TaskRepositoryTestSuite) TestDeleteTask() {
	

	// create a new task then remove it
	task := domain.Task{
		Title:       "new task delete",
		Description: "This is a new task to be deleted",
		Status:      "Not Started",
		DueDate:    time.Now().UTC().Truncate(24 * time.Hour),
	}

	newTask, err := suite.repo.CreateTask(context.Background(), task)
	suite.Require().NoError(err)

	// Verify the task is inserted
	collection := suite.db.Collection(suite.collection)
	var result domain.Task

	err = collection.FindOne(context.Background(), bson.M{"_id": newTask.ID}).Decode(&result)
	suite.Require().NoError(err)

	err = suite.repo.DeleteTask(context.Background(), newTask.ID.Hex())
	suite.Require().NoError(err)

}

func TestTaskRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(TaskRepositoryTestSuite))
}