package api

import (
	"dooqiniu/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UploadHandler handles file upload requests to Qiniu Cloud
func UploadHandler(c *gin.Context) {
	// Get parameters from request
	filePath := c.Query("filePath")
	objectName := c.Query("objectName")

	if filePath == "" || objectName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "filePath and objectName are required parameters",
		})
		return
	}

	// Initialize Qiniu uploader
	uploader := service.NewQiniuClient()

	// Perform the upload
	err := uploader.Upload(filePath, objectName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "upload failed: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "上传成功",
	})
}

// ListFilesHandler 处理获取七牛云文件列表的请求
func ListFilesHandler(c *gin.Context) {
	// 获取请求参数
	prefix := c.DefaultQuery("prefix", "") // 文件前缀
	marker := c.DefaultQuery("marker", "") // 游标，列举时继续读取上次的 marker
	limit := 1000                          // 默认每次最多列举 1000 个文件
	if c.Query("limit") != "" {
		// 如果有指定 limit，转换为整数
		parsedLimit, err := strconv.Atoi(c.DefaultQuery("limit", "1000"))
		if err == nil && parsedLimit > 0 && parsedLimit <= 1000 {
			limit = parsedLimit
		}
	}

	// 初始化七牛云客户端
	client := service.NewQiniuClient()

	// 获取文件列表
	files, nextMarker, err := client.ListFiles(prefix, marker, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "Error getting file list: " + err.Error(),
		})
		return
	}

	// 返回文件列表和下一页游标
	c.JSON(http.StatusOK, gin.H{
		"code":        http.StatusOK,
		"msg":         "File list retrieved successfully",
		"files":       files,
		"next_marker": nextMarker,
	})
}
