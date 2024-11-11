package router

import (
	"dooqiniu/internal/api"

	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置 Gin 路由
func SetupRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		v1.POST("/upload", api.UploadHandler)
		// v1.GET("/download", api.DownloadFileHandler)
		// v1.DELETE("/delete", api.DeleteFileHandler)
		v1.GET("/list", api.ListFilesHandler)
		// v1.POST("/copy", api.CopyFileHandler)
		// v1.POST("/move", api.MoveFileHandler)
	}
}
