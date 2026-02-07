package routes

import (
	"go-todo-api/internal/app/handler"
	"go-todo-api/internal/repository"
	"go-todo-api/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	//初始化仓库
	userReop := repository.NewUserRepository(db)

	//初始化服务
	authService := service.NewAuthService(userReop)

	//初始化处理器
	authHandler := handler.NewAuthHandler(authService)
	pubilc := r.Group("/api/v1")
	{
		pubilc.POST("auth/register", authHandler.Register)
		pubilc.POST("auth/login", authHandler.Login)
	}
	return r
}
