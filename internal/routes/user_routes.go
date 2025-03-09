package routes

import (
	"github.com/SRIRAMGJ007/Health-Sync/internal/handler/booking"
	"github.com/SRIRAMGJ007/Health-Sync/internal/handler/doctor"
	"github.com/SRIRAMGJ007/Health-Sync/internal/handler/user"
	"github.com/SRIRAMGJ007/Health-Sync/internal/middleware"
	"github.com/SRIRAMGJ007/Health-Sync/internal/repository"
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine, queries *repository.Queries) {
	userGroup := r.Group("/user")
	userGroup.Use(middleware.ValidateJWT())
	{

		userGroup.GET("/:userid/profile", func(ctx *gin.Context) {
			user.GetUserProfile(ctx, queries)
		})
		userGroup.PUT("/updateprofile/:id", func(ctx *gin.Context) {
			user.UpdateUserProfile(ctx, queries)
		})
		userGroup.GET("/doctors", func(ctx *gin.Context) {
			user.ListDoctorsHandler(ctx, queries)
		})
		userGroup.GET("/doctors/:doctorId", func(ctx *gin.Context) {
			user.GetDoctorByIDHandler(ctx, queries)
		})
		userGroup.GET("/doctors/:doctorId/availability", func(ctx *gin.Context) {
			doctor.GetAvailabilityByDoctorHandler(ctx, queries)
		})
		userGroup.POST("/bookings/users/:userId/availability/:availabilityId", func(ctx *gin.Context) {
			booking.CreateBookingHandler(ctx, queries)
		})
		userGroup.GET("/bookings/:bookingId", func(ctx *gin.Context) {
			booking.GetBookingByIDHandler(ctx, queries)
		})
		userGroup.GET("/bookings/users/:userId", func(ctx *gin.Context) {
			booking.GetBookingsByUserIDHandler(ctx, queries)
		})
		userGroup.DELETE("/bookings/:bookingId", func(ctx *gin.Context) {
			booking.DeleteBookingHandler(ctx, queries)
		})
		userGroup.POST("/:user_id/medications", func(ctx *gin.Context) { //Added medication schedule
			user.CreateMedicationHandler(ctx, queries)
		})
		userGroup.PUT("/:user_id/medications/:medication_id/read", func(ctx *gin.Context) { //Added medication mark read
			user.MarkMedicationAsReadHandler(ctx, queries)
		})
	}
}
