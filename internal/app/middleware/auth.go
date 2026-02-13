package middleware

import (
	"go-todo-api/pkg/jwt"
	"go-todo-api/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		//fmt.Println("Authorization header:", authHeader) //test
		if authHeader == "" {
			response.Unauthorized(c, "缺少认证令牌")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Unauthorized(c, "认证令牌格式错误")
			c.Abort()
			return
		}

		tokenString := parts[1]
		//fmt.Println(tokenString)//test
		claims, err := jwt.ParseToken(tokenString)
		if err != nil {
			//fmt.Println("AuthHandler ParseToken error:", err) // 新增
			response.Unauthorized(c, "认证令牌无效或已过期")
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
