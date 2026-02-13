package request

import (
	"go-todo-api/internal/repository"
)

type GetTodoBYIdRequest struct {
	ID uint `json:"id"`
}

type GetTodoRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Status   *uint8 `form:"status" binding:"omitempty,oneof=0 1 2"`
	Priority *uint8 `form:"priority" binding:"omitempty,oneof=1 2 3 4"`
	TagIds   []uint `form:"tag_ids" binding:"omitempty,dive,min=1"`
}

type CreateTodoRequest struct {
	Title       string  `json:"title" binding:"required,min=1,max=255"`
	Description *string `json:"description" binding:"omitempty,min=1,max=255"`
	Status      *uint8  `json:"status" binding:"required,oneof=0 1 2"`
	Priority    *uint8  `json:"priority" binding:"required,oneof=1 2 3 4"`
	DueDate     *string `json:"due_date" binding:"omitempty,min=1,max=255"`
	TagIDs      []uint  `json:"tag_ids" binding:"omitempty"`
}
type UpdateTodoRequest struct {
	ID          uint    `json:"id"`
	Title       *string `json:"title" binding:"omitempty,min=1,max=255"`
	Description *string `json:"description" binding:"omitempty,min=1,max=255"`
	Status      *uint8  `json:"status" binding:"omitempty,oneof=0 1 2"`
	Priority    *uint8  `json:"priority" binding:"omitempty,oneof=1 2 3 4"`
}

type BatchUpdateStatusRequest struct {
	Ids    []uint `json:"todo_ids" binding:"required,dive,min=1"`
	Status uint8  `json:"status" binding:"required,oneof=0 1 2"`
}

func (req *GetTodoRequest) ToTodoQuery() *repository.TodoQuery {
	//fmt.Println(req.Status, req.Priority)
	return &repository.TodoQuery{
		Page:     req.Page,
		PageSize: req.PageSize,
		Status:   req.Status,
		Priority: req.Priority,
		TagIDs:   req.TagIds,
	}

}
