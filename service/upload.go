package service

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"community/storage"

	"github.com/google/uuid"
)

// ─── 限制配置 ───

const (
	MaxImageSize = 5 << 20  // 5MB
	MaxFileSize  = 20 << 20 // 20MB
)

var allowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

// ─── 上传图片 ───

// UploadImage 校验 + 保存图片，返回公开 URL
func UploadImage(fileHeader *multipart.FileHeader) (string, error) {
	if fileHeader.Size > MaxImageSize {
		return "", &AppError{Code: http.StatusBadRequest, Message: "图片不能超过 5MB"}
	}

	// MIME 校验
	contentType := fileHeader.Header.Get("Content-Type")
	if !allowedImageTypes[contentType] {
		return "", &AppError{Code: http.StatusBadRequest,
			Message: fmt.Sprintf("不支持的图片格式: %s，仅支持 jpg/png/gif/webp", contentType)}
	}

	// 魔数校验：读前 512 字节检测真实类型
	file, err := fileHeader.Open()
	if err != nil {
		return "", &AppError{Code: http.StatusInternalServerError, Message: "读取文件失败"}
	}
	defer file.Close()

	if !isRealImage(file) {
		return "", &AppError{Code: http.StatusBadRequest, Message: "文件内容与扩展名不匹配"}
	}

	// 重置读指针
	file.Seek(0, io.SeekStart)

	// 生成随机文件名
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext == "" || ext == ".jpeg" {
		ext = ".jpg"
	}
	filename := uuid.New().String() + ext

	return storage.Store.SaveImage(filename, file)
}

// ─── 上传文件 ───

// UploadFile 校验 + 保存普通文件
func UploadFile(fileHeader *multipart.FileHeader) (string, error) {
	if fileHeader.Size > MaxFileSize {
		return "", &AppError{Code: http.StatusBadRequest, Message: "文件不能超过 20MB"}
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", &AppError{Code: http.StatusInternalServerError, Message: "读取文件失败"}
	}
	defer file.Close()

	// 生成安全文件名（保留原始扩展名，但文件名用 UUID）
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	safeName := uuid.New().String() + ext

	return storage.Store.SaveFile(safeName, file)
}

// ─── 魔数检测 ───

// isRealImage 读取文件头几个字节判断真实图片类型
func isRealImage(r io.ReadSeeker) bool {
	header := make([]byte, 12)
	n, _ := io.ReadFull(r, header)
	if n < 4 {
		return false
	}

	// JPEG: FF D8 FF
	if header[0] == 0xFF && header[1] == 0xD8 && header[2] == 0xFF {
		return true
	}
	// PNG: 89 50 4E 47
	if header[0] == 0x89 && header[1] == 0x50 && header[2] == 0x4E && header[3] == 0x47 {
		return true
	}
	// GIF: 47 49 46 38
	if header[0] == 0x47 && header[1] == 0x49 && header[2] == 0x46 && header[3] == 0x38 {
		return true
	}
	// WebP: 52 49 46 46 ... 57 45 42 50
	if n >= 12 && header[0] == 0x52 && header[1] == 0x49 && header[2] == 0x46 && header[3] == 0x46 &&
		header[8] == 0x57 && header[9] == 0x45 && header[10] == 0x42 && header[11] == 0x50 {
		return true
	}
	return false
}
