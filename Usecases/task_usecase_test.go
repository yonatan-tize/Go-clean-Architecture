package usecases

import (
	"context"
	domain "example/go-clean-architecture/Domain"
	"example/go-clean-architecture/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskUseCaseSuite struct {
	suite.Suite
	mockTaskRepo *mocks.TaskRepository
	taskUseCase  *taskUseCase
}

var taskID = primitive.NewObjectID()
var taskTitle = "Task Title"
var taskDescription = "Task Description"
var taskStatus = "Started"
var taskDueDate = time.Now()

func (suite *TaskUseCaseSuite) SetupTest() {
	suite.mockTaskRepo = new(mocks.TaskRepository)
	suite.taskUseCase = NewTaskUsecase(suite.mockTaskRepo, time.Second*2)
}

// TestGetAllTasks is a unit test function that tests the GetAllTasks method of the TaskUseCase struct.
// It creates mock tasks and asserts that the returned tasks match the expected tasks.
// The mockTaskRepo is expected to return the mock tasks without any error.
// This test ensures that GetAllTasks method retrieves all tasks successfully.
func (suite *TaskUseCaseSuite) TestGetAllTasks() {
	mockTasks := []domain.Task{
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

	suite.mockTaskRepo.On("FindAlltasks", mock.Anything).Return(mockTasks, nil)

	tasks, err := suite.taskUseCase.GetAllTasks(context.Background())
	assert.NoError(suite.T(), err)

	for idx, mockTask:= range mockTasks{
		assert.Equal(suite.T(), mockTask, tasks[idx])

	}
	suite.mockTaskRepo.AssertExpectations(suite.T())
}

// TestAddNewTask is a unit test function that tests the AddNewTask method of the TaskUseCase struct.
// It creates a new task with the given taskID, taskTitle, taskDescription, taskStatus, and taskDueDate.
// The mockTaskRepo's CreateTask method is expected to be called with the newTask parameter and return the createdTask and nil error.
// The AddNewTask method is then called with the newTask parameter, and the returned task and error are asserted.
// Finally, the expectations of the mockTaskRepo are asserted.
func (suite *TaskUseCaseSuite) TestAddNewTask() {
	newTask := domain.Task{
			ID:          taskID,
			Title:       taskTitle,
			Description: taskDescription,
			Status:      taskStatus,
			DueDate:    taskDueDate,
	}
	createdTask := domain.Task{
		ID:          taskID,
		Title:       taskTitle,
		Description: taskDescription,
		Status:      taskStatus,
		DueDate:     taskDueDate,
	}

	suite.mockTaskRepo.On("CreateTask", mock.Anything, newTask).Return(createdTask, nil)

	task, err := suite.taskUseCase.AddNewTask(context.Background(), newTask)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), createdTask, task)
	suite.mockTaskRepo.AssertExpectations(suite.T())
}

// TestGetTaskByID is a unit test function that tests the GetTaskByID method of the TaskUseCase struct.
// It creates a mock task and expects the task repository to return the same task when FindTaskById is called.
// The function then calls the GetTaskByID method and asserts that the returned task is equal to the expected task.
// Finally, it asserts that all expectations on the mock task repository have been met.
func (suite *TaskUseCaseSuite) TestGetTaskByID() {
    task := domain.Task{
        ID:          taskID,
        Title:       "Test Task",
        Description: "Test Description",
        Status:      "pending",
        DueDate:     time.Now(),
    }

    suite.mockTaskRepo.On("FindTaskById", mock.Anything, taskID.Hex()).Return(task, nil)

	result, err := suite.taskUseCase.GetTaskByID(context.Background(), taskID.Hex())

    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), task, result)

    suite.mockTaskRepo.AssertExpectations(suite.T())
}

// TestModifyTaskById is a unit test function that tests the ModifyTaskById method of the TaskUseCase struct.
// It verifies that the method correctly modifies a task by its ID and returns the updated task.
// The test creates an updated task with a new title, description, status, and due date.
// It then sets up a mock task repository to expect the UpdateTaskById method to be called with the updated task and the task ID.
// The method is called with the updated task and task ID, and the returned result and error are asserted.
// Finally, the expectations of the mock task repository are asserted.
func (suite *TaskUseCaseSuite) TestModifyTaskById() {
    updatedTask := domain.Task{
        ID:          taskID,
        Title:       "Updated Task",
        Description: "Updated Description",
        Status:      "completed",
        DueDate:     time.Now(),
    }

    suite.mockTaskRepo.On("UpdateTaskById", mock.Anything, updatedTask, taskID.Hex()).Return(updatedTask, nil)

    result, err := suite.taskUseCase.ModifyTaskById(context.Background(), updatedTask, taskID.Hex())

    assert.NoError(suite.T(), err)
    assert.Equal(suite.T(), updatedTask, result)
    suite.mockTaskRepo.AssertExpectations(suite.T())
}

// TestDeleteTask is a unit test function that tests the DeleteTaskById method of the TaskUseCase struct.
// It creates a new task ID using primitive.NewObjectID().Hex() and mocks the DeleteTask method of the task repository.
// The test asserts that no error is returned when calling the DeleteTaskById method and that the expectations of the mocked repository are met.
func (suite *TaskUseCaseSuite) TestDeleteTask() {
    taskID := primitive.NewObjectID().Hex()

    suite.mockTaskRepo.On("DeleteTask", mock.Anything, taskID).Return(nil)

    err := suite.taskUseCase.DeleteTaskById(context.Background(), taskID)

    assert.NoError(suite.T(), err)
    suite.mockTaskRepo.AssertExpectations(suite.T())
}

func TestTaskUseCaseSuite(t *testing.T) {
	suite.Run(t, new(TaskUseCaseSuite))
}