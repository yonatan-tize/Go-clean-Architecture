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