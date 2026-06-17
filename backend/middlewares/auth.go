package middlewares

import (
	"net/http"
	"strings"
	"time"

	"community/backend/config"
	"community/backend/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTSecret 签名密钥（生产环境应该放环境变量）
var JWTSecret = []byte("community-secret-key-2024")

// Claims 自定义 JWT 载荷
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT Token，有效期 7 天
func GenerateToken(userID uint, username string) (string, error) {
	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWTSecret)
}

// AuthRequired JWT 鉴权中间件（加载角色 + 封禁检查）
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			models.Error(c, http.StatusUnauthorized, "请先登录")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			models.Error(c, http.StatusUnauthorized, "Token 格式错误")
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return JWTSecret, nil
		})

		if err != nil || !token.Valid {
			models.Error(c, http.StatusUnauthorized, "Token 无效或已过期")
			c.Abort()
			return
		}

		// 查出用户 + 角色
		var user models.User
		if err := config.DB.Preload("Roles").First(&user, claims.UserID).Error; err != nil {
			models.Error(c, http.StatusUnauthorized, "用户不存在，请重新登录")
			c.Abort()
			return
		}

		// 全站封禁检查
		if user.IsBanned {
			models.Error(c, http.StatusForbidden, "账号已被封禁")
			c.Abort()
			return
		}

		// 注入上下文
		c.Set("user_id", user.ID)
		c.Set("username", claims.Username)
		c.Set("roles", user.RoleNames())
		c.Set("is_admin", user.IsAdminLike()) // 前端兼容
		c.Next()
	}
}
