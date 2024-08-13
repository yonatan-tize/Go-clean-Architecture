package main

import (
	"example/go-clean-architecture/Delivery/router"
	"example/go-clean-architecture/db"
	"time"

	"github.com/gin-gonic/gin"
)

// main is the entry point of the application.
// It sets up the Gin router, connects to the MongoDB database,
// sets up the router with the database connection, and starts the server.
// The server listens on localhost:8080.
func main() {
	r := gin.Default()
	databse := db.ConnectDB("mongodb://localhost:27017")
	defer db.DisconnectDB()

	router.SetUpRouter(r, *databse, 100 * time.Second)

	r.Run("localhost:8080")
}
