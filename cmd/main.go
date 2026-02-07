package main

import (
	"go-todo-api/config"
	"go-todo-api/pkg/database"
	"go-todo-api/routes"
	"log"
)

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
