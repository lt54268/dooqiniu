package config

import (
	"dooqiniu/internal/model"
	"os"
)

// 从环境变量加载配置信息
func LoadQiniuConfig() *model.Config {
	return &model.Config{
		Port:           os.Getenv("PORT"),            // 从环境变量读取端口
		QiniuRegion:    os.Getenv("QINIU_REGION"),    // 从环境变量读取区域
		QiniuEndpoint:  os.Getenv("QINIU_ENDPOINT"),  // 从环境变量读取 Endpoint
		QiniuBucket:    os.Getenv("QINIU_BUCKET"),    // 从环境变量读取 Bucket
		QiniuSecretId:  os.Getenv("QINIU_SECRETID"),  // 从环境变量读取 AccessKeyId
		QiniuSecretKey: os.Getenv("QINIU_SECRETKEY"), // 从环境变量读取 AccessKeySecret
	}
}
