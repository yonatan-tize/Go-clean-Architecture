package main

import (
	"example/go-clean-architecture/Delivery/router"
	"example/go-clean-architecture/db"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	databse := db.ConnectDB("mongodb://localhost:27017")
	defer db.DisconnectDB()

	router.SetUpRouter(r, *databse, 100 * time.Second)

	r.Run("localhost:8080")
}
