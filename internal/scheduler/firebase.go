package scheduler

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	// Replace with your actual import path
	//... other imports ...
)

var firebaseApp *firebase.App
var messagingClient *messaging.Client

var firebaseProjectID string // Global variable within the scheduler package

func init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file in scheduler")
	}

	// Get Firebase project ID from environment variables
	firebaseProjectID = os.Getenv("FIREBASE_PROJECT_ID")
	if firebaseProjectID == "" {
		log.Fatal("FIREBASE_PROJECT_ID not found in .env file in scheduler")
	}
}

// InitializeFirebase initializes the Firebase app and messaging client.
func InitializeFirebase() error {
	opt := option.WithCredentialsFile("/Health-Sync/internal/scheduler/health-sync-30494-be8768d3833e.json")

	projectID := os.Getenv("FIREBASE_PROJECT_ID")

	if projectID == "" {
		log.Println("Firebase project ID not found in environment.")
		return os.ErrInvalid
	}

	conf := &firebase.Config{ProjectID: projectID}

	app, err := firebase.NewApp(context.Background(), conf, opt)

	if err != nil {
		return err
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		return err
	}

	firebaseApp = app
	messagingClient = client
	return nil
}

// sendPushNotification sends a push notification using FCM.
func sendPushNotification(_, medicationName, dosage string) error {
	if firebaseApp == nil || messagingClient == nil {
		log.Println("Firebase not initialized")
		return os.ErrInvalid
	}

	message := &messaging.Message{
		Token: "fwrAaiXwTBqN2XWdJuAm2A:APA91bExbnr3-F9yBQn4M-uXmW348KrL8KAgFq4uRlL8dUI0uZHmbOZFaZLB0_wMCiT-qpk_jW1D6fJ5R_YMwp0BE_NtD7O6bwZRXS1PdywfrGHGpoqFlDA",
		Notification: &messaging.Notification{
			Title: "Medication Reminder",
			Body:  medicationName + " - Dosage: " + dosage,
		},
	}

	_, err := messagingClient.Send(context.Background(), message)
	if err != nil {
		return err
	}

	return nil
}
