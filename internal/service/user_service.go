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

func (s *userService) ChangePassword(ctx context.Context, userID uint, req *request.ChangePasswordRequest) (*response.UserResponse, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	if !encryption.CheckPasswordHash(req.OldPassword, user.PasswordHash) {
		return nil, errors.New("原密码不正确")
	}
	hashPassword, err := encryption.HashPassword(req.NewPassword)
	if err != nil {
		return nil, err
	}
	user.PasswordHash = hashPassword
	err = s.userRepo.Update(ctx, user)
	return s.generateUserResponse(user)
}

func (s *userService) generateUserResponse(user *model.User) (*response.UserResponse, error) {
	userResponse := response.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Avatar:   user.AvatarURL,
	}
	return &userResponse, nil
}
