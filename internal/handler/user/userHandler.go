package user

import (
	"context"
	"log"
	"net/http"

	"github.com/SRIRAMGJ007/Health-Sync/internal/repository"
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

type DoctorResponse struct {
	ID             pgtype.UUID `json:"id"`
	Name           string      `json:"name"`
	Specialization string      `json:"specialization"`
	Experience     int32       `json:"experience"`
	Qualification  string      `json:"qualification"`
	HospitalName   string      `json:"hospital_name"`
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
			ID:             doctor.ID,
			Name:           doctor.Name,
			Specialization: doctor.Specialization,
			Experience:     doctor.Experience,
			Qualification:  doctor.Qualification,
			HospitalName:   doctor.HospitalName,
		}
	}

	ctx.JSON(http.StatusOK, resp)

}
