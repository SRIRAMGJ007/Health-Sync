package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/SRIRAMGJ007/Health-Sync/internal/database"
	"github.com/SRIRAMGJ007/Health-Sync/internal/repository"
	"github.com/SRIRAMGJ007/Health-Sync/internal/routes"
	"github.com/SRIRAMGJ007/Health-Sync/internal/scheduler"
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
	// Initialize Database
	err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize SQLc Queries
	queries := repository.New(database.DB)

	// Initialize Firebase
	err = scheduler.InitializeFirebase()
	if err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}

	// Initialize Gin Router
	r := gin.Default()

	// Setup Routes
	routes.AuthRoutes(r, queries)
	routes.UserRoutes(r, queries)
	routes.DoctorRoutes(r, queries)
	routes.EMRRoutes(r, queries)

	// Start Scheduler (as a Go routine)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go scheduler.StartMedicationScheduler(ctx, queries)

	// Start Server
	httpServer := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: r,
	}

	httpsServer := &http.Server{
		Addr:    "0.0.0.0:8443",
		Handler: r,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start Http server: %v", err)
		}
	}()
	go func() {
		if err := httpsServer.ListenAndServeTLS("/Health-Sync/cert.pem", "/Health-Sync/key.pem"); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start Https server: %v", err)
		}
	}()

	fmt.Println("Health-Sync backend is running on port HTTP (8080) & HTTPS (8443)......")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a timeout context for shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	// Shutdown both servers gracefully
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP server forced to shutdown: %v", err)
	}

	if err := httpsServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTPS server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")

	log.Println("Server exiting")
}
