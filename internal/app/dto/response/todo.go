package response

import (
	"go-todo-api/internal/domain/model"
	"time"
)

type TodoResponse struct {
	ID           uint          `json:"id"`
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	Status       uint8         `json:"status"`
	StatusText   string        `json:"status_text"`
	Priority     uint8         `json:"priority"`
	PriorityText string        `json:"priority_text"`
	DueDate      string        `json:"due_date"`
	CompletedAt  *time.Time    `json:"completed_at"`
	UpdatedAt    *time.Time    `json:"updated_at"`
	Tags         []TagResponse `json:"tags"`
}

type Pagination struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int64 `json:"total_pages"`
}

type TodoListResponse struct {
	Todos      []TodoResponse `json:"todos"`
	Pagination Pagination     `json:"pagination"`
}

func GenerateTodoResponse(todo *model.Todo) (*TodoResponse, error) {
	todoResponse := ToTodoResponse(*todo)
	return &todoResponse, nil
}

func getStatusText(status uint8) string {
	switch status {
	case 0:
		return "待办"
	case 1:
		return "进行中"
	case 2:
		return "已完成"
	default:
		return ""
	}
}

func getPriorityText(status uint8) string {
	switch status {
	case 1:
		return "低"
	case 2:
		return "中"
	case 3:
		return "高"
	case 4:
		return "紧急"
	default:
		return ""
	}
}

func ToTodoResponse(todo model.Todo) TodoResponse {
	var dueDate string
	if todo.DueDate == nil {
		dueDate = ""
	} else {
		dueDate = todo.DueDate.Format("2006-01-02T15:04:05")
	}
	tags := make([]TagResponse, 0)
	for _, tag := range todo.Tags {
		tagResponse, _ := GenerateTagResponse(&tag)
		tags = append(tags, *tagResponse)
	}
	return TodoResponse{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: *todo.Description,
		Status:      todo.Status,
		Priority:    todo.Priority,
		DueDate:     dueDate,
		CompletedAt: todo.CompletedAt,
		UpdatedAt:   &todo.UpdatedAt,
		Tags:        tags,
	}
}

func ToTodosResponses(todos []model.Todo) []TodoResponse {
	responses := make([]TodoResponse, len(todos))
	for i, todo := range todos {
		responses[i] = ToTodoResponse(todo)
	}
	return responses
}

func ToTodoListResponse(todos []model.Todo, total int64, page, pageSize int) *TodoListResponse {
	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)
	return &TodoListResponse{
		Todos: ToTodosResponses(todos),
		Pagination: Pagination{
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
			Total:      total,
		},
	}
}
