package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"go-todo-api/internal/domain/model"
	"time"

	"github.com/redis/go-redis/v9"
)

type UserCache struct {
	rdb        *redis.Client
	defaultTTL time.Duration // 默认缓存时间，如30分钟
}

func NewUserCache(rdb *redis.Client) UserCache {
	return UserCache{
		rdb:        rdb,
		defaultTTL: 30 * time.Minute,
	}
}

// Get 获取用户缓存
func (uc *UserCache) Get(ctx context.Context, userID uint) (*model.User, error) {
	key := fmt.Sprintf("user:%d", userID)
	data, err := uc.rdb.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // 缓存未命中，不是错误
		}
		return nil, err // 网络等错误
	}

	var user model.User
	if err := json.Unmarshal(data, &user); err != nil {
		// 反序列化失败，可能是数据格式损坏，删除这个脏数据
		uc.rdb.Del(ctx, key)
		return nil, nil
	}
	return &user, nil
}

// Set 设置用户缓存
func (uc *UserCache) Set(ctx context.Context, user *model.User) error {
	key := fmt.Sprintf("user:%d", user.ID)
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return uc.rdb.Set(ctx, key, data, uc.defaultTTL).Err()
}

// Delete 删除用户缓存
func (uc *UserCache) Delete(ctx context.Context, userID uint) error {
	key := fmt.Sprintf("user:%d", userID)
	return uc.rdb.Del(ctx, key).Err()
}
