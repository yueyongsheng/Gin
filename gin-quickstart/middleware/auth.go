package middleware

import (
	"net/http"
	"strings"

	"gin-quickstart/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWT 密钥
const SECRET_KEY = "565962F5-970E-4F85-9B30-7DB65CFA8F07"

// AuthMiddleware JWT 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.HandleError(c, http.StatusUnauthorized, "需要授权头", nil)
			c.Abort()
			return
		}

		// 解析 Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.HandleError(c, http.StatusUnauthorized, "无效的授权格式,应为 'Bearer <token>'", nil)
			c.Abort()
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		})

		if err != nil {
			utils.HandleError(c, http.StatusUnauthorized, "令牌解析失败", err)
			c.Abort()
			return
		}

		if !token.Valid {
			utils.HandleError(c, http.StatusUnauthorized, "令牌无效或已过期", nil)
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("user_id", uint(claims["id"].(float64)))
			c.Set("username", claims["username"].(string))
		} else {
			utils.HandleError(c, http.StatusUnauthorized, "令牌内容格式错误", nil)
			c.Abort()
			return
		}

		c.Next()
	}
}
