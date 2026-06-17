package data

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"community/storage"
)

// 编译期检查：确保 LocalStorage 实现了 storage.Storage
var _ storage.Storage = (*LocalStorage)(nil)

// LocalStorage 将图片和文件分别存到不同目录
type LocalStorage struct {
	ImageBasePath string // 如 "./uploads/images"
	FileBasePath  string // 如 "./uploads/files"
	ImageBaseURL  string // 如 "/uploads/images"
	FileBaseURL   string // 如 "/uploads/files"
}

// NewLocalStorage 创建本地存储实例
func NewLocalStorage(imagePath, filePath, imageURL, fileURL string) *LocalStorage {
	return &LocalStorage{
		ImageBasePath: imagePath,
		FileBasePath:  filePath,
		ImageBaseURL:  imageURL,
		FileBaseURL:   fileURL,
	}
}

func (s *LocalStorage) SaveImage(filename string, reader io.Reader) (string, error) {
	return s.save(s.ImageBasePath, s.ImageBaseURL, filename, reader)
}

func (s *LocalStorage) SaveFile(filename string, reader io.Reader) (string, error) {
	return s.save(s.FileBasePath, s.FileBaseURL, filename, reader)
}

func (s *LocalStorage) DeleteImage(path string) error {
	return s.deleteFile(s.ImageBasePath, path)
}

func (s *LocalStorage) DeleteFile(path string) error {
	return s.deleteFile(s.FileBasePath, path)
}

// ─── 内部实现 ───

func (s *LocalStorage) save(basePath, baseURL, filename string, reader io.Reader) (string, error) {
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return "", fmt.Errorf("创建存储目录失败: %w", err)
	}

	fullPath := filepath.Join(basePath, filename)
	f, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("创建文件失败: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, reader); err != nil {
		os.Remove(fullPath)
		return "", fmt.Errorf("写入文件失败: %w", err)
	}

	return baseURL + "/" + filename, nil
}

func (s *LocalStorage) deleteFile(basePath, path string) error {
	filename := filepath.Base(path)
	fullPath := filepath.Join(basePath, filename)
	return os.Remove(fullPath)
}
