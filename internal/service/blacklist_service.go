package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type blacklistService struct {
	rdb *redis.Client
}

func NewBlacklistService(rdb *redis.Client) BlacklistService {
	return &blacklistService{rdb: rdb}
}

func getSignature(tokenString string) string {
	//parts := []string{tokenString}
	parts := strings.SplitN(tokenString, ".", 3)
	if len(parts) != 3 {
		return ""
	}
	return parts[2]
}

func (b blacklistService) AddtoBlacklist(ctx context.Context, tokenString string, expiresIn time.Duration) error {
	if expiresIn <= 0 {
		return nil
	}
	Signature := getSignature(tokenString)
	if Signature == "" {
		return errors.New("invalid token format")
	}
	key := fmt.Sprintf("blacklist.%s", Signature)
	err := b.rdb.Set(ctx, key, 1, expiresIn).Err()
	if err != nil {
		return err
	}
	return nil
}

func (b blacklistService) IsInBlacklist(ctx context.Context, tokenString string) (bool, error) {
	fmt.Println(tokenString)
	Signature := getSignature(tokenString)
	if Signature == "" {
		return false, errors.New("invalid token format")
	}
	key := fmt.Sprintf("blacklist.%s", Signature)
	exists, err := b.rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return exists > 0, nil
}
