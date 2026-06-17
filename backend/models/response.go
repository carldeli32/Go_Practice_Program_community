package models

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一返回结构体
// Code: 0 表示成功，非 0 表示错误（值即 HTTP 状态码）
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 返回成功响应（HTTP 200, code=0）
func Success(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: message,
		Data:    data,
	})
}

// Error 返回错误响应（HTTP 状态码即 code）
func Error(c *gin.Context, httpStatus int, message string) {
	c.JSON(httpStatus, Response{
		Code:    httpStatus,
		Message: message,
	})
}
