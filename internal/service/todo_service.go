package service

import (
	"context"
	"errors"
	"fmt"
	"go-todo-api/internal/app/dto/request"
	"go-todo-api/internal/app/dto/response"
	"go-todo-api/internal/domain/model"
	"go-todo-api/internal/repository"
	"time"
)

type todoService struct {
	todoRepo repository.TodoRepository
	tagRepo  repository.TagRepository
}

func NewTodoService(todoRepo repository.TodoRepository, tagRepo repository.TagRepository) TodoService {
	return &todoService{
		todoRepo: todoRepo,
		tagRepo:  tagRepo,
	}
}

func (t *todoService) GetTodoByID(ctx context.Context, userID uint, req *request.GetTodoBYIdRequest) (*response.TodoResponse, error) {
	//TODO implement me
	todo, err := t.todoRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, errors.New("查找待办事项失败")
	}
	return response.GenerateTodoResponse(todo)
}

func (t *todoService) CreateTodo(ctx context.Context, userID uint, req *request.CreateTodoRequest) (*response.TodoResponse, error) {

	tags, err := t.validateAndGetTags(ctx, userID, req.TagIDs)
	if err != nil {
		return nil, err
	}
	var dueDate *time.Time
	if *req.DueDate != "" {
		parsedDueDate, err := time.Parse("2006-01-02T15:04:05", *req.DueDate)
		if err != nil {
			return nil, fmt.Errorf("无效的截止日期格式: %w", err)
		}
		parsedDueDate = parsedDueDate.UTC()
		dueDate = &parsedDueDate
	}
	todo := &model.Todo{
		ID:          0,
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Status:      *req.Status,
		DueDate:     dueDate,
		Priority:    *req.Priority,
		Tags:        tags,
	}
	err = t.todoRepo.Create(ctx, todo)
	if err != nil {
		return nil, fmt.Errorf("创建待办事项失败: %w", err)
	}
	return response.GenerateTodoResponse(todo)
}

func (t *todoService) validateAndGetTags(ctx context.Context, userID uint, tagIDs []uint) ([]model.Tag, error) {
	uniqueTags := make(map[uint]bool)
	for _, tagID := range tagIDs {
		uniqueTags[tagID] = true
	}

	tags := make([]model.Tag, len(tagIDs))
	for tagID, _ := range uniqueTags {
		tag, err := t.tagRepo.GetByID(ctx, tagID)
		if err != nil {
			return nil, fmt.Errorf("标签ID %d 不存在", tagID)
		}
		if tag.UserID != userID {
			return nil, fmt.Errorf("标签ID %d 不属于当前用户", tagID)
		}
		tags = append(tags, *tag)
	}
	return tags, nil
}

func (t *todoService) Delete(ctx context.Context, userID uint, req *request.GetTodoBYIdRequest) error {
	todo, err := t.todoRepo.GetByID(ctx, req.ID)
	if err != nil {
		return fmt.Errorf("无法找到待办事项：%v", err.Error())
	}
	if todo.UserID != userID {
		return errors.New("待办事项不属于此用户")
	}
	err = t.todoRepo.Delete(ctx, req.ID)
	if err != nil {
		return fmt.Errorf("删除待办事项失败 %v", err.Error())
	}
	return nil
}

func (t *todoService) Update(ctx context.Context, userID uint, req *request.UpdateTodoRequest) (*response.TodoResponse, error) {
	todo, err := t.todoRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, errors.New("未找到待办事项")
	}
	if todo.UserID != userID {
		return nil, errors.New("待办事项不属于此用户")
	}
	if req.Title != nil {
		todo.Title = *req.Title
	}
	if req.Description != nil {
		todo.Description = req.Description
	}
	if req.Status != nil {
		todo.Status = *req.Status
	}
	if req.Priority != nil {
		todo.Priority = *req.Priority
	}
	err = t.todoRepo.Update(ctx, todo)
	if err != nil {
		return nil, fmt.Errorf("更新待办事项失败%v", err.Error())
	}
	return response.GenerateTodoResponse(todo)
}

func (t *todoService) GetTodos(ctx context.Context, userID uint, req *request.GetTodoRequest) (*response.TodoListResponse, error) {
	todos, total, err := t.todoRepo.ListByUser(ctx, userID, req.ToTodoQuery())
	fmt.Println(todos)
	if err != nil {
		return nil, fmt.Errorf("查询待办事项失败%v", err.Error())
	}
	todosResp := response.ToTodoListResponse(todos, total, req.Page, req.PageSize)
	return todosResp, nil
}

func (t *todoService) BatchUpdateStatus(ctx context.Context, userID uint, req *request.BatchUpdateStatusRequest) ([]response.TodoResponse, error) {
	uniqueIds := make(map[uint]bool)
	for _, id := range req.Ids {
		uniqueIds[id] = true
	}
	todos := make([]model.Todo, 0)
	for ids, _ := range uniqueIds {
		todo, err := t.todoRepo.GetByID(ctx, ids)
		if err != nil {
			return nil, fmt.Errorf("查找%v号待办事项失败：%v", ids, err.Error())
		}
		if todo.UserID != userID {
			return nil, fmt.Errorf("%d号待办事项%s不属于此用户", todo.UserID, todo.Title)
		}
		todo.Status = req.Status
		todos = append(todos, *todo)
	}
	for _, todo := range todos {
		fmt.Println(todo)
		err := t.todoRepo.Update(ctx, &todo)
		if err != nil {
			return nil, err
		}
	}
	return response.ToTodosResponses(todos), nil
}
