package service

import (
	"context"
	"errors"
	"go-todo-api/internal/app/dto/request"
	"go-todo-api/internal/app/dto/response"
	"go-todo-api/internal/domain/model"
	"go-todo-api/internal/repository"
	"go-todo-api/pkg/cache"
	"go-todo-api/pkg/encryption"
	"log"
	"time"
)

type userService struct {
	userRepo  repository.UserRepository
	userCache cache.UserCache
}

func NewUserService(userRepo repository.UserRepository, userCache cache.UserCache) UserService {
	return &userService{userRepo: userRepo, userCache: userCache}
}

func (s *userService) GetProfile(ctx context.Context, userID uint) (*response.UserProfileResponse, error) {
	cachedUser, err := s.userCache.Get(ctx, userID)
	if err != nil {
		// 记录日志，但不中断流程，降级到查数据库
		log.Printf("获取用户缓存失败: %v，降级查DB", err)
	} else if cachedUser != nil {
		// 缓存命中，直接返回
		return s.generateUserResponse(cachedUser)
	}
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	if user.Status == 0 {
		return nil, errors.New("用户已被禁用")
	}
	go func() { // 使用goroutine异步写入，不阻塞本次请求响应
		ctxAsync, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := s.userCache.Set(ctxAsync, user); err != nil {
			log.Printf("异步写入用户缓存失败: %v", err)
		}
	}()
	return s.generateUserResponse(user)
}

func (s *userService) ChangePassword(ctx context.Context, userID uint, req *request.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return errors.New("用户不存在")
	}
	if !encryption.CheckPasswordHash(req.OldPassword, user.PasswordHash) {
		return errors.New("原密码不正确")
	}
	if req.ConfirmPassword != req.NewPassword {
		return errors.New("新密码和确认密码不一致")
	}
	hashPassword, err := encryption.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = hashPassword
	err = s.userRepo.Update(ctx, user)
	go func() {
		ctxAsync, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := s.userCache.Delete(ctxAsync, userID); err != nil {
			log.Printf("删除用户缓存失败: %v", err)
		}
	}()
	if err != nil {
		return err
	}
	return nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID uint, req *request.UpdateProfileRequest) (*response.UserProfileResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	if !req.HasUpdate() {
		return nil, errors.New("未更新信息")
	}
	// 2. 检查是否有重复的用户名
	if req.Username != nil && *req.Username != user.Username {
		existingUser, _ := s.userRepo.GetByUsername(ctx, *req.Username)
		if existingUser != nil && existingUser.ID != userID {
			return nil, errors.New("用户名已存在")
		}
		user.Username = *req.Username
	}

	// 3. 检查是否有重复的邮箱
	if req.Email != nil && *req.Email != user.Email {
		existingUser, _ := s.userRepo.GetByEmail(ctx, *req.Email)
		if existingUser != nil && existingUser.ID != userID {
			return nil, errors.New("邮箱已存在")
		}
		user.Email = *req.Email
	}
	if req.Avatar != nil {
		user.AvatarURL = req.Avatar
	}
	err = s.userRepo.Update(ctx, user)
	go func() {
		ctxAsync, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if err := s.userCache.Delete(ctxAsync, userID); err != nil {
			log.Printf("删除用户缓存失败: %v", err)
		}
	}()
	if err != nil {
		return nil, err
	}
	return s.generateUserResponse(user)
}

func (s *userService) generateUserResponse(user *model.User) (*response.UserProfileResponse, error) {
	userResponse := response.UserProfileResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Avatar:    user.AvatarURL,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	return &userResponse, nil
}
