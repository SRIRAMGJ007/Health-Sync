package booking

import (
	"net/http"

	"github.com/SRIRAMGJ007/Health-Sync/internal/repository"
	"github.com/SRIRAMGJ007/Health-Sync/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CreateBookingRequest struct {
	AvailabilityID pgtype.UUID `json:"availability_id" binding:"required"`
}

type UpdateBookingStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type DoctorResponse struct {
	ID             pgtype.UUID `json:"id"`
	Name           string      `json:"name"`
	Specialization string      `json:"specialization"`
	Experience     int32       `json:"experience"`
	Qualification  string      `json:"qualification"`
	HospitalName   string      `json:"hospital_name"`
}

type BookingResponse struct {
	ID               pgtype.UUID    `json:"id"`
	UserID           pgtype.UUID    `json:"user_id"`
	DoctorID         pgtype.UUID    `json:"doctor_id"`
	AvailabilityID   pgtype.UUID    `json:"availability_id"`
	BookingDate      string         `json:"booking_date"`
	BookingStartTime string         `json:"booking_start_time"`
	BookingEndTime   string         `json:"booking_end_time"`
	Status           string         `json:"status"`
	Doctor           DoctorResponse `json:"doctor"`
}

type DoctorBookingResponse struct {
	BookingDate      pgtype.Date `json:"booking_date"`
	BookingStartTime string      `json:"booking_start_time"`
	BookingEndTime   string      `json:"booking_end_time"`
	Status           string      `json:"status"`
	PatientName      string      `json:"patient_name"`
	AvailabilityID   pgtype.UUID `json:"availability_id"`
}

func CreateBookingHandler(ctx *gin.Context, queries *repository.Queries) {

	userIDStr := ctx.Param("userId")
	availabilityIDStr := ctx.Param("availabilityId")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	parsedUserID := pgtype.UUID{Bytes: userID, Valid: true}

	availabilityID, err := uuid.Parse(availabilityIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid availability ID"})
		return
	}
	parsedAvailabilityID := pgtype.UUID{Bytes: availabilityID, Valid: true}

	availability, err := queries.GetDoctorAvailabilityByID(ctx, parsedAvailabilityID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "AvailabilityiD not found"})
		return
	}

	if *availability.IsBooked {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Availability slot is already booked"})
		return
	}

	booking, err := queries.CreateBooking(ctx, repository.CreateBookingParams{
		UserID:           parsedUserID,
		DoctorID:         availability.DoctorID,
		AvailabilityID:   parsedAvailabilityID,
		BookingDate:      availability.AvailabilityDate,
		BookingStartTime: availability.StartTime,
		BookingEndTime:   availability.EndTime,
		Status:           "pending",
	})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	bookingStatus := true
	err = queries.UpdateAvailabilityBookedStatus(ctx, repository.UpdateAvailabilityBookedStatusParams{
		ID:       parsedAvailabilityID,
		IsBooked: &bookingStatus,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update availability status"})
		return
	}

	resp := BookingResponse{
		ID:               booking.ID,
		UserID:           booking.UserID,
		DoctorID:         booking.DoctorID,
		AvailabilityID:   booking.AvailabilityID,
		BookingDate:      booking.BookingDate.Time.String(),
		BookingStartTime: utils.FormatTime(booking.BookingStartTime),
		BookingEndTime:   utils.FormatTime(booking.BookingEndTime),
		Status:           booking.Status,
	}

	ctx.JSON(http.StatusCreated, resp)

}
func GetBookingByIDHandler(ctx *gin.Context, queries *repository.Queries) {

	bookingIDStr := ctx.Param("bookingId")

	bookingID, err := uuid.Parse(bookingIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	parsedBookingID := pgtype.UUID{Bytes: bookingID, Valid: true}

	booking, err := queries.GetBookingByID(ctx, parsedBookingID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Booking not found"})
		return
	}

	resp := BookingResponse{
		ID:               booking.ID,
		UserID:           booking.UserID,
		DoctorID:         booking.DoctorID,
		AvailabilityID:   booking.AvailabilityID,
		BookingDate:      booking.BookingDate.Time.String(),
		BookingStartTime: utils.FormatTime(booking.BookingStartTime),
		BookingEndTime:   utils.FormatTime(booking.BookingEndTime),
		Status:           booking.Status,
	}

	ctx.JSON(http.StatusOK, resp)

}
func GetBookingsByUserIDHandler(ctx *gin.Context, queries *repository.Queries) {

	userIDStr := ctx.Param("userId")

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	parsedUserID := pgtype.UUID{Bytes: userID, Valid: true}

	bookings, err := queries.GetBookingsByUserID(ctx, parsedUserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve bookings"})
		return
	}

	resp := make([]BookingResponse, len(bookings))
	for i, booking := range bookings {
		doctor, err := queries.GetDoctorByID(ctx, booking.DoctorID) // Assuming you have GetDoctorByID
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve doctor details"})
			return
		}

		resp[i] = BookingResponse{
			ID:               booking.ID,
			UserID:           booking.UserID,
			DoctorID:         booking.DoctorID,
			AvailabilityID:   booking.AvailabilityID,
			BookingDate:      booking.BookingDate.Time.String(),
			BookingStartTime: utils.FormatTime(booking.BookingStartTime),
			BookingEndTime:   utils.FormatTime(booking.BookingEndTime),
			Status:           booking.Status,
			Doctor: DoctorResponse{
				ID:             doctor.ID,
				Name:           doctor.Name,
				Specialization: doctor.Specialization,
				HospitalName:   doctor.HospitalName,
				Experience:     doctor.Experience,
				Qualification:  doctor.Qualification,
			},
		}

	}
	ctx.JSON(http.StatusOK, resp)

}
func UpdateBookingStatusHandler(ctx *gin.Context, queries *repository.Queries) {

	bookingIDStr := ctx.Param("bookingId")

	bookingID, err := uuid.Parse(bookingIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	parsedBookingID := pgtype.UUID{Bytes: bookingID, Valid: true}
	var req UpdateBookingStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = queries.UpdateBookingStatus(ctx, repository.UpdateBookingStatusParams{
		ID:     parsedBookingID,
		Status: req.Status,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update booking status"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Booking status updated successfully"})

}
func GetBookingsByDoctorIDHandler(ctx *gin.Context, queries *repository.Queries) {

	doctorIDStr := ctx.Param("doctorId")

	doctorID, err := uuid.Parse(doctorIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid doctor ID"})
		return
	}

	parsedDoctorID := pgtype.UUID{Bytes: doctorID, Valid: true}

	bookings, err := queries.GetBookingsByDoctorID(ctx, parsedDoctorID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve bookings"})
		return
	}

	resp := make([]DoctorBookingResponse, len(bookings))
	for i, booking := range bookings {
		user, err := queries.GetUserByID(ctx, booking.UserID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve user details"})
			return
		}

		resp[i] = DoctorBookingResponse{
			BookingDate:      booking.BookingDate,
			BookingStartTime: utils.FormatTime(booking.BookingStartTime),
			BookingEndTime:   utils.FormatTime(booking.BookingEndTime),
			Status:           booking.Status,
			PatientName:      *user.Name,
			AvailabilityID:   booking.AvailabilityID,
		}
	}

	ctx.JSON(http.StatusOK, resp)

}
func DeleteBookingHandler(ctx *gin.Context, queries *repository.Queries) {

	bookingIDStr := ctx.Param("bookingId")

	bookingID, err := uuid.Parse(bookingIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}

	parsedBookingID := pgtype.UUID{Bytes: bookingID, Valid: true}

	err = queries.DeleteBooking(ctx, parsedBookingID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete booking"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Booking deleted successfully"})

}
