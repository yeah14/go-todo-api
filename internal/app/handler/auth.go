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

// Register 用户注册
// @Summary 用户注册
// @Description 使用用户名、邮箱和密码注册新用户
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.RegisterRequest true "注册请求参数"
// @Success 200 {object} response.AuthResponse "注册成功"
// @Failure 400 {object} response.Response "请求参数错误或用户名/邮箱已存在"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /auth/register [post]
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

// Login 用户登录
// @Summary 用户登录
// @Description 使用用户名和密码登录
// @Tags auth
// @Accept json
// @Produce json
// @Param request body request.LoginRequest true "登录请求参数"
// @Success 200 {object} response.AuthResponse "登录成功，返回令牌"
// @Failure 400 {object} response.Response "请求参数错误或用户名密码错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /auth/login [post]
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
