package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"training_eval_system/internal/dto/request"
	"training_eval_system/internal/middleware"
	"training_eval_system/internal/model"
	"training_eval_system/internal/service"
	"training_eval_system/pkg/response"
)

type TaskHandler struct {
	taskService *service.TaskService
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

func (h *TaskHandler) List(c *gin.Context) {
	var req request.TaskQuery
	if err := c.ShouldBindQuery(&req); err != nil {
		req = request.TaskQuery{}
	}
	var tasks []model.Task
	var err error
	switch middleware.GetRole(c) {
	case "teacher":
		tasks, err = h.taskService.ListByTeacher(middleware.GetUserID(c), &req)
	case "student":
		tasks, err = h.taskService.ListByStudent(middleware.GetUserID(c), &req)
	default:
		tasks, err = h.taskService.List(&req)
	}
	if err != nil {
		response.InternalError(c, "获取任务列表失败")
		return
	}
	response.Success(c, tasks)
}

func (h *TaskHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的任务ID")
		return
	}
	task, err := h.taskService.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "任务不存在")
		return
	}
	response.Success(c, task)
}

func (h *TaskHandler) Create(c *gin.Context) {
	var req request.CreateTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请检查输入信息")
		return
	}
	if err := h.taskService.Create(&req); err != nil {
		response.InternalError(c, "创建任务失败")
		return
	}
	response.Created(c, nil)
}

func (h *TaskHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的任务ID")
		return
	}
	var req request.UpdateTaskReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请检查输入信息")
		return
	}
	if err := h.taskService.Update(uint(id), &req); err != nil {
		response.InternalError(c, "更新失败")
		return
	}
	response.Success(c, nil)
}

func (h *TaskHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的任务ID")
		return
	}
	if err := h.taskService.Delete(uint(id)); err != nil {
		response.InternalError(c, "删除失败")
		return
	}
	response.Success(c, nil)
}
