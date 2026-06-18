package storage

import "io"

type Storage interface {
	// SaveImage 保存图片到指定帖子目录，返回公开访问 URL
	SaveImage(postID uint, filename string, reader io.Reader) (url string, err error)
	// SaveFile 保存普通文件到指定帖子目录
	SaveFile(postID uint, filename string, reader io.Reader) (url string, err error)
	// DeleteImage 删除图片（按完整 URL 路径）
	DeleteImage(urlPath string) error
	// DeleteFile 删除文件（按完整 URL 路径）
	DeleteFile(urlPath string) error
	// DeletePostDir 删除指定帖子的整个上传目录
	DeletePostDir(postID uint) error
}

var Store Storage
