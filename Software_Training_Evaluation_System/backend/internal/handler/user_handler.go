package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"training_eval_system/internal/dto/request"
	"training_eval_system/internal/middleware"
	"training_eval_system/internal/service"
	"training_eval_system/pkg/response"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	user, err := h.userService.GetProfile(userID)
	if err != nil {
		response.NotFound(c, "用户不存在")
		return
	}
	response.Success(c, user)
}

func (h *UserHandler) List(c *gin.Context) {
	var req request.UserQuery
	if err := c.ShouldBindQuery(&req); err != nil {
		req = request.UserQuery{}
	}
	users, err := h.userService.List(&req)
	if err != nil {
		response.InternalError(c, "获取用户列表失败")
		return
	}
	response.Success(c, users)
}

func (h *UserHandler) Create(c *gin.Context) {
	var req request.CreateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请检查输入信息")
		return
	}
	if err := h.userService.Create(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, nil)
}

func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}
	var req request.UpdateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请检查输入信息")
		return
	}
	if err := h.userService.Update(uint(id), &req); err != nil {
		response.InternalError(c, "更新失败")
		return
	}
	response.Success(c, nil)
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的用户ID")
		return
	}
	if err := h.userService.Delete(uint(id)); err != nil {
		response.InternalError(c, "删除失败")
		return
	}
	response.Success(c, nil)
}
