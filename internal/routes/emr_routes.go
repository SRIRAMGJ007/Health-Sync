package routes

import (
	"github.com/SRIRAMGJ007/Health-Sync/internal/handler/records"
	"github.com/SRIRAMGJ007/Health-Sync/internal/middleware"
	"github.com/SRIRAMGJ007/Health-Sync/internal/repository"
	"github.com/gin-gonic/gin"
)

func EMRRoutes(r *gin.Engine, queries *repository.Queries) {
	recordGroup := r.Group("/EMR")
	recordGroup.Use(middleware.ValidateJWT())
	{
		recordGroup.POST("/:userid/emr/upload", func(ctx *gin.Context) {
			records.UploadMedicalRecord(ctx, queries)
		})
		recordGroup.GET("/:userid/emr/list", func(ctx *gin.Context) {
			records.ListMedicalRecords(ctx, queries)
		})
		recordGroup.GET("/:userid/record/:fileid/emr/download", func(ctx *gin.Context) {
			records.DownloadMedicalRecord(ctx, queries)
		})
		recordGroup.GET("/:userid/emr/view", func(ctx *gin.Context) {
			records.ViewMedicalRecord(ctx, queries)
		})
	}
}
