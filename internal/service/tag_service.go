package service

import (
	"context"
	"errors"
	"go-todo-api/internal/app/dto/request"
	"go-todo-api/internal/app/dto/response"
	"go-todo-api/internal/domain/model"
	"go-todo-api/internal/repository"
)

type tagService struct {
	tagRepo repository.TagRepository
}

func NewTagService(tagRepo repository.TagRepository) TagService {
	return &tagService{tagRepo: tagRepo}
}

func (t tagService) Create(ctx context.Context, userID uint, req *request.CreateTagRequest) (*response.TagResponse, error) {
	tags := &model.Tag{
		Name:   req.Name,
		Color:  req.Color,
		UserID: userID,
	}
	err := t.tagRepo.Create(ctx, tags)
	if err != nil {
		return nil, errors.New("创建标签失败")
	}
	return response.GenerateTagResponse(tags)
}

func (t tagService) GetTags(ctx context.Context, userID uint) (*response.TagListResponse, error) {
	tags, err := t.tagRepo.GetAll(ctx, userID)
	if err != nil {
		return nil, errors.New("查找标签列表失败")
	}
	return response.GenerateTagListResponse(tags)
}

func (t tagService) Update(ctx context.Context, userID uint, req *request.UpdateTagRequest) (*response.TagResponse, error) {
	tag, err := t.tagRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, errors.New("查找标签失败")
	}
	if tag.UserID != userID {
		return nil, errors.New("此标签不属于该用户")
	}
	if req.Name != nil {
		tag.Name = *req.Name
	}
	if req.Color != nil {
		tag.Color = *req.Color
	}
	return response.GenerateTagResponse(tag)
}

func (t tagService) Delete(ctx context.Context, userID uint, req *request.DeleteTagRequest) error {
	tag, err := t.tagRepo.GetByID(ctx, req.ID)
	if err != nil {
		return errors.New("查找标签失败")
	}
	if tag.UserID != userID {
		return errors.New("标签不属于该用户")
	}
	err = t.tagRepo.Delete(ctx, tag.ID)
	if err != nil {
		return errors.New("删除标签失败")
	}
	return nil
}
