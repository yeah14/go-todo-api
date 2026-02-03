package routes

import (
	"go-todo-api/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	userReop := repository.NewUserRepository(db)
	return r
}
