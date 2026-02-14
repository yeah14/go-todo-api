package response

import (
	"go-todo-api/internal/domain/model"
	"time"
)

type TagResponse struct {
	ID       uint       `json:"id"`
	Name     string     `json:"name"`
	Color    string     `json:"color"`
	CreateAt *time.Time `json:"create_at"`
	UpdateAt *time.Time `json:"update_at"`
}

type TagListResponse struct {
	Tags []TagResponse `json:"tags"`
}

func GenerateTagResponse(tag *model.Tag) (*TagResponse, error) {
	tagResponse := TagResponse{
		ID:       tag.ID,
		Name:     tag.Name,
		Color:    tag.Color,
		CreateAt: &tag.CreatedAt,
		UpdateAt: &tag.UpdatedAt,
	}
	return &tagResponse, nil
}

func GenerateTagListResponse(tags []model.Tag) (*TagListResponse, error) {
	tagListResponse := TagListResponse{
		Tags: make([]TagResponse, 0),
	}
	for _, tag := range tags {
		tagResponse, _ := GenerateTagResponse(&tag)
		tagListResponse.Tags = append(tagListResponse.Tags, *tagResponse)
	}
	return &tagListResponse, nil
}
