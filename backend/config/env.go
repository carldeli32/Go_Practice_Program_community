package config

import (
	"log"
	"os"
	"strings"

	"community/backend/models"

	"github.com/joho/godotenv"
)

// 环境变量键名常量
const (
	EnvDBDSN             = "DB_DSN"
	EnvJWTSecret         = "JWT_SECRET"
	EnvPort              = "PORT"
	EnvUploadImg         = "UPLOAD_IMG_DIR"
	EnvUploadFile        = "UPLOAD_FILE_DIR"
	EnvServeImg          = "SERVE_IMG_PREFIX"
	EnvServeFile         = "SERVE_FILE_PREFIX"
	EnvRootPassword      = "ROOT_PASSWORD"
	EnvRootPasswordFile  = "ROOT_PASSWORD_FILE"
	EnvDBMaxOpenConns    = "DB_MAX_OPEN_CONNS"
	EnvDBMaxIdleConns    = "DB_MAX_IDLE_CONNS"
	EnvDBConnMaxLifetime = "DB_CONN_MAX_LIFETIME"
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

// RootPassword 获取 root 初始密码（仅首次部署需要）
// 优先级：ROOT_PASSWORD_FILE > ROOT_PASSWORD
// 非首次部署（root 用户已存在）返回空字符串
func RootPassword() string {
	// 1. 文件优先（Docker secrets / K8s secrets）
	if file := os.Getenv(EnvRootPasswordFile); file != "" {
		data, err := os.ReadFile(file)
		if err != nil {
			panic("无法读取 " + EnvRootPasswordFile + ": " + err.Error())
		}
		return strings.TrimSpace(string(data))
	}

	// 2. 直接环境变量
	if pw := os.Getenv(EnvRootPassword); pw != "" {
		return pw
	}

	// 3. 非首次部署（root 用户已存在）不需要密码
	if DB != nil {
		var count int64
		DB.Model(&models.User{}).Where("username = ?", "root").Count(&count)
		if count > 0 {
			return ""
		}
	}

	// 4. 首次部署但没设密码
	panic("首次部署需要设置 ROOT_PASSWORD 或 ROOT_PASSWORD_FILE 环境变量")
}

// DBMaxOpenConns 返回数据库最大连接数
func DBMaxOpenConns() int {
	return parseIntEnv(EnvDBMaxOpenConns, 25)
}

// DBMaxIdleConns 返回数据库最大空闲连接数
func DBMaxIdleConns() int {
	return parseIntEnv(EnvDBMaxIdleConns, 10)
}

// DBConnMaxLifetime 返回连接最大存活秒数
func DBConnMaxLifetime() int {
	return parseIntEnv(EnvDBConnMaxLifetime, 300)
}

// parseIntEnv 解析整数环境变量，失败返回默认值
func parseIntEnv(key string, defaultVal int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	n := 0
	for _, c := range val {
		if c < '0' || c > '9' {
			return defaultVal
		}
		n = n*10 + int(c-'0')
	}
	return n
}
