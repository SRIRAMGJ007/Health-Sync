package user

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/SRIRAMGJ007/Health-Sync/internal/repository"
	"github.com/SRIRAMGJ007/Health-Sync/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type UpdateUserRequest struct {
	Name                     *string `json:"name"`
	Age                      *int32  `json:"age"`
	Gender                   *string `json:"gender"`
	BloodGroup               *string `json:"blood_group"`
	EmergencyContactNumber   *string `json:"emergency_contact_number"`
	EmergencyContactRelation *string `json:"emergency_contact_relationship"`
}

type AvailabilityResponse struct {
	ID               pgtype.UUID `json:"id"`
	DoctorID         pgtype.UUID `json:"doctor_id"`
	AvailabilityDate pgtype.Date `json:"availability_date"`
	StartTime        string      `json:"start_time"`
	EndTime          string      `json:"end_time"`
	IsBooked         bool        `json:"is_booked"`
}

type UserProfileResponse struct {
	ID                           string `json:"id"`
	Email                        string `json:"email"`
	Name                         string `json:"name,omitempty"`
	Age                          int32  `json:"age,omitempty"`
	Gender                       string `json:"gender,omitempty"`
	BloodGroup                   string `json:"blood_group,omitempty"`
	EmergencyContactNumber       string `json:"emergency_contact_number,omitempty"`
	EmergencyContactRelationship string `json:"emergency_contact_relationship,omitempty"`
}

type DoctorResponse struct {
	ID             string                 `json:"id"`
	Name           string                 `json:"name"`
	Specialization string                 `json:"specialization"`
	Experience     int32                  `json:"experience"`
	Qualification  string                 `json:"qualification"`
	HospitalName   string                 `json:"hospital_name"`
	Availability   []AvailabilityResponse `json:"availability,omitempty"`
}

type CreateMedicationRequest struct {
	MedicationName string `json:"medication_name" binding:"required"`
	Dosage         string `json:"dosage" binding:"required"`
	TimeToNotify   string `json:"time_to_notify" binding:"required"`
	Frequency      string `json:"frequency" binding:"required,oneof=daily weekly"` // Use oneof for validation
}

func safeDerefString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func safeDerefInt32(i *int32) int32 {
	if i != nil {
		return *i
	}
	return 0 // Or -1 for "unset" indication
}

func UpdateUserProfile(ctx *gin.Context, queries *repository.Queries) {
	var req UpdateUserRequest

	userID := ctx.Param("id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "userid missing"})
		log.Println("update user profile failed -> query parameter id missing, bad reuest")
		return
	}

	parsedID, err := uuid.Parse(userID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		log.Printf("update user profile failed -> invalid userid, bad reuest: { %v }", err)
		return
	}

	pgUUID := pgtype.UUID{Bytes: parsedID, Valid: true}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "invalid request body"})
		log.Printf("update user profile failed -> invalid request body, bad reuest : { %v }", err)
		return
	}

	err = queries.UpdateUserProfile(context.Background(), repository.UpdateUserProfileParams{
		ID:                           pgUUID,
		Name:                         req.Name,
		Age:                          req.Age,
		BloodGroup:                   req.BloodGroup,
		Gender:                       req.Gender,
		EmergencyContactNumber:       req.EmergencyContactNumber,
		EmergencyContactRelationship: req.EmergencyContactRelation,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user service"})
		log.Printf("update user profile failed : { %v }", err)
		return
	}

	updatedUser, err := queries.GetUserByID(context.Background(), pgUUID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch updated user "})
		log.Printf("failed to fetch updated user profile data : { %v }", err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "user profile updated sucessfully",
		"user": gin.H{
			"id":                             updatedUser.ID,
			"name":                           updatedUser.Name,
			"age":                            updatedUser.Age,
			"blood_group":                    updatedUser.BloodGroup,
			"gender":                         updatedUser.Gender,
			"emergency_contact_number":       updatedUser.EmergencyContactNumber,
			"emergency_contact_relationship": updatedUser.EmergencyContactRelationship,
		},
	})

}

func ListDoctorsHandler(ctx *gin.Context, queries *repository.Queries) {

	doctors, err := queries.ListDoctors(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve doctors"})
		return
	}

	resp := make([]DoctorResponse, len(doctors))
	for i, doctor := range doctors {
		resp[i] = DoctorResponse{
			ID:             doctor.ID.String(),
			Name:           doctor.Name,
			Specialization: doctor.Specialization,
			Experience:     doctor.Experience,
			Qualification:  doctor.Qualification,
			HospitalName:   doctor.HospitalName,
		}
	}

	ctx.JSON(http.StatusOK, resp)

}

func GetDoctorByIDHandler(ctx *gin.Context, queries *repository.Queries) {
	log.Println("GetDoctorByIDHandler: Request received")
	doctorIDStr := ctx.Param("doctorId")

	doctorID, err := uuid.Parse(doctorIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID"})
		return
	}

	parsedDoctorID := pgtype.UUID{Bytes: doctorID, Valid: true}

	doctor, err := queries.GetDoctorByID(ctx, parsedDoctorID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Doctor not found"})
		return
	}

	availability, err := queries.GetDoctorAvailabilityByDoctor(ctx, parsedDoctorID) // Assuming you have this function
	if err != nil {
		log.Printf("Error fetching availability: %v", err)
		// We'll continue even if fetching availability fails
	}

	var availabilityResponses []AvailabilityResponse
	for _, avail := range availability {
		availabilityResponses = append(availabilityResponses, AvailabilityResponse{
			ID:               avail.ID,
			DoctorID:         avail.DoctorID,
			AvailabilityDate: avail.AvailabilityDate,
			StartTime:        utils.FormatTime(avail.StartTime),
			EndTime:          utils.FormatTime(avail.EndTime),
			IsBooked:         *avail.IsBooked,
		})
	}

	resp := DoctorResponse{
		ID:             doctor.ID.String(),
		Name:           doctor.Name,
		Specialization: doctor.Specialization,
		Experience:     doctor.Experience,
		Qualification:  doctor.Qualification,
		HospitalName:   doctor.HospitalName,
		Availability:   availabilityResponses,
	}

	ctx.JSON(http.StatusOK, resp)
}

func GetUserProfile(ctx *gin.Context, queries *repository.Queries) {

	log.Println("GetUserProfile: Request received")
	userIDStr := ctx.Param("userid")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	parsedUserID := pgtype.UUID{Bytes: userID, Valid: true}

	user, err := queries.GetUserProfileByID(ctx, parsedUserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	resp := UserProfileResponse{
		ID:                           user.ID.String(),
		Email:                        user.Email,
		Name:                         safeDerefString(user.Name),
		Age:                          safeDerefInt32(user.Age),
		Gender:                       safeDerefString(user.Gender),
		BloodGroup:                   safeDerefString(user.BloodGroup),
		EmergencyContactNumber:       safeDerefString(user.EmergencyContactNumber),
		EmergencyContactRelationship: safeDerefString(user.EmergencyContactRelationship),
	}

	ctx.JSON(http.StatusOK, resp)

}

func CreateMedicationHandler(ctx *gin.Context, queries *repository.Queries) {
	userIDStr := ctx.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		log.Printf("CreateMedicationHandler: Invalid user ID: %v", err)
		return
	}

	var req CreateMedicationRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		log.Printf("CreateMedicationHandler: Invalid request body: %v", err)
		return
	}

	// Validate TimeToNotify format
	_, err = time.Parse("15:04:05", req.TimeToNotify)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time_to_notify format. Use HH:MM:SS"})
		log.Printf("CreateMedicationHandler: Invalid time format: %v", err)
		return
	}

	// Parse time in IST
	ist, err := time.LoadLocation("Asia/Kolkata") // Load IST location
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load IST location"})
		log.Printf("CreateMedicationHandler: Failed to load IST location: %v", err)
		return
	}

	timeToNotify, err := time.ParseInLocation("15:04:05", req.TimeToNotify, ist) // Parse with IST location
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time_to_notify format. Use HH:MM:SS in IST"})
		log.Printf("CreateMedicationHandler: Invalid time format: %v", err)
		return
	}

	// Convert time to pgtype.Time (UTC)
	utcTime := timeToNotify.UTC() // convert time to UTC for storing in database.
	pgxTime := pgtype.Time{
		Microseconds: int64((utcTime.Hour()*3600 + utcTime.Minute()*60 + utcTime.Second()) * 1e6),
		Valid:        true,
	}

	medication, err := queries.CreateMedication(ctx, repository.CreateMedicationParams{
		UserID:         pgtype.UUID{Bytes: userID, Valid: true},
		MedicationName: req.MedicationName,
		Dosage:         req.Dosage,
		TimeToNotify:   pgxTime,
		Frequency:      req.Frequency,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create medication"})
		log.Printf("CreateMedicationHandler: Failed to create medication: %v", err)
		return
	}

	readableTime := time.Unix(0, medication.TimeToNotify.Microseconds*1000).UTC().Format("15:04:05")
	response := gin.H{
		"medication": gin.H{
			"ID":             medication.ID,
			"UserID":         medication.UserID,
			"MedicationName": medication.MedicationName,
			"Dosage":         medication.Dosage,
			"TimeToNotify":   readableTime, // Use the readable time
			"Frequency":      medication.Frequency,
			"IsReadbyuser":   medication.IsReadbyuser,
			"CreatedAt":      medication.CreatedAt,
			"UpdatedAt":      medication.UpdatedAt,
		},
		"message": "Medication scheduled successfully",
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Medication scheduled successfully", "medication": response})
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

func GetMedicationsByUserIDHandler(ctx *gin.Context, queries *repository.Queries) {
	userIDStr := ctx.Param("user_id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		log.Printf("GetMedicationsByUserIDHandler: Invalid user ID: %v", err)
		return
	}

	medications, err := queries.GetMedicationsByUserID(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve medications"})
		log.Printf("GetMedicationsByUserIDHandler: Failed to retrieve medications: %v", err)
		return
	}

	// Format the response
	response := make([]gin.H, len(medications))
	for i, medication := range medications {
		readableTime := time.Unix(0, medication.TimeToNotify.Microseconds*1000).UTC().Format("15:04:05") //Format time
		response[i] = gin.H{
			"ID":             medication.ID,
			"MedicationName": medication.MedicationName,
			"Dosage":         medication.Dosage,
			"TimeToNotify":   readableTime, // Format time
			"Frequency":      medication.Frequency,
			"IsReadbyuser":   medication.IsReadbyuser,
			"CreatedAt":      medication.CreatedAt,
			"UpdatedAt":      medication.UpdatedAt,
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"medications": response})
}
