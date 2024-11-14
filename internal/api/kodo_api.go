package api

import (
	"dooqiniu/internal/model"
	"dooqiniu/internal/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// UploadHandler 文件上传接口
// @Summary 上传文件至七牛云
// @Description 根据文件路径和目标对象名称，将文件上传至七牛云存储
// @Tags 文件管理
// @Accept json
// @Produce json
// @Param filePath query string true "本地文件路径"
// @Param objectName query string true "目标对象名称"
// @Success 200 {object} map[string]interface{} "上传成功，返回文件信息"
// @Failure 400 {object} map[string]interface{} "缺少必要参数 filePath 和 objectName"
// @Failure 500 {object} map[string]interface{} "上传失败"
// @Router /api/v1/upload [get]
func UploadHandler(c *gin.Context) {
	// Get parameters from request
	filePath := c.Query("filePath")
	objectName := c.Query("objectName")

	if filePath == "" || objectName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "缺少filePath 和 objectName参数",
		})
		return
	}

	// Initialize Qiniu uploader
	uploader := service.NewQiniuClient()

	// Perform the upload and get file info
	uploadResponse, err := uploader.Upload(filePath, objectName)
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
		"data": uploadResponse,
	})
}

// DownloadFileHandler 生成文件下载链接接口
// @Summary 生成文件下载链接
// @Description 根据文件名生成私有或公共的下载链接
// @Tags 文件管理
// @Accept json
// @Produce json
// @Param objectName query string true "文件名"
// @Param accessType query string false "访问类型 ('public' 或 'private')" 默认 "private"
// @Success 200 {object} map[string]interface{} "生成下载链接成功，返回下载链接"
// @Failure 400 {object} map[string]interface{} "缺少必要参数 objectName"
// @Router /api/v1/download [get]
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

// DeleteFileHandler 删除文件接口
// @Summary 删除文件
// @Description 根据文件名删除七牛云中的文件
// @Tags 文件管理
// @Accept json
// @Produce json
// @Param objectName query string true "文件名"
// @Success 200 {object} map[string]interface{} "文件删除成功"
// @Failure 400 {object} map[string]interface{} "缺少必要参数 objectName"
// @Failure 500 {object} map[string]interface{} "文件删除失败"
// @Router /api/v1/delete [delete]
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

// ListFilesHandler 获取文件列表接口
// @Summary 获取文件列表
// @Description 列出七牛云存储空间中的文件
// @Tags 文件管理
// @Accept json
// @Produce json
// @Param prefix query string false "文件名前缀筛选条件"
// @Param marker query string false "游标，继续从上次读取的标记处开始列出"
// @Param limit query int false "每次列举的最大文件数量 (1-1000)"
// @Success 200 {object} map[string]interface{} "文件列表获取成功，返回文件信息及下一页游标"
// @Failure 500 {object} map[string]interface{} "文件列表获取失败"
// @Router /api/v1/list [get]
func ListFilesHandler(c *gin.Context) {
	// 获取请求参数
	prefix := c.DefaultQuery("prefix", "") // 文件前缀
	marker := c.DefaultQuery("marker", "") // 游标，列举时继续读取上次的 marker
	limit := 1000                          // 默认每次最多列举 1000 个文件
	if c.Query("limit") != "" {
		// 如果有指定 limit，转换为整数
		parsedLimit, err := strconv.Atoi(c.Query("limit"))
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

	// Map original file list to a limited field response
	var limitedFiles []model.FileInfo
	for _, file := range files {
		limitedFiles = append(limitedFiles, model.FileInfo{
			Key:           file.Key,
			ContentLength: file.Fsize,
			ETag:          file.Hash,
			LastModified:  time.Unix(file.PutTime/1e7, 0).UTC(), // Convert timestamp to time.Time
		})
	}

	// 返回文件列表和下一页游标
	c.JSON(http.StatusOK, gin.H{
		"code":        http.StatusOK,
		"msg":         "文件列表获取成功",
		"files":       limitedFiles,
		"next_marker": nextMarker,
	})
}

// CopyFileHandler 复制文件接口
// @Summary 复制文件
// @Description 将七牛云存储空间中的文件从一个位置复制到另一个位置
// @Tags 文件管理
// @Accept json
// @Produce json
// @Param srcObject query string true "源文件名"
// @Param destObject query string true "目标文件名"
// @Param force query bool false "是否强制覆盖目标文件（true/false，默认为 false）"
// @Success 200 {object} map[string]interface{} "文件复制成功"
// @Failure 400 {object} map[string]interface{} "srcKey、destKey 缺失或 force 参数无效"
// @Failure 500 {object} map[string]interface{} "文件复制失败"
// @Router /api/v1/copy [post]
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

// MoveFileHandler 移动文件接口
// @Summary 移动文件
// @Description 将七牛云存储空间中的文件从一个位置移动到另一个位置
// @Tags 文件管理
// @Accept json
// @Produce json
// @Param srcObject query string true "源文件名"
// @Param destObject query string true "目标文件名"
// @Param force query bool false "是否强制覆盖目标文件（true/false，默认为 false）"
// @Success 200 {object} map[string]interface{} "文件移动成功"
// @Failure 400 {object} map[string]interface{} "srcKey、destKey 缺失或 force 参数无效"
// @Failure 500 {object} map[string]interface{} "文件移动失败"
// @Router /api/v1/move [post]
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
