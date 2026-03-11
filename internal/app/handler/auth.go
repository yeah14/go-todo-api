package handler

import (
	"fmt"
	"go-todo-api/internal/app/dto/request"
	"go-todo-api/internal/service"
	"strings"
	"time"

	"go-todo-api/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type authHandler struct {
	authService  service.AuthService
	blacklistSvc service.BlacklistService
}

func NewAuthHandler(authService service.AuthService, blacklistSvc service.BlacklistService) *authHandler {
	return &authHandler{authService: authService, blacklistSvc: blacklistSvc}
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

func (h *authHandler) Logout(c *gin.Context) {
	tokenString, exists := c.Get("tokenString")
	if !exists {
		// 如果中间件没存，也可以从Header再取一次
		authHeader := c.GetHeader("Authorization")
		parts := strings.SplitN(authHeader, " ", 2)
		tokenString = parts[1]
	}

	token, _ := jwt.Parse(tokenString.(string), nil) // 这里不验证签名，只解析
	claims, _ := token.Claims.(jwt.MapClaims)
	exp := claims["exp"].(float64)
	expiresAt := time.Unix(int64(exp), 0)
	remainingTime := time.Until(expiresAt)
	fmt.Println(tokenString.(string))
	if remainingTime > 0 {

		err := h.blacklistSvc.AddtoBlacklist(c.Request.Context(), tokenString.(string), remainingTime)

		if err != nil {
			response.InternalServerError(c, err.Error()+tokenString.(string))
			return
		}

	}
	response.Success(c, "退出成功")
}
