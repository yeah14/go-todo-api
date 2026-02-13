package repository

import (
	"context"
	"errors"
	"go-todo-api/internal/domain/model"

	"gorm.io/gorm"
)

type tagRepo struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) TagRepository {
	return &tagRepo{db: db}
}

func (t tagRepo) Create(ctx context.Context, tag *model.Tag) error {
	//TODO implement me
	panic("implement me")
}

func (t tagRepo) GetByID(ctx context.Context, id uint) (*model.Tag, error) {
	//TODO implement me
	Tag := model.Tag{}
	err := t.db.WithContext(ctx).First(&Tag, id).Error
	if err != nil {
		return nil, errors.New("未找到标签")
	}
	return &Tag, nil
}

func (t tagRepo) Delete(ctx context.Context, id uint) error {
	//TODO implement me
	panic("implement me")
}

func (t tagRepo) Update(ctx context.Context, tag *model.Tag) error {
	//TODO implement me
	panic("implement me")
}
