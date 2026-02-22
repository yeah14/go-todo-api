package handler

import (
	"go-todo-api/internal/app/dto/request"
	"go-todo-api/internal/service"
	"go-todo-api/pkg/response"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *userHandler {
	return &userHandler{userService}
}

// GetProfile 获取当前用户信息
// @Summary      获取个人资料
// @Description  获取当前登录用户的详细信息
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.UserProfileResponse    "用户信息"
// @Failure      401  {object}  response.Response      "未认证"
// @Failure      403  {object}  response.Response      "用户已被禁用"
// @Failure      404  {object}  response.Response      "用户不存在"
// @Failure      500  {object}  response.Response      "服务器内部错误"
// @Router       /user/me [get]
func (h *userHandler) GetProfile(c *gin.Context) {
	userID, exist := c.Get("userID")
	if !exist {
		response.Unauthorized(c, "用户未认证")
		return
	}
	profile, err := h.userService.GetProfile(c.Request.Context(), userID.(uint))
	if err != nil {
		switch err.Error() {
		case "用户不存在":
			response.NotFound(c, err.Error())
		case "用户已被禁用":
			response.Forbidden(c, err.Error())
		default:
			response.InternalServerError(c, err.Error())
		}
		return
	}
	response.Success(c, profile)
}

func (h *userHandler) UpdateProfile(c *gin.Context) {
	var req *request.UpdateProfileRequest
	userID, exist := c.Get("userID")
	if !exist {
		response.Unauthorized(c, "用户未认证")
		return
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if !req.HasUpdate() {
		response.BadRequest(c, "未提供更新信息")
		return
	}

	// 4. 调用服务层
	updatedUser, err := h.userService.UpdateProfile(c.Request.Context(), userID.(uint), req)
	if err != nil {
		switch err.Error() {
		case "用户不存在":
			response.NotFound(c, err.Error())
		case "邮箱已存在", "用户名已存在":
			response.BadRequest(c, err.Error())
		default:
			response.InternalServerError(c, "更新失败: "+err.Error())
		}
		return
	}
	response.Success(c, updatedUser)
}

func (h *userHandler) ChangePassword(c *gin.Context) {
	var req *request.ChangePasswordRequest
	userID, exist := c.Get("userID")
	if !exist {
		response.Unauthorized(c, "用户未认证")
		return
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	err := h.userService.ChangePassword(c.Request.Context(), userID.(uint), req)
	if err != nil {
		switch err.Error() {
		case "用户不存在":
			response.NotFound(c, err.Error())
		case "原密码不正确":
			response.BadRequest(c, err.Error())
		case "新密码和确认密码不一致":
			response.BadRequest(c, err.Error())
		default:
			response.InternalServerError(c, err.Error())
		}
		return
	}
	response.Success(c, gin.H{
		"message": "修改密码成功",
		"user_id": userID,
	})

}
