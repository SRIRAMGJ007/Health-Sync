package routes

import (
	"github.com/SRIRAMGJ007/Health-Sync/internal/handler/booking"
	"github.com/SRIRAMGJ007/Health-Sync/internal/handler/doctor"
	"github.com/SRIRAMGJ007/Health-Sync/internal/middleware"
	"github.com/SRIRAMGJ007/Health-Sync/internal/repository"
	"github.com/gin-gonic/gin"
)

func DoctorRoutes(r *gin.Engine, queries *repository.Queries) {
	doctorGroup := r.Group("/doctors/:doctorId/availability")
	doctorGroup.Use(middleware.ValidateJWT())
	{
		doctorGroup.POST("/", func(ctx *gin.Context) {
			doctor.CreateAvailabilityHandler(ctx, queries)
		})
		doctorGroup.GET("/", func(ctx *gin.Context) {
			doctor.GetAvailabilityByDoctorHandler(ctx, queries)
		})
		doctorGroup.GET("/date/:date", func(ctx *gin.Context) {
			doctor.GetAvailabilityByDoctorAndDateHandler(ctx, queries)
		})
		doctorGroup.PUT("/:availabilityId", func(ctx *gin.Context) {
			doctor.UpdateAvailabilityHandler(ctx, queries)
		})
		doctorGroup.GET("/bookings", func(ctx *gin.Context) {
			booking.GetBookingsByDoctorIDHandler(ctx, queries)
		})
		doctorGroup.PUT("/bookings/:bookingId/status", func(ctx *gin.Context) {
			booking.UpdateBookingStatusHandler(ctx, queries)
		})
	}
}
