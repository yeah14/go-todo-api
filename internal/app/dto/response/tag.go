package response

import "go-todo-api/internal/domain/model"

type TagResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

func GenerateTagResponse(tag *model.Tag) (*TagResponse, error) {
	tagResponse := TagResponse{
		ID:    tag.ID,
		Name:  tag.Name,
		Color: tag.Color,
	}
	return &tagResponse, nil
}
