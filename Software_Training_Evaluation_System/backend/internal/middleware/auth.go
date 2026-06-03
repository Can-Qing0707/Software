package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	jwtpkg "github.com/golang-jwt/jwt/v5"

	"training_eval_system/config"
	"training_eval_system/pkg/jwt"
	"training_eval_system/pkg/response"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "未提供认证令牌")
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "认证格式错误")
			c.Abort()
			return
		}
		claims, err := jwt.ParseToken(config.AppConfig.JWT.Secret, parts[1])
		if err != nil {
			response.Unauthorized(c, "令牌无效或已过期")
			c.Abort()
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func RoleRequired(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, _ := c.Get("role")
		roleStr, _ := userRole.(string)
		for _, r := range roles {
			if r == roleStr {
				c.Next()
				return
			}
		}
		response.Forbidden(c, "权限不足")
		c.Abort()
	}
}

func GetUserID(c *gin.Context) uint {
	id, exists := c.Get("user_id")
	if !exists {
		return 0
	}
	return id.(uint)
}

func GetRole(c *gin.Context) string {
	role, exists := c.Get("role")
	if !exists {
		return ""
	}
	return role.(string)
}

func IsAdmin(c *gin.Context) bool {
	return GetRole(c) == "admin"
}

func IsTeacher(c *gin.Context) bool {
	return GetRole(c) == "teacher"
}

func IsStudent(c *gin.Context) bool {
	return GetRole(c) == "student"
}

var _ = jwtpkg.ErrInvalidKey
