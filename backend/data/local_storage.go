package data

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"community/backend/storage"
)

var _ storage.Storage = (*LocalStorage)(nil)

type LocalStorage struct {
	ImageBasePath string
	FileBasePath  string
	ImageBaseURL  string
	FileBaseURL   string
}

func NewLocalStorage(imagePath, filePath, imageURL, fileURL string) *LocalStorage {
	return &LocalStorage{
		ImageBasePath: imagePath,
		FileBasePath:  filePath,
		ImageBaseURL:  imageURL,
		FileBaseURL:   fileURL,
	}
}

// ─── 保存到帖子目录 ───

func (s *LocalStorage) SaveImage(postID uint, filename string, reader io.Reader) (string, error) {
	return s.saveForPost(s.ImageBasePath, s.ImageBaseURL, postID, filename, reader)
}

func (s *LocalStorage) SaveFile(postID uint, filename string, reader io.Reader) (string, error) {
	return s.saveForPost(s.FileBasePath, s.FileBaseURL, postID, filename, reader)
}

func (s *LocalStorage) saveForPost(basePath, baseURL string, postID uint, filename string, reader io.Reader) (string, error) {
	subDir := strconv.FormatUint(uint64(postID), 10)
	fullDir := filepath.Join(basePath, subDir)
	if err := os.MkdirAll(fullDir, 0755); err != nil {
		return "", fmt.Errorf("创建存储目录失败: %w", err)
	}

	fullPath := filepath.Join(fullDir, filename)
	f, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("创建文件失败: %w", err)
	}
	defer f.Close()

	if _, err := io.Copy(f, reader); err != nil {
		os.Remove(fullPath)
		return "", fmt.Errorf("写入文件失败: %w", err)
	}

	return filepath.Join(baseURL, subDir, filename), nil
}

// ─── 删除 ───

func (s *LocalStorage) DeleteImage(urlPath string) error {
	return s.deleteByURL(s.ImageBasePath, s.ImageBaseURL, urlPath)
}

func (s *LocalStorage) DeleteFile(urlPath string) error {
	return s.deleteByURL(s.FileBasePath, s.FileBaseURL, urlPath)
}

func (s *LocalStorage) deleteByURL(basePath, baseURL, urlPath string) error {
	rel := urlPath
	if len(urlPath) > len(baseURL) && urlPath[:len(baseURL)] == baseURL {
		rel = urlPath[len(baseURL)+1:]
	}
	fullPath := filepath.Join(basePath, rel)
	return os.Remove(fullPath)
}

// ─── 帖子目录管理 ───

func (s *LocalStorage) DeletePostDir(postID uint) error {
	cid := strconv.FormatUint(uint64(postID), 10)
	_ = os.RemoveAll(filepath.Join(s.ImageBasePath, cid))
	_ = os.RemoveAll(filepath.Join(s.FileBasePath, cid))
	return nil
}
