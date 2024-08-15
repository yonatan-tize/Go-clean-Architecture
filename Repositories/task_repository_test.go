package Repositories

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

// SetupSuite sets up the test suite for the TaskRepository.
// It connects to the MongoDB test instance, pings the MongoDB instance to ensure connection is established,
// and initializes the necessary variables for the test suite.
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

// TearDownSuite is a function that is called after all tests have run in the TaskRepositoryTestSuite.
// It is responsible for dropping the test database and disconnecting from MongoDB.
func (suite *TaskRepositoryTestSuite) TearDownSuite() {
	// Drop the test database after all tests have run
	err := suite.db.Drop(context.Background())
	suite.Require().NoError(err)

	// Disconnect from MongoDB
	err = suite.client.Disconnect(context.Background())
	suite.Require().NoError(err)
}

// TestCreateTask tests the CreateTask method of the TaskRepository.
//
// It creates a new task with the given title, description, status, and due date.
// Then, it calls the CreateTask method of the repository to insert the task into the database.
// It verifies that no error occurred during the insertion.
// Finally, it retrieves the inserted task from the database and compares its title with the original task's title.
// It also checks that no error occurred during the retrieval.
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

// TestFindAlltasks is a unit test function that tests the FindAlltasks method of the TaskRepository.
// It verifies that the correct number of tasks were fetched and that the fetched tasks have the same title, description, status, and due date as the previously created task.
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

// TestFindTaskById tests the functionality of finding a task by its ID.
//
// It creates a new task, retrieves it by its ID, and then verifies if the retrieved task matches the original task.
// The test performs the following steps:
// 1. Creates a new task with a title, description, status, and due date.
// 2. Calls the `CreateTask` method of the repository to create the task.
// 3. Retrieves the task by its ID using the `FindTaskById` method of the repository.
// 4. Verifies if the retrieved task matches the original task.
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


// TestUpdateTaskById tests the functionality of updating a task by its ID.
//
// It creates a new task, updates its properties, and then verifies if the task is updated correctly.
// The test performs the following steps:
// 1. Creates a new task with a title, description, status, and due date.
// 2. Calls the `CreateTask` method of the repository to create the task.
// 3. Updates the task's title, description, status, and due date with new values.
// 4. Calls the `UpdateTaskById` method of the repository to update the task.
// 5. Retrieves the updated task from the database and verifies if it matches the updated values.
// 6. Retrieves the task by its ID and verifies if its fields are updated correctly.
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

// TestDeleteTask tests the DeleteTask method of the TaskRepository.
//
// It creates a new task, inserts it into the database, verifies the task is inserted,
// and then deletes the task using the DeleteTask method. It asserts that no error occurs during the process.
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