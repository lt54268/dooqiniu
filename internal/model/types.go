package model

import (
	"mime/multipart"
	"time"
)

// Config 用于存储配置信息
type Config struct {
	Port           string
	QiniuRegion    string
	QiniuEndpoint  string
	QiniuBucket    string
	QiniuSecretId  string
	QiniuSecretKey string
}

// Uploader 定义上传接口
type Uploader interface {
	Upload(file multipart.File, objectName string) (string, error)
}

// FileInfo 包含文件基本信息
type FileInfo struct {
	Key           string    `json:"key"`
	ContentLength int64     `json:"content-length"`
	ETag          string    `json:"etag"`
	LastModified  time.Time `json:"last_modified"`
}

// UploadResponse 包含上传文件后的响应信息
type UploadResponse struct {
	ContentLength int64     `json:"content-length"`
	ETag          string    `json:"etag"`
	LastModified  time.Time `json:"last-modified"`
}
