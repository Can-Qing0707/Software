package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"training_eval_system/internal/service"
	"training_eval_system/pkg/response"
)

type VerifyHandler struct {
	verifyService *service.VerifyService
}

func NewVerifyHandler(verifyService *service.VerifyService) *VerifyHandler {
	return &VerifyHandler{verifyService: verifyService}
}

func (h *VerifyHandler) Verify(c *gin.Context) {
	submissionID, err := strconv.ParseUint(c.Param("submissionId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的提交ID")
		return
	}
	vr, err := h.verifyService.Verify(uint(submissionID))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, vr)
}

func (h *VerifyHandler) GetVerification(c *gin.Context) {
	submissionID, err := strconv.ParseUint(c.Param("submissionId"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的提交ID")
		return
	}
	vr, err := h.verifyService.GetVerification(uint(submissionID))
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Success(c, vr)
}
