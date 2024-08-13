package usecases

import (
	"context"
	"time"
	domain "example/go-clean-architecture/Domain"
)


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

func (tu *taskUseCase) GetAllTasks(c context.Context) ([]domain.Task, error){
	ctx, close := context.WithTimeout(c, tu.contextTimeout)
	defer close()

	return tu.taskRepository.FindAlltasks(ctx)
}

func (tu *taskUseCase) GetTaskByID(c context.Context, taskId string) (domain.Task, error){
	ctx, close := context.WithTimeout(c, tu.contextTimeout)
	defer close()

	return tu.taskRepository.FindTaskById(ctx , taskId )
}
func (tu *taskUseCase) AddNewTask(c context.Context, task domain.Task) (domain.Task, error){
//check the validity of the incoming request(task)
    ctx, close := context.WithTimeout(c, tu.contextTimeout)
	defer close()

	return tu.taskRepository.CreateTask(ctx, task)
}
func (tu *taskUseCase) ModifyTaskById(c context.Context, task domain.Task,  taskId string) (domain.Task, error){
	ctx, close := context.WithTimeout(c, tu.contextTimeout)
	defer close()

	return tu.taskRepository.UpdateTaskById(ctx, task, taskId)
}
func (tu *taskUseCase) DeleteTaskById(c context.Context, taskId string)error{
	ctx, close := context.WithTimeout(c, tu.contextTimeout)
	defer close()

	return tu.taskRepository.DeleteTask(ctx , taskId)
}