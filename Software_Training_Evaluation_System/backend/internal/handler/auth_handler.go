package handler

import (
	"github.com/gin-gonic/gin"

	"training_eval_system/internal/dto/request"
	"training_eval_system/internal/service"
	"training_eval_system/pkg/response"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请输入用户名和密码")
		return
	}
	resp, err := h.authService.Login(&req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, resp)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请检查输入信息")
		return
	}
	if err := h.authService.Register(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, nil)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	response.Success(c, nil)
}
