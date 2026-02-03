package database

import (
	"go-todo-api/config"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectMysql(t *testing.T) {
	//cfg := config.DatabaseConfig{
	//	Host:     "localhost",
	//	Port:     "3306",
	//	Username: "root",
	//	Password: "123456",
	//	DBname:   "todo_db",
	//	Charset:  "utf8mb4",
	//}
	config.InitConfig("../../config/dev.yaml")
	db, err := ConnectMysql(config.GlobalConfig.Database)
	assert.NoError(t, err, "数据库链接成功")
	assert.NotNil(t, db, "数据库的对象不为空")

	if db != nil {
		var result int
		err = db.Raw("SELECT 1").Scan(&result).Error
		assert.NoError(t, err, "应该能执行SQL查询")
		assert.Equal(t, 1, result, "查询结果应该是1")

		// 关闭连接
		sqlDB, _ := db.DB()
		sqlDB.Close()
	}
}

func TestConnectMySQL_Failure(t *testing.T) {
	// 使用错误配置测试连接失败
	cfg := config.DatabaseConfig{
		Host:     "wrong_host",
		Port:     "3306",
		User:     "wrong_user",
		Password: "wrong_password",
		DBName:   "wrong_db",
		Charset:  "utf8mb4",
	}

	db, err := ConnectMysql(cfg)
	assert.Error(t, err, "应该返回连接错误")
	assert.Nil(t, db, "数据库连接对象应该为nil")
}
