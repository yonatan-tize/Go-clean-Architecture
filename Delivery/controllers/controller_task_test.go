package controllers

import (
	"bytes"
	"encoding/json"
	"example/go-clean-architecture/Domain"
	"example/go-clean-architecture/mocks"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var taskID = primitive.NewObjectID()
var taskTitle = "Task Title"
var taskDescription = "Task Description"
var taskStatus = "Started"
var taskDueDate = time.Now()

// TestSuite struct that holds all the mocks and controllers
type TestSuite struct {
	suite.Suite
	mockUserUseCase *mocks.UserUseCase
	mockTaskUseCase *mocks.TaskUseCase
	userController  UserController
	taskController  TaskController
}

// SetupTest initializes the test suite before each test
func (suite *TestSuite) SetupTest() {
	suite.mockUserUseCase = new(mocks.UserUseCase)
	suite.mockTaskUseCase = new(mocks.TaskUseCase)
	suite.userController = UserController{
		UserUseCase: suite.mockUserUseCase,
	}
	suite.taskController = TaskController{
		TaskUseCase: suite.mockTaskUseCase,
	}
}

// TestGetTasks tests the GetTasks method
func (suite *TestSuite) TestGetTasks() {
	// Mock data
	tasks := []Domain.Task{
		{
			ID:          taskID,
			Title:       taskTitle,
			Description: taskDescription,
			Status:      taskStatus,
			DueDate:    taskDueDate,
		},
		{
			ID:           primitive.NewObjectID(),
			Title:       "second task",
			Description: "description for second task",
			Status:      "ongoing",
			DueDate:    taskDueDate,
		},
	}

	// Mock the GetAllTasks method
	suite.mockTaskUseCase.On("GetAllTasks", mock.Anything).Return(tasks, nil)

	// Create a new gin context
	gin.SetMode(gin.TestMode)
	req, _ := http.NewRequest(http.MethodGet, "/tasks", nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	suite.taskController.GetTasks(c)

	// Assert the status code and response
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.mockTaskUseCase.AssertExpectations(suite.T())
}

// TestGetTask tests the GetTask method
func (suite *TestSuite) TestGetTask() {
	// Mock data
	task := Domain.Task{
		ID:          taskID,
		Title:       taskTitle,
		Description: taskDescription,
		Status:      taskStatus,
		DueDate:    taskDueDate,
	}

	// Mock the GetTaskByID method
	suite.mockTaskUseCase.On("GetTaskByID", mock.Anything, taskID.Hex()).Return(task, nil)

	// Create a new gin context
	gin.SetMode(gin.TestMode)
	req, _ := http.NewRequest(http.MethodGet, "/tasks/" + taskID.Hex(), nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: taskID.Hex()}}

	suite.taskController.GetTask(c)

	// Assert the status code and response
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.mockTaskUseCase.AssertExpectations(suite.T())
}

// TestCreateTask tests the CreateTask method
func (suite *TestSuite) TestCreateTask() {
	// Mock data

	newTask := Domain.Task{
		Title: "New Task", 
		Description: "This is a new task",
		DueDate:  time.Date(2024, time.August, 15, 21, 19, 31, 703415278, time.UTC),
	}
	createdTask := Domain.Task{
		ID: primitive.NewObjectID(), 
		Title: "New Task", 
		Description: "This is a new task",
		DueDate:  time.Date(2024, time.August, 15, 21, 19, 31, 703415278, time.UTC),
	}

	// Mock the AddNewTask method
	suite.mockTaskUseCase.On("AddNewTask", mock.AnythingOfType("*gin.Context"), newTask).Return(createdTask, nil)

	// Create a new gin context
	gin.SetMode(gin.TestMode)
	jsonValue, _ := json.Marshal(newTask)
	req, _ := http.NewRequest(http.MethodPost, "/tasks", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	suite.taskController.CreateTask(c)

	// Assert the status code and response
	assert.Equal(suite.T(), http.StatusCreated, w.Code)
	suite.mockTaskUseCase.AssertExpectations(suite.T())
}

// TestUpdateTask tests the UpdateTask method
func (suite *TestSuite) TestUpdateTask() {
	// Mock data
	updatedTask := Domain.Task{Title: "Updated Task", Description: "This task is updated"}

	// Mock the ModifyTaskById method
	suite.mockTaskUseCase.On("ModifyTaskById", mock.Anything, updatedTask,  taskID.Hex()).Return(updatedTask, nil)

	// Create a new gin context
	gin.SetMode(gin.TestMode)
	jsonValue, _ := json.Marshal(updatedTask)
	req, _ := http.NewRequest(http.MethodPut, "/tasks/" +  taskID.Hex(), bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value:  taskID.Hex()}}

	suite.taskController.UpdateTask(c)

	// Assert the status code and response
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.mockTaskUseCase.AssertExpectations(suite.T())
}

// TestDeleteTask tests the DeleteTask method
func (suite *TestSuite) TestDeleteTask() {
	// Mock the DeleteTaskById method
	suite.mockTaskUseCase.On("DeleteTaskById", mock.Anything, taskID.Hex()).Return(nil)

	// Create a new gin context
	gin.SetMode(gin.TestMode)
	req, _ := http.NewRequest(http.MethodDelete, "/tasks/" + taskID.Hex(), nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "id", Value: taskID.Hex()}}

	suite.taskController.DeleteTask(c)

	// Assert the status code and response
	assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.mockTaskUseCase.AssertExpectations(suite.T())
}

// Run the test suite
func TestControllerSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}
