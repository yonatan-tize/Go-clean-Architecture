package usecases

import (
	"context"
	"time"
	domain "example/go-clean-architecture/Domain"
)


// taskUseCase represents the use case for managing tasks.
type taskUseCase struct {
	taskRepository domain.TaskRepository
	contextTimeout time.Duration
}
var _ domain.TaskUseCase = &taskUseCase{}
func NewTaskUsecase(taskRepository domain.TaskRepository, timeout time.Duration) *taskUseCase {
	return &taskUseCase{
		taskRepository: taskRepository,
		contextTimeout: timeout,
	}
}

// GetAllTasks retrieves all tasks from the task repository.
// It takes a context as input and returns a slice of domain.Task and an error.
// The context is used to control the execution timeout.
func (tu *taskUseCase) GetAllTasks(c context.Context) ([]domain.Task, error){
	ctx, close := context.WithTimeout(c, tu.contextTimeout)
	defer close()

	return tu.taskRepository.FindAlltasks(ctx)
}

// GetTaskByID retrieves a task by its ID.
// It takes a context.Context and a taskId string as parameters.
// It returns a domain.Task and an error.
// The context.Context is used for managing the execution context of the function.
// The taskId is the unique identifier of the task to be retrieved.
// The function returns the found task and any error that occurred during the retrieval process.
func (tu *taskUseCase) GetTaskByID(c context.Context, taskId string) (domain.Task, error){
	ctx, close := context.WithTimeout(c, tu.contextTimeout)
	defer close()

	return tu.taskRepository.FindTaskById(ctx , taskId )
}


// AddNewTask adds a new task to the system.
// It takes a context and a task as input parameters and returns the created task and an error (if any).
func (tu *taskUseCase) AddNewTask(c context.Context, task domain.Task) (domain.Task, error){
    ctx, close := context.WithTimeout(c, tu.contextTimeout)
	defer close()

	return tu.taskRepository.CreateTask(ctx, task)
}


// ModifyTaskById modifies a task by its ID.
// It takes a context.Context, a task domain.Task, and a taskId string as parameters.
// It returns the modified task and an error, if any.
func (tu *taskUseCase) ModifyTaskById(c context.Context, task domain.Task,  taskId string) (domain.Task, error){
	ctx, close := context.WithTimeout(c, tu.contextTimeout)
	defer close()

	return tu.taskRepository.UpdateTaskById(ctx, task, taskId)
}


// DeleteTaskById deletes a task by its ID.
// It takes a context.Context and a taskId string as parameters.
// It returns an error if the task deletion fails.
func (tu *taskUseCase) DeleteTaskById(c context.Context, taskId string)error{
	ctx, close := context.WithTimeout(c, tu.contextTimeout)
	defer close()

	return tu.taskRepository.DeleteTask(ctx , taskId)
}