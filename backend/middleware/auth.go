package middleware

import (
	"e-device-recycle-backend/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWT中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少认证令牌"})
			c.Abort()
			return
		}

		// 检查Bearer前缀
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "认证令牌格式错误"})
			c.Abort()
			return
		}

		// 验证token
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的认证令牌"})
			c.Abort()
			return
		}

		// 将用户信息保存到上下文
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// 管理员权限中间件
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "需要管理员权限"})
			c.Abort()
			return
		}
		c.Next()
	}
}

// CORS中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
