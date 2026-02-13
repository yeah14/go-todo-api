package handler

import (
	"go-todo-api/internal/app/dto/request"
	"go-todo-api/internal/service"
	"go-todo-api/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type todoHandler struct {
	todoSevice service.TodoService
}

func NewTodoHandler(todoService service.TodoService) *todoHandler {
	return &todoHandler{todoSevice: todoService}
}

func (h *todoHandler) CreateTodo(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		response.Unauthorized(c, "用户未认证")
		return
	}
	req := new(request.CreateTodoRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if req.Title == "" {
		response.BadRequest(c, "标题不能为空")
		return
	}
	todoResp, err := h.todoSevice.CreateTodo(c.Request.Context(), userID.(uint), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, todoResp)
}

func (h *todoHandler) GetTodoById(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		response.BadRequest(c, "用户未认证")
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	req := &request.GetTodoBYIdRequest{
		ID: uint(id),
	}
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	todo, err := h.todoSevice.GetTodoByID(c.Request.Context(), userID.(uint), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, todo)
}

func (h *todoHandler) Delete(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		response.BadRequest(c, "用户未认证")
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, err.Error())
	}
	req := &request.GetTodoBYIdRequest{ID: uint(id)}
	err = h.todoSevice.Delete(c.Request.Context(), userID.(uint), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, gin.H{
		"message": "删除事项成功",
		"user_id": userID,
		"id":      id,
	})
}

func (h *todoHandler) Update(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		response.BadRequest(c, "用户未认证")
		return
	}
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	req := new(request.UpdateTodoRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	req.ID = uint(id)
	todoResp, err := h.todoSevice.Update(c.Request.Context(), userID.(uint), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, todoResp)
}

func (h *todoHandler) GetTodos(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		response.BadRequest(c, "用户未认证")
		return
	}
	var req request.GetTodoRequest
	err := c.ShouldBindQuery(&req)
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	todos, err := h.todoSevice.GetTodos(c.Request.Context(), userID.(uint), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, todos)
}

func (h *todoHandler) BatchUpdateStatus(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		response.BadRequest(c, "用户未认证")
		return
	}
	req := new(request.BatchUpdateStatusRequest)
	if err := c.ShouldBindJSON(req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	todosResp, err := h.todoSevice.BatchUpdateStatus(c.Request.Context(), userID.(uint), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, todosResp)
}
