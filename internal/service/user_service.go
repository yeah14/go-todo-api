package service

import (
	"context"
	"errors"
	"go-todo-api/internal/app/dto/request"
	"go-todo-api/internal/app/dto/response"
	"go-todo-api/internal/domain/model"
	"go-todo-api/internal/repository"
	"go-todo-api/pkg/encryption"
)

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetProfile(ctx context.Context, userID uint) (*response.UserProfileResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	if user.Status == 0 {
		return nil, errors.New("用户已被禁用")
	}
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
