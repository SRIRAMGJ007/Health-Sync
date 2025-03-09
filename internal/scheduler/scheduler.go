package scheduler

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/SRIRAMGJ007/Health-Sync/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// checkMedicationsToNotify checks for medications that need to be notified.
func checkMedicationsToNotify(ctx context.Context, queries *repository.Queries) {
	now := time.Now().UTC() // Get current UTC time.
	log.Printf("now time %v", now)
	currentTime := pgtype.Time{
		Microseconds: int64((now.Hour()*3600 + now.Minute()*60 + now.Second()) * 1e6),
		Valid:        true,
	}
	readableTime := time.Unix(0, currentTime.Microseconds*1000).UTC().Format("15:04:05")
	log.Printf("Checking medications for time: %v", readableTime)

	medications, err := queries.GetMedicationsToNotify(ctx, currentTime)
	if err != nil {
		log.Printf("Error retrieving medications: %v", err)
		return
	}

	for _, medication := range medications {
		go func(medication repository.Medication) { // Launch a Go routine for each medication
			// Retrieve user's FCM token
			fcmToken, err := queries.GetUserFCMToken(ctx, medication.UserID)
			if err != nil {
				log.Printf("Error retrieving FCM token: %v", err)
				return
			}

			// Send push notification using FCM
			err = sendPushNotification(*fcmToken, medication.MedicationName, medication.Dosage) // Replace with your FCM logic
			if err != nil {
				log.Printf("Error sending push notification: %v", err)
				return
			}

			// Update is_readbyuser flag
			err = queries.UpdateMedicationReadStatus(ctx, medication.ID)
			if err != nil {
				log.Printf("Error updating medication read status: %v", err)
				return
			}

			log.Printf("Notification sent for medication: %s", medication.MedicationName)
		}(medication)
	}
}

// StartMedicationScheduler initializes and starts the medication scheduler.
func StartMedicationScheduler(ctx context.Context, queries *repository.Queries) {
	log.Println("Starting medication scheduler...")

	ticker := time.NewTicker(1 * time.Minute) // Check every minute
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("Scheduler tick: Checking medications...")
			checkMedicationsToNotify(ctx, queries)
		case <-ctx.Done():
			log.Println("Medication scheduler stopped.")
			return
		}
	}
}

func MarkMedicationAsReadHandler(ctx *gin.Context, queries *repository.Queries) {
	medicationIDStr := ctx.Param("medication_id")

	medicationID, err := uuid.Parse(medicationIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid medication ID"})
		log.Printf("MarkMedicationAsReadHandler: Invalid medication ID: %v", err)
		return
	}

	err = queries.UpdateMedicationReadStatus(ctx, pgtype.UUID{Bytes: medicationID, Valid: true})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark medication as read"})
		log.Printf("MarkMedicationAsReadHandler: Failed to mark medication as read: %v", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Medication marked as read"})
}
