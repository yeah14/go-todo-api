package repository

import (
	"context"
	"go-todo-api/internal/domain/model"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByID(ctx context.Context, id uint) (*model.User, error)
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

type TodoRepository interface {
	Create(ctx context.Context, todo *model.Todo) error
	GetByID(ctx context.Context, id uint) (*model.Todo, error)
	Update(ctx context.Context, todo *model.Todo) error
	Delete(ctx context.Context, id uint) error
	ListByUser(ctx context.Context, userID uint, query *TodoQuery) ([]model.Todo, int64, error)
	UpdateStatus(ctx context.Context, id uint, status uint) error
	//ListWithTafs(ctx context.Context, userID uint, query *TodoQuery) ([]model.Todo, int64, error)
}

type TagRepository interface {
	Create(ctx context.Context, tag *model.Tag) error
	GetByID(ctx context.Context, id uint) (*model.Tag, error)
	Delete(ctx context.Context, id uint) error
	Update(ctx context.Context, tag *model.Tag) error
	GetAll(ctx context.Context, userID uint) ([]model.Tag, error)
}

type TodoQuery struct {
	Page     int
	PageSize int
	Status   *uint8
	Priority *uint8
	Keyword  *string
	TagIDs   []uint
	DueDate  *time.Time
	SortBy   string
	Order    string // "asc" or "desc"
}
