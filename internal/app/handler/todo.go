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

// CreateTodo 创建待办事项
// @Summary      创建待办事项
// @Description  创建一个新的待办事项
// @Tags         todos
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body request.CreateTodoRequest true "待办事项信息"
// @Success      200  {object}  response.TodoResponse "创建成功"
// @Failure      400  {object}  response.Response   "请求参数错误或标题为空"
// @Failure      401  {object}  response.Response   "用户未认证"
// @Router       /todos [post]
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

// GetTodoById 根据 ID 获取待办事项详情
// @Summary      获取待办事项详情
// @Description  根据待办事项 ID 获取详细信息
// @Tags         todos
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "待办事项 ID"
// @Success      200  {object}  response.TodoResponse "待办事项详情"
// @Failure      400  {object}  response.Response   "无效的ID或请求失败"
// @Failure      401  {object}  response.Response   "用户未认证"
// @Router       /todos/{id} [get]
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

// Delete 删除待办事项
// @Summary      删除待办事项
// @Description  根据 ID 删除待办事项
// @Tags         todos
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "待办事项 ID"
// @Success      200  {object}  map[string]interface{} "删除成功，返回消息和ID"
// @Failure      400  {object}  response.Response   "无效的ID或删除失败"
// @Failure      401  {object}  response.Response   "用户未认证"
// @Router       /todos/{id} [delete]
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

// Update 更新待办事项
// @Summary      更新待办事项
// @Description  根据 ID 更新待办事项信息
// @Tags         todos
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path int                         true "待办事项 ID"
// @Param        request body request.UpdateTodoRequest   true "更新信息"
// @Success      200     {object} response.TodoResponse   "更新成功"
// @Failure      400     {object} response.Response       "请求参数错误或更新失败"
// @Failure      401     {object} response.Response       "用户未认证"
// @Router       /todos/{id} [put]
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

// GetTodos 获取待办事项列表（支持分页和筛选）
// @Summary      获取待办事项列表
// @Description  分页获取待办事项列表，支持按状态、优先级、标签筛选
// @Tags         todos
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page      query int false "页码，默认1"  minimum(1)
// @Param        page_size query int false "每页数量，默认10" minimum(1) maximum(100)
// @Param        status    query int false "状态筛选：0-待办，1-进行中，2-已完成" Enums(0,1,2)
// @Param        priority  query int false "优先级筛选：1-低，2-中，3-高，4-紧急" Enums(1,2,3,4)
// @Param        tag_id    query int false "标签ID筛选"
// @Success      200  {object}  response.TodoListResponse "待办事项列表（包含分页信息）"
// @Failure      400  {object}  response.Response   "请求参数错误"
// @Failure      401  {object}  response.Response   "用户未认证"
// @Router       /todos [get]
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

// BatchUpdateStatus 批量更新待办事项状态
// @Summary      批量更新状态
// @Description  批量更新多个待办事项的状态
// @Tags         todos
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body request.BatchUpdateStatusRequest true "批量更新请求"
// @Success      200  {array} response.TodoResponse    "更新结果"
// @Failure      400  {object}  response.Response              "请求参数错误或更新失败"
// @Failure      401  {object}  response.Response              "用户未认证"
// @Router       /todos/batch/status [put]
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
