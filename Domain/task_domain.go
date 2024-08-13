package Domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Task struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	Title       string             `json:"title" bson:"title" validate:"required"`
	Description string             `json:"description" bson:"description" validate:"required"`
	DueDate     time.Time          `json:"due_date" bson:"due_date" validate:"required"`
	Status      string             `json:"status" bson:"status" validate:"required"`
}

type TaskRepository interface {
	FindAlltasks(ctx context.Context) ([]Task, error)
	FindTaskById(ctx context.Context, taskId string) (Task, error)
	CreateTask(ctx context.Context, task Task) (Task, error)
	UpdateTaskById(ctx context.Context, task Task, id string) (Task, error)
	DeleteTask(ctx context.Context, taskId string) error
}

type TaskUseCase interface {
	GetAllTasks(ctx context.Context) ([]Task, error)
	GetTaskByID(ctx context.Context, taskId string) (Task, error)
	AddNewTask(ctx context.Context, task Task) (Task, error)
	ModifyTaskById(ctx context.Context, task Task, taskId string) (Task, error)
	DeleteTaskById(ctx context.Context, taskId string) error
}
