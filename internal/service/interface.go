package service

import (
	"context"
	"go-todo-api/internal/app/dto/request"
	"go-todo-api/internal/app/dto/response"
	"time"
)

type AuthService interface {
	Register(ctx context.Context, req *request.RegisterRequest) (*response.AuthResponse, error)
	Login(ctx context.Context, req *request.LoginRequest) (*response.AuthResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*response.AuthResponse, error)
}

type UserService interface {
	ChangePassword(ctx context.Context, userID uint, req *request.ChangePasswordRequest) error
	UpdateProfile(ctx context.Context, userID uint, req *request.UpdateProfileRequest) (*response.UserProfileResponse, error)
	GetProfile(ctx context.Context, userID uint) (*response.UserProfileResponse, error)
}

type TodoService interface {
	GetTodoByID(ctx context.Context, userID uint, req *request.GetTodoBYIdRequest) (*response.TodoResponse, error)
	CreateTodo(ctx context.Context, userID uint, req *request.CreateTodoRequest) (*response.TodoResponse, error)
	Delete(ctx context.Context, userID uint, req *request.GetTodoBYIdRequest) error
	Update(ctx context.Context, userID uint, req *request.UpdateTodoRequest) (*response.TodoResponse, error)
	GetTodos(ctx context.Context, userID uint, req *request.GetTodoRequest) (*response.TodoListResponse, error)
	BatchUpdateStatus(ctx context.Context, userID uint, req *request.BatchUpdateStatusRequest) ([]response.TodoResponse, error)
}

type TagService interface {
	Create(ctx context.Context, userID uint, req *request.CreateTagRequest) (*response.TagResponse, error)
	GetTags(ctx context.Context, userID uint) (*response.TagListResponse, error)
	Update(ctx context.Context, userID uint, req *request.UpdateTagRequest) (*response.TagResponse, error)
	Delete(ctx context.Context, userID uint, req *request.DeleteTagRequest) error
}

type BlacklistService interface {
	AddtoBlacklist(ctx context.Context, tokenString string, expiresIn time.Duration) error
	IsInBlacklist(ctx context.Context, tokenString string) (bool, error)
}
