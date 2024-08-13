package router

import (
	"time"
	controllers "example/go-clean-architecture/Delivery/controllers"
	infrastructure "example/go-clean-architecture/Infrastructure"
	usecases "example/go-clean-architecture/Usecases"
	repository "example/go-clean-architecture/repositories"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)
func SetUpRouter(router *gin.Engine, db mongo.Database, time time.Duration) {

	tr := repository.NewTaskRepository(db, "User")
	ur := repository.NewUserRepository(db, "Tasks")

	tc := controllers.TaskController{
		TaskUseCase: usecases.NewTaskUsecase(tr, time),
	}
	uc := controllers.UserController{
		UserUseCase: usecases.NewUserUsecase(ur, time),
	}
	// Public routes
	public := router.Group("/")
	{
		public.POST("/register", uc.CreateAccount)
		public.POST("/login", uc.Login)
	}

	// Authenticated routes
	authorized := router.Group("/")
	authorized.Use(infrastructure.AuthMiddleware())
	{
		authorized.GET("/tasks", tc.GetTasks)
		authorized.GET("/tasks/:id", tc.GetTask)
	}

	// Admin routes (require admin privileges)
	admin := router.Group("/admin")
	admin.Use(infrastructure.AuthMiddleware(), infrastructure.AuthAdminMiddleware())
	{
		admin.PUT("/promote/:id", uc.PromoteUser)
		admin.POST("/tasks", tc.CreateTask)
		admin.PUT("/tasks/:id", tc.UpdateTask)
		admin.DELETE("/tasks/:id", tc.DeleteTask)
	}
}
