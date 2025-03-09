package main

import (
	"fmt"
	"log"

	"github.com/SRIRAMGJ007/Health-Sync/internal/database"
	"github.com/SRIRAMGJ007/Health-Sync/internal/repository"
	"github.com/SRIRAMGJ007/Health-Sync/internal/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {

	err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	queries := repository.New(database.DB)

	r := gin.Default()

	routes.AuthRoutes(r, queries)
	routes.UserRoutes(r, queries)
	routes.DoctorRoutes(r, queries)

	fmt.Println("Health-Sync backend is running on port 8080...")
	err = r.Run(":8080")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
