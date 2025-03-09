package routes

import (
	"github.com/SRIRAMGJ007/Health-Sync/internal/handler/auth"
	"github.com/SRIRAMGJ007/Health-Sync/internal/repository"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine, queries *repository.Queries) {
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register/user", func(ctx *gin.Context) {
			auth.UserRegisterHandler(ctx, queries)
		})

		authGroup.POST("/register/doctor", func(ctx *gin.Context) {
			auth.DoctorRegisterHandler(ctx, queries)
		})
		authGroup.POST("/login/user", func(ctx *gin.Context) {
			auth.UserLoginHandler(ctx, queries)
		})
		authGroup.POST("/login/doctor", func(ctx *gin.Context) {
			auth.DoctorLoginHandler(ctx, queries)
		})
		authGroup.GET("/google/callback", func(ctx *gin.Context) {
			auth.GoogleAuthhandler(ctx, queries)
		})
	}
}
