package storage

import "io"

// Storage 文件存储抽象接口。
// 当前使用本地磁盘实现，将来可替换为 S3/OSS/COS 等云存储，
// 只需实现本接口，上层 service 代码无需改动。
type Storage interface {
	// SaveImage 保存图片，返回公开访问 URL
	SaveImage(filename string, reader io.Reader) (url string, err error)
	// SaveFile 保存普通文件，返回公开访问 URL
	SaveFile(filename string, reader io.Reader) (url string, err error)
	// DeleteImage 删除图片
	DeleteImage(path string) error
	// DeleteFile 删除文件
	DeleteFile(path string) error
}

// Store 全局存储实例，由 main.go 初始化
var Store Storage
