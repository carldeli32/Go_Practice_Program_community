package controllers

import (
	"net/http"
	"strconv"

	"community/backend/data"
	"community/backend/models"
	"community/backend/service"

	"github.com/gin-gonic/gin"
)

// ─── 创建帖子 ───
func CreatePost(c *gin.Context) {
	var req struct {
		Title      string `json:"title"`
		Content    string `json:"content"`
		CategoryID *uint  `json:"category_id"`
		Status     string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		models.Error(c, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	categoryID := uint(1)
	if req.CategoryID != nil {
		categoryID = *req.CategoryID
	}

	userID := c.GetUint("user_id")
	status := req.Status
	if status == "" {
		status = "published"
	}
	post, err := service.CreatePost(userID, req.Title, req.Content, categoryID, status)
	if err != nil {
		code, msg := service.ToHTTP(err)
		models.Error(c, code, msg)
		return
	}

	models.Success(c, "发帖成功 ✍️", gin.H{"post": post})
}

// ─── 帖子列表 ───
func GetPosts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	posts, total, err := data.ListPosts(page, pageSize,
		c.DefaultQuery("category_id", ""),
		c.DefaultQuery("q", ""))
	if err != nil {
		models.Error(c, http.StatusInternalServerError, "查询失败")
		return
	}

	models.Success(c, "获取成功", gin.H{
		"posts":     posts,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// ─── 帖子详情 ───
func GetPost(c *gin.Context) {
	pid := c.Param("id")
	post, err := data.FindPostByID(strToUint(pid))
	if err != nil {
		models.Error(c, http.StatusNotFound, "帖子不存在")
		return
	}
	models.Success(c, "获取成功", gin.H{"post": post})
}

// ─── 更新帖子 ───
func UpdatePost(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetUint("user_id")
	roles := c.GetStringSlice("roles")

	var req struct {
		Title      string `json:"title"`
		Content    string `json:"content"`
		CategoryID *uint  `json:"category_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		models.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	categoryID := uint(0)
	if req.CategoryID != nil {
		categoryID = *req.CategoryID
	}

	post, err := service.UpdatePost(strToUint(id), userID, req.Title, req.Content, categoryID, roles)
	if err != nil {
		code, msg := service.ToHTTP(err)
		models.Error(c, code, msg)
		return
	}

	models.Success(c, "更新成功", gin.H{"post": post})
}

// ─── 删除帖子 ───
func DeletePost(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetUint("user_id")
	roles := c.GetStringSlice("roles")

	if err := service.DeletePost(strToUint(id), userID, roles); err != nil {
		code, msg := service.ToHTTP(err)
		models.Error(c, code, msg)
		return
	}

	models.Success(c, "删除成功", nil)
}

func strToUint(s string) uint {
	v, _ := strconv.Atoi(s)
	return uint(v)
}
