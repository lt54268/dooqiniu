package service

import (
	"context"
	"dooqiniu/internal/model"
	"fmt"
	"os"
	"path/filepath"

	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/http_client"
	"github.com/qiniu/go-sdk/v7/storagev2/uploader"
)

type QiniuCommoner struct {
	accessKey  string
	secretKey  string
	bucketName string
}

func NewQiniuClient() *QiniuCommoner {
	return &QiniuCommoner{
		accessKey:  os.Getenv("QINIU_ACCESSKEY"),
		secretKey:  os.Getenv("QINIU_SECRETKEY"),
		bucketName: os.Getenv("QINIU_BUCKET"),
	}
}

// Upload uploads a file to Qiniu Cloud
func (q *QiniuCommoner) Upload(filePath, objectName string) (*model.UploadResponse, error) {
	// Initialize credentials
	mac := credentials.NewCredentials(q.accessKey, q.secretKey)

	// Create upload manager with credentials
	options := uploader.UploadManagerOptions{
		Options: http_client.Options{
			Credentials: mac,
		},
	}
	uploadManager := uploader.NewUploadManager(&options)

	// Open the file to upload
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Ensure file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %s", filePath)
	}

	// Set up object options with target path
	objectOptions := &uploader.ObjectOptions{
		BucketName: q.bucketName,
		ObjectName: &objectName,
		FileName:   filepath.Base(filePath),
	}

	// Perform upload
	err = uploadManager.UploadFile(context.Background(), filePath, objectOptions, nil)
	if err != nil {
		return nil, fmt.Errorf("upload failed: %w", err)
	}
	model.FileInfo := files[0]

	return &model.UploadResponse{
		Fsize:   fileInfo.Fsize,   // 文件大小
		ETag:    fileInfo.Hash,    // 文件的 ETag
		PutTime: fileInfo.PutTime, // 上传时间
	}, nil

}

// ListFiles 列出七牛云桶中的文件
func (q *QiniuCommoner) ListFiles(prefix, marker string, limit int) ([]storage.ListItem, string, error) {
	// 初始化认证
	mac := auth.New(q.accessKey, q.secretKey)

	// 创建 bucketManager 实例
	bucketManager := storage.NewBucketManager(mac, &storage.Config{})

	// 获取文件列表
	entries, _, nextMarker, hasNext, err := bucketManager.ListFiles(q.bucketName, prefix, "", marker, limit)
	if err != nil {
		return nil, "", fmt.Errorf("failed to list files: %v", err)
	}

	// 这里我们只返回文件项列表和下一页的 marker
	if !hasNext {
		nextMarker = ""
	}

	// 返回文件列表和下一页的游标
	return entries, nextMarker, nil
}
