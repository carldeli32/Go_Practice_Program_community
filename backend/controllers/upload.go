package controllers

import (
	"net/http"
	"strconv"

	"community/backend/models"
	"community/backend/service"

	"github.com/gin-gonic/gin"
)

// ─── 上传图片 ───
// POST /api/upload/image  (需登录)
// Form: file + post_id
func UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		models.Error(c, http.StatusBadRequest, "请选择图片文件")
		return
	}

	postID, _ := strconv.Atoi(c.PostForm("post_id"))
	url, err := service.UploadImage(uint(postID), file)
	if err != nil {
		code, msg := service.ToHTTP(err)
		models.Error(c, code, msg)
		return
	}

	models.Success(c, "上传成功 🖼️", gin.H{
		"url":      url,
		"filename": file.Filename,
		"size":     file.Size,
	})
}

// ─── 上传文件 ───
// POST /api/upload/file  (需登录)
// Form: file + post_id
func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		models.Error(c, http.StatusBadRequest, "请选择文件")
		return
	}

	postID, _ := strconv.Atoi(c.PostForm("post_id"))
	url, err := service.UploadFile(uint(postID), file)
	if err != nil {
		code, msg := service.ToHTTP(err)
		models.Error(c, code, msg)
		return
	}

	models.Success(c, "上传成功 📎", gin.H{
		"url":      url,
		"filename": file.Filename,
		"size":     file.Size,
	})
}
