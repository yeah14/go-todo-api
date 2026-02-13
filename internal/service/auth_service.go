package service

import (
	"context"
	"errors"
	"go-todo-api/internal/app/dto/request"
	"go-todo-api/internal/app/dto/response"
	"go-todo-api/internal/domain/model"
	"go-todo-api/internal/repository"
	"go-todo-api/pkg/encryption"
	"go-todo-api/pkg/jwt"
	"time"
)

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Register(ctx context.Context, req *request.RegisterRequest) (*response.AuthResponse, error) {
	existingUser, _ := s.userRepo.GetByUsername(ctx, req.Username)
	if existingUser != nil {
		return nil, errors.New("用户名已存在")
	}
	existingUser, _ = s.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("邮箱已存在")
	}
	hashPassword, err := encryption.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	user := &model.User{
		ID:           0,
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashPassword,
		Status:       1,
	}
	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	return s.generateAuthResponse(user)
}

func (s *authService) Login(ctx context.Context, req *request.LoginRequest) (*response.AuthResponse, error) {
	user, err := s.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	if !encryption.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("用户名或密码错误")
	}

	if user.Status == 0 {
		return nil, errors.New("用户已被禁用")
	}

	return s.generateAuthResponse(user)
}

func (s authService) RefreshToken(ctx context.Context, refreshToken string) (*response.AuthResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *authService) generateAuthResponse(user *model.User) (*response.AuthResponse, error) {
	// 生成访问令牌
	accessToken, err := jwt.GenerateToken(user.ID, user.Username, false)
	if err != nil {
		return nil, err
	}

	// 生成刷新令牌
	refreshToken, err := jwt.GenerateToken(user.ID, user.Username, true)
	if err != nil {
		return nil, err
	}

	return &response.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(time.Hour * 24 * 7).Unix(), // 7天有效期
		User: response.UserResponse{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Avatar:   user.AvatarURL,
		},
	}, nil
}
