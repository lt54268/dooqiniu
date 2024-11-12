package api

import (
	"dooqiniu/internal/service"
	"net/http"
	"strconv"
	"time"

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

// DownloadHandler 处理文件下载请求
func DownloadFileHandler(c *gin.Context) {
	objectName := c.Query("objectName")
	accessType := c.DefaultQuery("accessType", "private") // "public" or "private"

	if objectName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "objectName is a required parameter",
		})
		return
	}

	// 初始化 Qiniu 客户端
	client := service.NewQiniuClient()

	var downloadURL string
	if accessType == "private" {
		// 设置链接的有效期为2小时
		expiryTime := time.Now().Add(2 * time.Hour).Unix()
		downloadURL = client.GeneratePrivateURL(objectName, expiryTime)
	} else {
		downloadURL = client.GeneratePublicURL(objectName)
	}

	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "生成下载链接成功",
		"data": gin.H{
			"downloadURL": downloadURL,
		},
	})
}

// DeleteFileHandler handles file delete requests from Qiniu Cloud
func DeleteFileHandler(c *gin.Context) {
	objectName := c.Query("objectName")
	if objectName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "objectName is a required parameter",
		})
		return
	}

	// 初始化七牛云客户端
	client := service.NewQiniuClient()

	// 调用 Delete 方法删除文件
	err := client.Delete(objectName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "failed to delete file: " + err.Error(),
		})
		return
	}

	// 删除成功，返回响应
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "文件删除成功",
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

// CopyFileHandler handles file copy requests on Qiniu Cloud
func CopyFileHandler(c *gin.Context) {
	// 获取源文件名和目标文件名
	srcKey := c.Query("srcObject")
	destKey := c.Query("destObject")
	forceStr := c.DefaultQuery("force", "false")

	if srcKey == "" || destKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "srcKey and destKey are required parameters",
		})
		return
	}

	// 将 force 参数转换为布尔值
	force, err := strconv.ParseBool(forceStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "invalid value for force parameter",
		})
		return
	}

	// 初始化七牛云客户端
	client := service.NewQiniuClient()

	// 执行复制操作
	err = client.Copy(srcKey, destKey, force)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "failed to copy file: " + err.Error(),
		})
		return
	}

	// 返回复制成功的响应
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "文件复制成功",
	})
}

// MoveFileHandler 处理文件移动请求
func MoveFileHandler(c *gin.Context) {
	// 获取源文件名和目标文件名
	srcKey := c.Query("srcObject")
	destKey := c.Query("destObject")
	forceStr := c.DefaultQuery("force", "false")

	if srcKey == "" || destKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "srcKey and destKey are required parameters",
		})
		return
	}

	// 将 force 参数转换为布尔值
	force, err := strconv.ParseBool(forceStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "invalid value for force parameter",
		})
		return
	}

	// 初始化七牛云客户端
	client := service.NewQiniuClient()

	// 执行移动操作
	err = client.Move(srcKey, destKey, force)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "failed to move file: " + err.Error(),
		})
		return
	}

	// 返回移动成功的响应
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "文件移动成功",
	})
}
