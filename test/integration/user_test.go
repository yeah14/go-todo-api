package integration

import (
	"context"
	"go-todo-api/config"
	"go-todo-api/internal/domain/model"
	"go-todo-api/internal/repository"
	"go-todo-api/pkg/database"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type UserIntegrationTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo repository.UserRepository
	ctx  context.Context
}

func (suite *UserIntegrationTestSuite) SetupSuite() {
	// 初始化配置
	config.InitConfig("../../config/dev.yaml")

	// 连接数据库
	var err error
	suite.db, err = database.ConnectMysql(config.GlobalConfig.Database)
	assert.NoError(suite.T(), err)

	// 创建仓库
	suite.repo = repository.NewUserRepository(suite.db)
	suite.ctx = context.Background()

	// 清理测试数据
	suite.db.Exec("DELETE FROM users WHERE username LIKE 'test_%'")
}

func (suite *UserIntegrationTestSuite) TearDownSuite() {
	if suite.db != nil {
		sqlDB, _ := suite.db.DB()
		sqlDB.Close()
	}
}

func (suite *UserIntegrationTestSuite) TestFullUserCRUD() {
	// 创建用户
	user := &model.User{
		Username:     "test_john",
		Email:        "test_john@example.com",
		PasswordHash: "$2a$10$hashedpassword",
	}

	err := suite.repo.Create(suite.ctx, user)
	suite.NoError(err)
	suite.NotZero(user.ID)

	// 通过ID查询
	retrieved, err := suite.repo.GetByID(suite.ctx, user.ID)
	suite.NoError(err)
	suite.NotNil(retrieved)
	suite.Equal(user.Username, retrieved.Username)

	// 通过用户名查询
	byUsername, err := suite.repo.GetByUsername(suite.ctx, "test_john")
	suite.NoError(err)
	suite.NotNil(byUsername)
	suite.Equal(user.Email, byUsername.Email)

	// 通过邮箱查询
	byEmail, err := suite.repo.GetByEmail(suite.ctx, "test_john@example.com")
	suite.NoError(err)
	suite.NotNil(byEmail)
	suite.Equal(user.Username, byUsername.Username)

	// 更新用户
	newAvatar := "https://new-avatar.com/pic.jpg"
	retrieved.AvatarURL = &newAvatar
	err = suite.repo.Update(suite.ctx, retrieved)
	suite.NoError(err)

	// 验证更新
	updated, err := suite.repo.GetByID(suite.ctx, user.ID)
	suite.NoError(err)
	suite.NotNil(updated.AvatarURL)
	suite.Equal(newAvatar, *updated.AvatarURL)

	// 删除用户
	err = suite.repo.Delete(suite.ctx, user.ID)
	suite.NoError(err)

	// 验证删除
	deleted, err := suite.repo.GetByID(suite.ctx, user.ID)
	suite.Error(err)
	suite.Nil(deleted)
}

func TestUserIntegrationTestSuite(t *testing.T) {
	// 只在有数据库连接时才运行集成测试
	if testing.Short() {
		t.Skip("跳过集成测试")
	}
	suite.Run(t, new(UserIntegrationTestSuite))
}
