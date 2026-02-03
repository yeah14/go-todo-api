package repository

import (
	"context"
	"errors"
	"go-todo-api/internal/domain/model"

	"gorm.io/gorm"
)

type userrep struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userrep{db: db}
}

func (u userrep) Create(ctx context.Context, user *model.User) error {
	if user == nil {
		return errors.New("用户为空")
	}
	return u.db.WithContext(ctx).Create(user).Error
}

func (u userrep) GetByID(ctx context.Context, id uint) (*model.User, error) {
	user := new(model.User)
	err := u.db.WithContext(ctx).First(user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u userrep) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	user := new(model.User)
	err := u.db.WithContext(ctx).First(user, "username = ?", username).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u userrep) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	user := new(model.User)
	err := u.db.WithContext(ctx).First(user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u userrep) Update(ctx context.Context, user *model.User) error {
	if user == nil {
		return errors.New("用户为空")
	}
	return u.db.WithContext(ctx).Save(user).Error
}

func (u userrep) Delete(ctx context.Context, id uint) error {
	user := new(model.User)
	err := u.db.WithContext(ctx).First(user, "id = ?", id).Error
	if err != nil {
		return err
	}
	return u.db.WithContext(ctx).Delete(user).Error
}

func (u userrep) List(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize
	err := u.db.WithContext(ctx).Model(&model.User{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	err = u.db.WithContext(ctx).Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}
	return users, total, nil
}

func (u userrep) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	var count int64
	err := u.db.WithContext(ctx).Model(&model.User{}).Where("username = ?", username).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (u userrep) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var count int64
	err := u.db.WithContext(ctx).Model(&model.User{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
