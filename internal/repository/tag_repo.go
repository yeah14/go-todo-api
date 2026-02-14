package repository

import (
	"context"
	"errors"
	"fmt"
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
	if tag == nil {
		return errors.New("标签为空")
	}
	err := t.db.WithContext(ctx).Create(tag).Error
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (t tagRepo) GetByID(ctx context.Context, id uint) (*model.Tag, error) {
	Tag := model.Tag{}
	err := t.db.WithContext(ctx).First(&Tag, id).Error
	if err != nil {
		return nil, errors.New("未找到标签")
	}
	return &Tag, nil
}

func (t tagRepo) Delete(ctx context.Context, id uint) error {
	if err := t.db.WithContext(ctx).Delete(&model.Tag{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (t tagRepo) Update(ctx context.Context, tag *model.Tag) error {
	if tag == nil {
		return errors.New("待办事项为空")
	}
	return t.db.WithContext(ctx).Save(tag).Error
}

func (t tagRepo) GetAll(ctx context.Context, userID uint) ([]model.Tag, error) {
	var tags []model.Tag
	err := t.db.WithContext(ctx).Find(&tags).Where("user_id = ?", userID).Error
	if err != nil {
		return nil, err
	}
	return tags, nil
}
