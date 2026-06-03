package handler

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	"training_eval_system/internal/dto/request"
	"training_eval_system/internal/middleware"
	"training_eval_system/internal/model"
	"training_eval_system/internal/service"
	"training_eval_system/pkg/response"
)

type SubmissionHandler struct {
	submissionService *service.SubmissionService
	fileService       *service.FileService
}

func NewSubmissionHandler(submissionService *service.SubmissionService, fileService *service.FileService) *SubmissionHandler {
	return &SubmissionHandler{submissionService: submissionService, fileService: fileService}
}

func (h *SubmissionHandler) List(c *gin.Context) {
	var req request.SubmissionQuery
	if err := c.ShouldBindQuery(&req); err != nil {
		req = request.SubmissionQuery{}
	}
	role := middleware.GetRole(c)
	userID := middleware.GetUserID(c)
	list, err := h.submissionService.List(role, userID, &req)
	if err != nil {
		response.InternalError(c, "获取提交列表失败")
		return
	}
	response.Success(c, list)
}

func (h *SubmissionHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的提交ID")
		return
	}
	sub, err := h.submissionService.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "提交记录不存在")
		return
	}
	response.Success(c, sub)
}

func (h *SubmissionHandler) Create(c *gin.Context) {
	var req request.CreateSubmissionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请检查输入信息")
		return
	}
	studentID := middleware.GetUserID(c)
	if err := h.submissionService.Create(&req, studentID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Created(c, nil)
}

func (h *SubmissionHandler) Resubmit(c *gin.Context) {
	var req request.CreateSubmissionReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请检查输入信息")
		return
	}
	studentID := middleware.GetUserID(c)
	if err := h.submissionService.Resubmit(&req, studentID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.Success(c, "重新提交成功")
}

func (h *SubmissionHandler) DownloadFile(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的提交ID")
		return
	}
	idx, err := strconv.Atoi(c.Param("idx"))
	if err != nil {
		response.BadRequest(c, "无效的文件索引")
		return
	}

	sub, err := h.submissionService.GetByID(uint(id))
	if err != nil {
		response.NotFound(c, "提交记录不存在")
		return
	}

	role := middleware.GetRole(c)
	userID := middleware.GetUserID(c)
	subStudentID, _ := sub["student_id"].(uint)
	if role == "student" && subStudentID != userID {
		response.Forbidden(c, "无权下载此文件")
		return
	}

	files, ok := sub["files"].(model.FileList)
	if !ok || idx < 0 || idx >= len(files) {
		response.NotFound(c, "文件不存在")
		return
	}

	fileItem := files[idx]
	filePath := h.fileService.GetFilePath(fileItem.URL)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		response.NotFound(c, "文件不存在")
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename*=UTF-8''%s`, fileItem.Name))
	c.Header("Content-Type", "application/octet-stream")
	c.File(filePath)
}
