package main

import (
	"go-todo-api/config"
	"go-todo-api/pkg/database"
	"go-todo-api/routes"
	"log"
)

// @title           Go-Todo-API
// @version         1.0
// @description     基于 Gin 的待办事项管理系统 API
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 Type "Bearer" followed by a space and JWT token.
func main() {
	config.InitConfig("config/dev.yaml")
	db, err := database.ConnectMysql(config.GlobalConfig.Database)
	if err != nil {
		log.Fatalf("数据库链接失败 %v", err)
	}
	r := routes.SetupRouter(db)
	port := config.GlobalConfig.Server.Port
	log.Printf("服务启动在:https://localhost:%s\n", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("服务启动失败： %v", err)
	}
}
