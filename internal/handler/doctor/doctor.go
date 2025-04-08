package doctor

import (
	"log"
	"net/http"
	"time"

	"github.com/SRIRAMGJ007/Health-Sync/internal/repository"
	"github.com/SRIRAMGJ007/Health-Sync/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/jackc/pgx/v5/pgtype"
)

type CreateAvailabilityRequest struct {
	AvailabilityDate string `json:"availability_date" binding:"required"`
	StartTime        string `json:"start_time" binding:"required"`
	EndTime          string `json:"end_time" binding:"required"`
}

type UpdateDoctorRequest struct {
	Name           string `json:"name"`
	Specialization string `json:"specialization"`
	Experience     string `json:"experience"`
	Qualification  string `json:"qualification"`
	HospitalName   string `json:"hospital_name"`
}

type AvailabilityResponse struct {
	ID               pgtype.UUID `json:"id"`
	DoctorID         pgtype.UUID `json:"doctor_id"`
	AvailabilityDate pgtype.Date `json:"availability_date"`
	StartTime        string      `json:"start_time"`
	EndTime          string      `json:"end_time"`
	IsBooked         bool        `json:"is_booked"`
}

type UpdateAvailabilityRequest struct {
	StartTime string `json:"start_time,omitempty" time_format:"15:04:05"`
	EndTime   string `json:"end_time,omitempty" time_format:"15:04:05"`
	IsBooked  bool   `json:"is_booked,omitempty"`
}

// func UpdateDoctorHandler(ctx *gin.Context, queries *repository.Queries) {

// 	doctorIDStr := ctx.Param("doctorId")

// 	doctorID, err := uuid.Parse(doctorIDStr)
// 	if err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID"})
// 		return
// 	}

// 	var req UpdateDoctorRequest
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	err = queries.UpdateDoctor(ctx, repository.UpdateDoctorParams{
// 		ID:             doctorID,
// 		Name:           req.Name,
// 		Specialization: req.Specialization,
// 		Experience:     req.Experience,
// 		Qualification:  req.Qualification,
// 		HospitalName:   req.HospitalName,
// 	})

// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update doctor"})
// 		return
// 	}

// 	ctx.JSON(http.StatusOK, gin.H{"message": "Doctor updated successfully"})

// }

func CreateAvailabilityHandler(ctx *gin.Context, queries *repository.Queries) {
	log.Println("CreateAvailabilityHandler: Request received")
	doctorID := ctx.Param("doctorId")
	uuid, err := uuid.Parse(doctorID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID"})
		return
	}
	parsedid := pgtype.UUID{Bytes: uuid, Valid: true}

	var req CreateAvailabilityRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse AvailabilityDate
	availabilityDate, err := time.Parse("2006-01-02", req.AvailabilityDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid availability_date format"})
		return
	}
	availabilityDatePg := pgtype.Date{Time: availabilityDate, Valid: true}

	// Parse StartTime
	startTime, err := time.Parse("15:04:05", req.StartTime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_time format"})
		return
	}

	// Parse EndTime
	endTime, err := time.Parse("15:04:05", req.EndTime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_time format"})
		return
	}

	// Calculate microseconds since midnight
	startTimeMicro := (startTime.Hour()*3600 + startTime.Minute()*60 + startTime.Second()) * 1e6
	endTimeMicro := (endTime.Hour()*3600 + endTime.Minute()*60 + endTime.Second()) * 1e6

	startTimePg := pgtype.Time{Microseconds: int64(startTimeMicro), Valid: true}
	endTimePg := pgtype.Time{Microseconds: int64(endTimeMicro), Valid: true}

	availability, err := queries.CreateDoctorAvailability(ctx, repository.CreateDoctorAvailabilityParams{
		DoctorID:         parsedid,
		AvailabilityDate: availabilityDatePg,
		StartTime:        startTimePg,
		EndTime:          endTimePg,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := AvailabilityResponse{
		ID:               availability.ID,
		DoctorID:         availability.DoctorID,
		AvailabilityDate: availability.AvailabilityDate,
		StartTime:        utils.FormatTime(availability.StartTime), // Formatted time
		EndTime:          utils.FormatTime(availability.EndTime),   // Formatted time
		IsBooked:         *availability.IsBooked,
	}

	ctx.JSON(http.StatusCreated, resp)

}

func UpdateAvailabilityHandler(ctx *gin.Context, queries *repository.Queries) {
	log.Println("UpdateAvailabilityHandler: Request received")
	doctorID := ctx.Param("doctorId")
	availabilityID := ctx.Param("availabilityId")

	doctorUUID, err := uuid.Parse(doctorID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID"})
		return
	}
	parsedDoctorID := pgtype.UUID{Bytes: doctorUUID, Valid: true}

	availabilityUUID, err := uuid.Parse(availabilityID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid availability ID"})
		return
	}
	parsedAvailabilityID := pgtype.UUID{Bytes: availabilityUUID, Valid: true}

	var req UpdateAvailabilityRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startTime, err := time.Parse("15:04:05", req.StartTime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_time format"})
		return
	}

	// Parse EndTime
	endTime, err := time.Parse("15:04:05", req.EndTime)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_time format"})
		return
	}

	// Calculate microseconds since midnight
	startTimeMicro := (startTime.Hour()*3600 + startTime.Minute()*60 + startTime.Second()) * 1e6
	endTimeMicro := (endTime.Hour()*3600 + endTime.Minute()*60 + endTime.Second()) * 1e6

	startTimePg := pgtype.Time{Microseconds: int64(startTimeMicro), Valid: true}
	endTimePg := pgtype.Time{Microseconds: int64(endTimeMicro), Valid: true}

	err = queries.UpdateDoctorAvailability(ctx, repository.UpdateDoctorAvailabilityParams{
		StartTime: startTimePg,
		EndTime:   endTimePg,
		IsBooked:  &req.IsBooked,
		ID:        parsedAvailabilityID,
		DoctorID:  parsedDoctorID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Availability updated successfully",
		"updates": gin.H{
			"start_time": req.StartTime,
			"end_time":   req.EndTime,
		},
	})

}

func GetAvailabilityByDoctorAndDateHandler(ctx *gin.Context, queries *repository.Queries) {

	doctorID := ctx.Param("doctorId")
	dateStr := ctx.Param("date")

	doctorUUID, err := uuid.Parse(doctorID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID"})
		return
	}
	parsedDoctorID := pgtype.UUID{Bytes: doctorUUID, Valid: true}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format (YYYY-MM-DD)"})
		return
	}

	datePg := pgtype.Date{Time: date, Valid: true}

	availabilitySlots, err := queries.GetDoctorAvailabilityByDoctorAndDate(ctx, repository.GetDoctorAvailabilityByDoctorAndDateParams{
		DoctorID:         parsedDoctorID,
		AvailabilityDate: datePg,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp := make([]AvailabilityResponse, len(availabilitySlots))
	for i, slot := range availabilitySlots {
		resp[i] = AvailabilityResponse{
			ID:               slot.ID,
			DoctorID:         slot.DoctorID,
			AvailabilityDate: slot.AvailabilityDate,
			StartTime:        utils.FormatTime(slot.StartTime),
			EndTime:          utils.FormatTime(slot.EndTime),
			IsBooked:         *slot.IsBooked,
		}
	}

	ctx.JSON(http.StatusOK, resp)

}

func GetAvailabilityByDoctorHandler(ctx *gin.Context, queries *repository.Queries) {

	doctorID := ctx.Param("doctorId")

	doctorUUID, err := uuid.Parse(doctorID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID"})
		log.Println("Error parsing doctor ID:", err)
		return
	}
	parsedDoctorID := pgtype.UUID{Bytes: doctorUUID, Valid: true}

	availabilitySlots, err := queries.GetDoctorAvailabilityByDoctor(ctx, parsedDoctorID)
	log.Printf("response: %v ----> error is %v", availabilitySlots, err)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println("Database error:", err)
		return
	}

	resp := make([]AvailabilityResponse, len(availabilitySlots))
	log.Printf("response: %v", resp)
	for i, slot := range availabilitySlots {
		resp[i] = AvailabilityResponse{
			ID:               slot.ID,
			DoctorID:         slot.DoctorID,
			AvailabilityDate: slot.AvailabilityDate,
			StartTime:        utils.FormatTime(slot.StartTime),
			EndTime:          utils.FormatTime(slot.EndTime),
			IsBooked:         *slot.IsBooked,
		}
	}
	log.Printf("response: %v", resp)

	ctx.JSON(http.StatusOK, resp)

}

func DeleteAvailabilityHandler(ctx *gin.Context, queries *repository.Queries) {

	doctorID := ctx.Param("doctorId")
	availabilityID := ctx.Param("availabilityId")

	doctorUUID, err := uuid.Parse(doctorID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID"})
		return
	}
	parsedDoctorID := pgtype.UUID{Bytes: doctorUUID, Valid: true}

	availabilityUUID, err := uuid.Parse(availabilityID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid availability ID"})
		return
	}
	parsedAvailabilityID := pgtype.UUID{Bytes: availabilityUUID, Valid: true}

	bookings, err := queries.GetBookingsByAvailabilityID(ctx, parsedAvailabilityID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check bookings"})
		return
	}

	if len(bookings) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete availability slot with bookings"})
		return
	}

	err = queries.DeleteDoctorAvailability(ctx, repository.DeleteDoctorAvailabilityParams{
		ID:       parsedAvailabilityID,
		DoctorID: parsedDoctorID,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Availability deleted successfully"})

}

// func DeleteDoctorHandler(queries *repository.Queries) gin.HandlerFunc {
// 	return func(ctx *gin.Context) {
// 		doctorIDStr := ctx.Param("doctorId")

// 		doctorID, err := uuid.Parse(doctorIDStr)
// 		if err != nil {
// 			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID"})
// 			return
// 		}

// 		err = queries.DeleteDoctor(ctx, doctorID)
// 		if err != nil {
// 			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete doctor"})
// 			return
// 		}

// 		ctx.JSON(http.StatusOK, gin.H{"message": "Doctor deleted successfully"})
// 	}
// }
