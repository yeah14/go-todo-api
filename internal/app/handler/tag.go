package handler

import (
	"go-todo-api/internal/app/dto/request"
	"go-todo-api/internal/service"
	"go-todo-api/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	tagService service.TagService
}

func NewTagHandler(tagService service.TagService) *TagHandler {
	return &TagHandler{tagService: tagService}
}

// Create 创建标签
// @Summary      创建新标签
// @Description  为当前用户创建一个新的标签
// @Tags         tags
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body request.CreateTagRequest true "标签信息"
// @Success      200  {object}  response.TagResponse "创建成功"
// @Failure      400  {object}  response.Response   "请求参数错误或用户未认证"
func (h *TagHandler) Create(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.BadRequest(c, "用户未认证")
		return
	}
	var req request.CreateTagRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	tagResp, err := h.tagService.Create(c.Request.Context(), userID.(uint), &req)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.Success(c, tagResp)
}

// GetTags 获取标签列表
// @Summary      获取当前用户的标签列表
// @Description  返回当前用户创建的所有标签
// @Tags         tags
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} response.TagListResponse "标签列表"
// @Failure      400  {object}  response.Response   "用户未认证"
// @Failure      500  {object}  response.Response   "服务器内部错误"
// @Router       /tags [get]
func (h *TagHandler) GetTags(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.BadRequest(c, "用户未认证")
		return
	}
	tagsResp, err := h.tagService.GetTags(c.Request.Context(), userID.(uint))
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.Success(c, tagsResp)
}

// Update 更新标签
// @Summary      更新标签
// @Description  根据标签ID更新标签信息（名称、颜色等）
// @Tags         tags
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path int                         true "标签ID"
// @Param        request body request.UpdateTagRequest    true "更新信息"
// @Success      200     {object} response.TagResponse    "更新成功"
// @Failure      400     {object} response.Response       "请求参数错误或用户未认证"
// @Failure      500     {object} response.Response       "服务器内部错误"
// @Router       /tags/{id} [put]
func (h *TagHandler) Update(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.BadRequest(c, "用户未认证")
		return
	}
	var req request.UpdateTagRequest
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	req.ID = uint(id)
	err = c.ShouldBindJSON(&req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	tagResp, err := h.tagService.Update(c.Request.Context(), userID.(uint), &req)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.Success(c, tagResp)
}

// Delete 删除标签
// @Summary      删除标签
// @Description  根据标签ID删除标签（会解除与待办事项的关联）
// @Tags         tags
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "标签ID"
// @Success      200 {object} response.Response "删除成功（数据为null）"
// @Failure      400 {object} response.Response "请求参数错误或用户未认证"
// @Failure      500 {object} response.Response "服务器内部错误"
// @Router       /tags/{id} [delete]
func (h *TagHandler) Delete(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.BadRequest(c, "用户未认证")
		return
	}
	var req request.DeleteTagRequest
	id, err := strconv.ParseInt(c.Param("id"), 10, 32)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	req.ID = uint(id)
	err = h.tagService.Delete(c.Request.Context(), userID.(uint), &req)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}
	response.Success(c, nil)
}
