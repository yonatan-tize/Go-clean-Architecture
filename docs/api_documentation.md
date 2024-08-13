# Task Manager Project

## Overview
This is an implementation of task_manager using Go Clean Architecture
The Task Manager project is a straightforward task management system developed using Go. It enables users to create, read, update, and delete tasks. The project is built with the Gin framework for the web server and leverages the official MongoDB driver for database operations.

Upon registration, the first user is automatically assigned an Admin role, while subsequent users are given a standard User role. Admins have the ability to promote other users to the Admin role.

The system includes authentication and authorization features, ensuring that users must be logged in to perform any actions. Depending on their role, users are granted different levels of access. Admins can create, update, and delete tasks, as well as retrieve all tasks or view a specific task by its ID. Regular users, however, are restricted to viewing all tasks or retrieving a specific task by its ID.

## Setup Instructions

### Prerequisites

- Go (version 1.16 or higher)
- A running MongoDB instance

### Installation


1. **Set up the database:**

    Update the database connection string in `data/task_service.go` with your own MongoDB connection parameters.

    ```go
    // Example connection string
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    ```

2. **Install dependencies:**

    ```sh
    go mod tidy
    ```

3. **Run the application:**

    ```sh
    go run main.go
    ```

## API Documentation

You can refer to the detailed API documentation using the link below:

[Postman API Documentation](https://documenter.getpostman.com/view/37171778/2sA3s3JC8E)

### Endpoints

folder structure

task-manager/
├── Delivery/
│   ├── main.go
│   ├── controllers/
│   │   └── controller.go
│   └── routers/
│       └── router.go
├── Domain/
│   └── domain.go
|__ db
|   |__database.go
├── Infrastructure/
│   ├── auth_middleWare.go
│   ├── jwt_service.go
│   └── password_service.go
├── Repositories/
│   ├── task_repository.go
│   └── user_repository.go
├── Usecases/
│   ├── task_usecases.go
│   └── user_usecases.go
├── docs/
│   └── api_documentation.md
└── README.md
