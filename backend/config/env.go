package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// 环境变量键名常量
const (
	EnvDBDSN      = "DB_DSN"
	EnvJWTSecret  = "JWT_SECRET"
	EnvPort       = "PORT"
	EnvUploadImg  = "UPLOAD_IMG_DIR"
	EnvUploadFile = "UPLOAD_FILE_DIR"
	EnvServeImg   = "SERVE_IMG_PREFIX"
	EnvServeFile  = "SERVE_FILE_PREFIX"
)

// LoadEnv 加载 .env 文件
// 优先级：backend/.env > ../.env（项目根目录）
func LoadEnv() {
	if err := godotenv.Load(".env"); err != nil {
		// 尝试项目根目录
		if err2 := godotenv.Load("../.env"); err2 != nil {
			log.Println("⚠️  未找到 .env 文件，使用系统环境变量或默认值")
		}
	}
}

// GetEnv 获取环境变量，如未设置返回默认值
func GetEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// DSN 返回数据库连接串
func DSN() string {
	return GetEnv(EnvDBDSN,
		"root:password@tcp(127.0.0.1:3306)/community?charset=utf8mb4&parseTime=True")
}

// JWTKey 返回 JWT 签名密钥
func JWTKey() []byte {
	return []byte(GetEnv(EnvJWTSecret, "community-secret-key-change-me"))
}

// ServerPort 返回服务监听端口（含冒号前缀）
func ServerPort() string {
	port := GetEnv(EnvPort, "8080")
	return ":" + port
}

// UploadImgDir 返回图片上传目录
func UploadImgDir() string {
	return GetEnv(EnvUploadImg, "../uploads/images")
}

// UploadFileDir 返回文件上传目录
func UploadFileDir() string {
	return GetEnv(EnvUploadFile, "../uploads/files")
}

// ServeImgPrefix 返回图片访问 URL 前缀
func ServeImgPrefix() string {
	return GetEnv(EnvServeImg, "/uploads/images")
}

// ServeFilePrefix 返回文件访问 URL 前缀
func ServeFilePrefix() string {
	return GetEnv(EnvServeFile, "/uploads/files")
}
