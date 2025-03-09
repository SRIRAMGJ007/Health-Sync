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
		userGroup.PUT("/profile", func(ctx *gin.Context) {
			user.UpdateUserProfile(ctx, queries)
		})
		userGroup.GET("/doctors", func(ctx *gin.Context) {
			user.ListDoctorsHandler(ctx, queries)
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
	}
}
