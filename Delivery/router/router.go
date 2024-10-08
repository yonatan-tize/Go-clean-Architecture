package router

import (
	controllers "example/go-clean-architecture/Delivery/controllers"
	infrastructure "example/go-clean-architecture/Infrastructure"
	repository "example/go-clean-architecture/Repositories"
	usecases "example/go-clean-architecture/Usecases"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// SetUpRouter sets up the router for the application.
// It configures the routes and middleware for different endpoints.
// The router parameter is a pointer to a gin.Engine instance.
// The db parameter is a mongo.Database instance representing the database connection.
// The time parameter is a time.Duration value representing the duration for certain operations.
func SetUpRouter(router *gin.Engine, db mongo.Database, time time.Duration) {

	tr := repository.NewTaskRepository(db, "tasks")
	ur := repository.NewUserRepository(db, "users")

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
