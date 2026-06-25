package middlewares

import (
	"strconv"

	"community/backend/data"

	"github.com/gin-gonic/gin"
)

// PostLoader 加载帖子资源：返回 ownerID + categoryID + Post 对象
func PostLoader(c *gin.Context) (ownerID, categoryID uint, resource interface{}, err error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return 0, 0, nil, err
	}

	post, err := data.FindPostByIDRaw(uint(id))
	if err != nil {
		return 0, 0, nil, err
	}
	return post.UserID, post.CategoryID, post, nil
}

// CommentLoader 加载评论资源：返回 ownerID + categoryID（通过 JOIN post 获取）+ Comment 对象
func CommentLoader(c *gin.Context) (ownerID, categoryID uint, resource interface{}, err error) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return 0, 0, nil, err
	}

	result, err := data.FindCommentByIDWithCategory(uint(id))
	if err != nil {
		return 0, 0, nil, err
	}
	return result.Comment.UserID, result.CategoryID, &result.Comment, nil
}

// StrToUint 字符串转 uint，解析失败返回 0
func StrToUint(s string) uint {
	v, _ := strconv.Atoi(s)
	return uint(v)
}
