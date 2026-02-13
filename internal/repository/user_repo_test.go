package repository

import (
	"context"
	"go-todo-api/internal/domain/model"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB 创建测试数据库
func setupTestDB(t *testing.T) *gorm.DB {
	// 使用 SQLite 内存数据库进行测试
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// 自动迁移表结构
	err = db.AutoMigrate(&model.User{})
	assert.NoError(t, err)

	return db
}

func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	// 创建测试用户
	avatarURL := "avatar.jpg"
	user := &model.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
		AvatarURL:    &avatarURL,
		Status:       1,
	}

	// 测试创建
	err := repo.Create(ctx, user)
	assert.NoError(t, err)
	assert.NotZero(t, user.ID, "创建后应该生成ID")

	// 验证数据
	retrieved, err := repo.GetByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.NotNil(t, retrieved)
	assert.Equal(t, user.Username, retrieved.Username)
	assert.Equal(t, user.Email, retrieved.Email)
}

func TestUserRepository_GetByUsername(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	// 先创建用户
	user := &model.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hash",
	}
	repo.Create(ctx, user)

	// 测试通过用户名查询
	found, err := repo.GetByUsername(ctx, "testuser")
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, "testuser", found.Username)

	// 测试不存在的用户
	notFound, err := repo.GetByUsername(ctx, "nonexistent")
	assert.Error(t, err)
	assert.Nil(t, notFound)
}

func TestUserRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	// 创建用户
	user := &model.User{
		Username:     "original",
		Email:        "original@example.com",
		PasswordHash: "hash",
	}
	repo.Create(ctx, user)

	// 更新用户
	newEmail := "updated@example.com"
	user.Email = newEmail
	err := repo.Update(ctx, user)
	assert.NoError(t, err)

	// 验证更新
	updated, err := repo.GetByID(ctx, user.ID)
	assert.NoError(t, err)
	assert.Equal(t, newEmail, updated.Email)
}

func TestUserRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	// 创建用户
	user := &model.User{
		Username:     "todelete",
		Email:        "delete@example.com",
		PasswordHash: "hash",
	}
	repo.Create(ctx, user)

	// 删除用户
	err := repo.Delete(ctx, user.ID)
	assert.NoError(t, err)

	// 验证已删除
	deleted, err := repo.GetByID(ctx, user.ID)
	assert.Error(t, err)
	assert.Nil(t, deleted)
}
