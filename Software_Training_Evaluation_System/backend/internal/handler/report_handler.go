package handler

import (
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"

	"training_eval_system/internal/dto/request"
	"training_eval_system/internal/middleware"
	"training_eval_system/internal/service"
	"training_eval_system/pkg/response"
)

type ReportHandler struct {
	reportService *service.ReportService
}

func NewReportHandler(reportService *service.ReportService) *ReportHandler {
	return &ReportHandler{reportService: reportService}
}

func (h *ReportHandler) List(c *gin.Context) {
	var req request.ReportQuery
	if err := c.ShouldBindQuery(&req); err != nil {
		req = request.ReportQuery{}
	}
	reports, err := h.reportService.List(&req)
	if err != nil {
		response.InternalError(c, "获取报表列表失败")
		return
	}
	response.Success(c, reports)
}

func (h *ReportHandler) Generate(c *gin.Context) {
	var req request.GenerateReportReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请检查输入信息，需要type和format字段")
		return
	}
	userID := middleware.GetUserID(c)
	report, err := h.reportService.Generate(&req, userID)
	if err != nil {
		response.InternalError(c, err.Error())
		return
	}
	response.Created(c, report)
}

func (h *ReportHandler) Export(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的报表ID")
		return
	}

	filePath, contentType, err := h.reportService.GetExportPath(uint(id))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filepath.Base(filePath)))
	c.File(filePath)
}
