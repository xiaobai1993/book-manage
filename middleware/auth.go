package middleware

import (
	"book-manage/utils"
	"bytes"
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// AuthMiddleware 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从Header、Query参数或请求体中获取token
		var token string

		// 优先从Header获取
		authHeader := c.GetHeader("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token = parts[1]
			}
		}

		// 如果Header中没有，尝试从Query参数获取
		if token == "" {
			token = c.Query("token")
		}

		// 如果Query中也没有，尝试从请求体获取（不消耗请求体）
		if token == "" && c.Request.Body != nil {
			// 读取请求体
			bodyBytes, err := io.ReadAll(c.Request.Body)
			if err == nil {
				// 恢复请求体
				c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
				
				// 解析JSON获取token
				var jsonData map[string]interface{}
				if json.Unmarshal(bodyBytes, &jsonData) == nil {
					if tokenVal, ok := jsonData["token"].(string); ok && tokenVal != "" {
						token = tokenVal
					}
				}
			}
		}

		if token == "" {
			utils.Error(c, 10001, "缺少token")
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(token)
		if err != nil {
			utils.Error(c, 10001, "token无效或已过期")
			c.Abort()
			return
		}

		// 将用户信息存储到上下文
		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)

		c.Next()
	}
}

// AdminMiddleware 管理员权限中间件
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("user_role")
		if !exists || role != "admin" {
			utils.Error(c, 10009, "权限不足")
			c.Abort()
			return
		}
		c.Next()
	}
}
