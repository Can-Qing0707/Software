package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"training_eval_system/internal/dto/request"
	"training_eval_system/internal/service"
	"training_eval_system/pkg/response"
)

type EvaluationHandler struct {
	evalService *service.EvaluationService
}

func NewEvaluationHandler(evalService *service.EvaluationService) *EvaluationHandler {
	return &EvaluationHandler{evalService: evalService}
}

func (h *EvaluationHandler) ListIndicators(c *gin.Context) {
	indicators, err := h.evalService.ListIndicators()
	if err != nil {
		response.InternalError(c, "获取指标列表失败")
		return
	}
	response.Success(c, indicators)
}

func (h *EvaluationHandler) GetAllIndicators(c *gin.Context) {
	indicators, err := h.evalService.GetAllIndicators()
	if err != nil {
		response.InternalError(c, "获取指标列表失败")
		return
	}
	response.Success(c, indicators)
}

func (h *EvaluationHandler) CreateIndicator(c *gin.Context) {
	var req request.CreateIndicatorReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请检查输入信息")
		return
	}
	if err := h.evalService.CreateIndicator(&req); err != nil {
		response.InternalError(c, "创建指标失败")
		return
	}
	response.Created(c, nil)
}

func (h *EvaluationHandler) UpdateIndicator(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的指标ID")
		return
	}
	var req request.UpdateIndicatorReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请检查输入信息")
		return
	}
	if err := h.evalService.UpdateIndicator(uint(id), &req); err != nil {
		response.InternalError(c, "更新失败")
		return
	}
	response.Success(c, nil)
}

func (h *EvaluationHandler) DeleteIndicator(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的指标ID")
		return
	}
	if err := h.evalService.DeleteIndicator(uint(id)); err != nil {
		response.InternalError(c, "删除失败")
		return
	}
	response.Success(c, nil)
}

func (h *EvaluationHandler) GetTaskIndicators(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("taskId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的任务ID")
		return
	}
	indicators, err := h.evalService.GetTaskIndicators(uint(taskID))
	if err != nil {
		response.InternalError(c, "获取任务指标失败")
		return
	}
	response.Success(c, indicators)
}

func (h *EvaluationHandler) SaveTaskIndicators(c *gin.Context) {
	taskID, err := strconv.ParseUint(c.Param("taskId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的任务ID")
		return
	}
	var req request.SaveTaskIndicatorsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请检查输入信息")
		return
	}
	if err := h.evalService.SaveTaskIndicators(uint(taskID), &req); err != nil {
		response.InternalError(c, "保存任务指标失败")
		return
	}
	response.Success(c, nil)
}

func (h *EvaluationHandler) GetCourseIndicators(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Param("courseId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的课程ID")
		return
	}
	indicators, err := h.evalService.GetCourseIndicators(uint(courseID))
	if err != nil {
		response.InternalError(c, "获取课程指标失败")
		return
	}
	response.Success(c, indicators)
}

func (h *EvaluationHandler) SaveCourseIndicators(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Param("courseId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的课程ID")
		return
	}
	var req request.SaveCourseIndicatorsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请检查输入信息")
		return
	}
	if err := h.evalService.SaveCourseIndicators(uint(courseID), &req); err != nil {
		response.InternalError(c, "保存课程指标失败")
		return
	}
	response.Success(c, nil)
}

func (h *EvaluationHandler) GetScores(c *gin.Context) {
	submissionID, err := strconv.ParseUint(c.Param("submissionId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的提交ID")
		return
	}
	scores, err := h.evalService.GetScores(uint(submissionID))
	if err != nil {
		response.InternalError(c, "获取评分失败")
		return
	}
	response.Success(c, scores)
}

func (h *EvaluationHandler) SubmitTeacherScore(c *gin.Context) {
	var req request.TeacherScoreReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请检查输入信息")
		return
	}
	if err := h.evalService.SubmitTeacherScore(&req); err != nil {
		response.InternalError(c, "保存评分失败")
		return
	}
	response.Success(c, nil)
}

func (h *EvaluationHandler) ScoreByLLM(c *gin.Context) {
	submissionID, err := strconv.ParseUint(c.Param("submissionId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的提交ID")
		return
	}
	if err := h.evalService.ScoreByLLM(uint(submissionID)); err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, "LLM评分完成")
}
