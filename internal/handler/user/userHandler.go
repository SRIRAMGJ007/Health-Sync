package user

import (
	"context"
	"log"
	"net/http"

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

type DoctorResponse struct {
	ID             string                 `json:"id"`
	Name           string                 `json:"name"`
	Specialization string                 `json:"specialization"`
	Experience     int32                  `json:"experience"`
	Qualification  string                 `json:"qualification"`
	HospitalName   string                 `json:"hospital_name"`
	Availability   []AvailabilityResponse `json:"availability,omitempty"`
}

func UpdateUserProfile(ctx *gin.Context, queries *repository.Queries) {
	var req UpdateUserRequest

	userID := ctx.Query("id")
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
