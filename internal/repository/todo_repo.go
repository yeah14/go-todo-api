package repository

import (
	"context"
	"errors"
	"fmt"
	"go-todo-api/internal/domain/model"
	"strings"
	"time"

	"gorm.io/gorm"
)

type todoRepo struct {
	db *gorm.DB
}

func NewTodoRepo(db *gorm.DB) TodoRepository {
	return &todoRepo{db: db}
}

func (t *todoRepo) Create(ctx context.Context, todo *model.Todo) error {
	if todo == nil {
		return errors.New("用户为空")
	}
	return t.db.WithContext(ctx).Create(todo).Error

}

func (t *todoRepo) GetByID(ctx context.Context, id uint) (*model.Todo, error) {
	todo := new(model.Todo)
	err := t.db.WithContext(ctx).First(todo, "id=?", id).Error
	if err != nil {
		return nil, err
	}
	return todo, err
}

func (t *todoRepo) Update(ctx context.Context, todo *model.Todo) error {
	if todo == nil {
		return errors.New("用户为空")
	}
	return t.db.WithContext(ctx).Save(todo).Error
}

func (t *todoRepo) Delete(ctx context.Context, id uint) error {
	todo := new(model.Todo)
	err := t.db.WithContext(ctx).First(todo, "id=?", id).Error
	if err != nil {
		return err
	}
	return t.db.WithContext(ctx).Delete(todo).Error
}

func (t *todoRepo) ListByUser(ctx context.Context, userID uint, query *TodoQuery) ([]model.Todo, int64, error) {
	var total int64
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 {
		query.PageSize = 10
	}
	offset := (query.Page - 1) * query.PageSize
	db := t.db.WithContext(ctx).Model(&model.Todo{}).Where("user_id=?", userID)
	db = t.applyQueryConditions(db, query)
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("统计待办事项失败: %w", err)
	}
	db = db.Offset(offset).Limit(query.PageSize)
	var todos []model.Todo
	if err := db.Find(&todos).Error; err != nil {
		return nil, 0, fmt.Errorf("查找待办事项失败: %w", err)
	}

	return todos, total, nil
}

func (t *todoRepo) UpdateStatus(ctx context.Context, id uint, status uint) error {
	err := t.db.WithContext(ctx).Model(&model.User{}).Where("id=?", id).Update("status", status).Error
	if err != nil {
		return err
	}
	return nil
}

func (t *todoRepo) applyQueryConditions(db *gorm.DB, query *TodoQuery) *gorm.DB {
	// 状态筛选
	if query.Status != nil {
		//fmt.Println(*query.Status)
		db = db.Where("status = ?", *query.Status)
	}

	// 优先级筛选
	if query.Priority != nil {
		//fmt.Println(*query.Priority)
		db = db.Where("priority = ?", *query.Priority)
	}

	// 关键词搜索（标题和描述）
	if query.Keyword != nil {
		keyword := "%" + strings.ToLower(*query.Keyword) + "%"
		db = db.Where("(LOWER(title) LIKE ? OR LOWER(description) LIKE ?)", keyword, keyword)
	}

	// 截止日期筛选
	if query.DueDate != nil {
		startOfDay := time.Date(query.DueDate.Year(), query.DueDate.Month(), query.DueDate.Day(), 0, 0, 0, 0, query.DueDate.Location())
		endOfDay := startOfDay.Add(24 * time.Hour)
		db = db.Where("due_date >= ? AND due_date < ?", startOfDay, endOfDay)
	}

	// 标签筛选
	if len(query.TagIDs) > 0 {
		db = db.Joins("JOIN todo_tags ON todo_tags.todo_id = todos.id").
			Where("todo_tags.tag_id IN ?", query.TagIDs).
			Group("todos.id")
	}

	return db
}
