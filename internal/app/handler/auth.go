package handler

import (
	"go-todo-api/internal/app/dto/request"
	"go-todo-api/internal/service"

	"go-todo-api/pkg/response"

	"github.com/gin-gonic/gin"
)

type authHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *authHandler {
	return &authHandler{authService: authService}
}

func (h *authHandler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userResp, err := h.authService.Register(c, &req)
	if err != nil {
		if err.Error() == "用户名已存在" || err.Error() == "邮箱已存在" {
			response.BadRequest(c, err.Error())
		} else {
			response.InternalServerError(c, err.Error())
		}
		return
	}
	response.Success(c, userResp)
}

func (h *authHandler) Login(c *gin.Context) {

	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userResp, err := h.authService.Login(c, &req)
	if err != nil {
		if err.Error() == "用户名或密码错误" {
			response.BadRequest(c, err.Error())
		} else {
			response.InternalServerError(c, err.Error())
		}
		return
	}
	response.Success(c, userResp)
}
