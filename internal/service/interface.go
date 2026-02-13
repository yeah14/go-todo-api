package service

import (
	"context"
	"go-todo-api/internal/app/dto/request"
	"go-todo-api/internal/app/dto/response"
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
