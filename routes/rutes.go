package routes

import (
	"go-todo-api/internal/app/handler"
	"go-todo-api/internal/app/middleware"
	"go-todo-api/internal/repository"
	"go-todo-api/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	//初始化仓库
	userReop := repository.NewUserRepository(db)
	todoRepo := repository.NewTodoRepo(db)
	tagReop := repository.NewTagRepository(db)

	//初始化服务
	authService := service.NewAuthService(userReop)
	userService := service.NewUserService(userReop)
	todoService := service.NewTodoService(todoRepo, tagReop)
	tagService := service.NewTagService(tagReop)
	//初始化处理器
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	todoHandler := handler.NewTodoHandler(todoService)
	tagHandler := handler.NewTagHandler(tagService)
	pubilc := r.Group("/api/v1")
	{
		pubilc.POST("auth/register", authHandler.Register)
		pubilc.POST("auth/login", authHandler.Login)
	}
	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthHandler())
	{
		user := protected.Group("/user")
		{
			user.GET("/me", userHandler.GetProfile)
			user.POST("/me", userHandler.UpdateProfile)
			user.POST("/me/password", userHandler.ChangePassword)
		}

		todos := protected.Group("/todos")
		{
			todos.GET("/:id", todoHandler.GetTodoById)
			todos.POST("", todoHandler.CreateTodo)
			todos.DELETE("/:id", todoHandler.Delete)
			todos.PUT("/:id", todoHandler.Update)
			todos.GET("", todoHandler.GetTodos)
			todos.PUT("/batch/status", todoHandler.BatchUpdateStatus)
		}

		tags := protected.Group("/tags")
		{
			tags.GET("", tagHandler.GetTags)
			tags.POST("", tagHandler.Create)
			tags.PUT("/:id", tagHandler.Update)
			tags.DELETE("/:id", tagHandler.Delete)
		}
	}
	return r
}
