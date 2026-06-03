package handler

import (
	"github.com/gin-gonic/gin"

	"training_eval_system/config"
	"training_eval_system/internal/service"
	"training_eval_system/pkg/response"
)

type UploadHandler struct {
	fileService *service.FileService
}

func NewUploadHandler(fileService *service.FileService) *UploadHandler {
	return &UploadHandler{fileService: fileService}
}

func (h *UploadHandler) Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "请选择要上传的文件")
		return
	}

	maxSize := config.AppConfig.Upload.MaxSize * 1024 * 1024
	if file.Size > int64(maxSize) {
		response.BadRequest(c, "文件大小超过限制")
		return
	}

	f, err := file.Open()
	if err != nil {
		response.InternalError(c, "文件读取失败")
		return
	}
	defer f.Close()

	url, err := h.fileService.SaveFile(f, file.Filename)
	if err != nil {
		response.InternalError(c, "文件保存失败")
		return
	}

	response.Created(c, gin.H{
		"name": file.Filename,
		"url":  url,
		"type": file.Header.Get("Content-Type"),
		"size": file.Size,
	})
}
