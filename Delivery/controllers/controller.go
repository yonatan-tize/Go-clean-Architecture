package controllers

import (
	domain "example/go-clean-architecture/Domain"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type TaskController struct {
	TaskUseCase domain.TaskUseCase
}

type UserController struct {
	UserUseCase domain.UserUseCase
}
var validate = validator.New()

func (uc *UserController) CreateAccount(c *gin.Context) {

	var newUser domain.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := validate.Struct(newUser)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	user, err := uc.UserUseCase.CreateAccount(c, &newUser)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (uc *UserController) Login(c *gin.Context) {
	// by using username and password
	var user domain.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, token, err := uc.UserUseCase.AuthenticateUser(c, user.Username, user.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return 
	}
	c.JSON(http.StatusAccepted, gin.H{"token": token, "user": user})
}

func (uc *UserController) PromoteUser(c *gin.Context) {
	userId := c.Param("id")
	err := uc.UserUseCase.UpdateUserRole(c, userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "promoted to admin"})
}

// GetTasks retrieves all tasks from the database.
// It returns a JSON response with the fetched tasks.
// If there is an error fetching the tasks, it returns a JSON response with an error message.
func (tc *TaskController) GetTasks(c *gin.Context) {
	tasks, err := tc.TaskUseCase.GetAllTasks(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to fetch tasks"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "the data is", "tasks" :tasks })
}

// GetTask retrieves a task by its ID.
//
// Parameters:
// - c: The gin context.
//
// Returns:
// - task: The retrieved task.
// - err: An error if the task retrieval fails.
func (tc *TaskController) GetTask(c *gin.Context) {
	id := c.Param("id")
	task, err := tc.TaskUseCase.GetTaskByID(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}

// CreateTask handles the HTTP POST request to create a new task.
// It binds the JSON data from the request body to a newTask variable.
// If the JSON binding fails, it returns a JSON response with a bad request status code and an error message.
// Otherwise, it calls the AddNewTask method of the TaskUseCase to add the new task.
// If an error occurs during the task creation, it returns a JSON response with an internal server error status code and an error message.
// Finally, it returns a JSON response with a created status code and the created task.
func (tc *TaskController) CreateTask(c *gin.Context) {
	var newTask domain.Task

	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	createdTask, err := tc.TaskUseCase.AddNewTask(c, newTask)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, createdTask)
}

// UpdateTask updates a task with the specified ID.
//
// Parameters:
// - c: The gin context.
//
// Returns:
// - None.
func (tc *TaskController) UpdateTask(c *gin.Context) {

	id := c.Param("id")

	var newTask domain.Task
	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	res, err := tc.TaskUseCase.ModifyTaskById(c, newTask, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// DeleteTask deletes a task by its ID.
//
// Parameters:
// - c: The gin context.
// - id: The ID of the task to be deleted.
//
// Returns:
// - An error if the task deletion fails.
func (tc *TaskController) DeleteTask(c *gin.Context) {
	id := c.Param("id")

	err := tc.TaskUseCase.DeleteTaskById(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted successfully"})
}
