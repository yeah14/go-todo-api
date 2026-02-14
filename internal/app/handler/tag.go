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
