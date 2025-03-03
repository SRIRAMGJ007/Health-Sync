package routes

import (
	"github.com/SRIRAMGJ007/Health-Sync/internal/handler/auth"
	"github.com/SRIRAMGJ007/Health-Sync/internal/repository"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.Engine, queries *repository.Queries) {
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", func(ctx *gin.Context) {
			auth.RegisterHandler(ctx, queries)
		})
		authGroup.POST("/login", func(ctx *gin.Context) {
			auth.LoginHandler(ctx, queries)
		})
		authGroup.GET("/google/callback", func(ctx *gin.Context) {
			auth.GoogleAuthhandler(ctx, queries)
		})
	}
}
