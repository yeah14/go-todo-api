package middleware

import (
	"go-todo-api/internal/service"
	"go-todo-api/pkg/jwt"
	"go-todo-api/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthHandler(blacklistSvc service.BlacklistService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 从Header中获取Token
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

		//2,检查黑名单
		inBlackList, err := blacklistSvc.IsInBlacklist(c.Request.Context(), tokenString)
		if err != nil {
			response.Unauthorized(c, err.Error())
			c.Abort()
			return
		}
		if inBlackList {
			response.Unauthorized(c, "令牌状态验证失败")
			c.Abort()
			return
		}

		// 3. 解析并验证JWT令牌
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
		c.Set("tokenString", tokenString)
		c.Next()
	}
}
