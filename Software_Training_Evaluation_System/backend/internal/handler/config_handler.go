package handler

import (
	"encoding/json"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"training_eval_system/config"
	"training_eval_system/internal/dto/request"
	"training_eval_system/internal/model"
	"training_eval_system/internal/repository"
	"training_eval_system/pkg/response"
)

type ConfigHandler struct {
	configRepo *repository.ConfigRepo
}

func NewConfigHandler(configRepo *repository.ConfigRepo) *ConfigHandler {
	return &ConfigHandler{configRepo: configRepo}
}

func (h *ConfigHandler) GetLlmConfig(c *gin.Context) {
	cfg, err := h.configRepo.GetByKey("llm_provider")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Success(c, config.AppConfig.LLM)
			return
		}
		response.InternalError(c, "获取配置失败")
		return
	}
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(cfg.ConfigValue), &result); err == nil {
		response.Success(c, result)
	} else {
		response.Success(c, cfg.ConfigValue)
	}
}

func (h *ConfigHandler) UpdateLlmConfig(c *gin.Context) {
	var req request.LlmConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "请检查输入信息")
		return
	}
	jsonData, err := json.Marshal(req)
	if err != nil {
		response.InternalError(c, "配置序列化失败")
		return
	}
	sysConfig := &model.SystemConfig{
		ConfigKey:   "llm_provider",
		ConfigValue: string(jsonData),
		Description: "LLM服务配置",
	}
	if err := h.configRepo.Upsert(sysConfig); err != nil {
		response.InternalError(c, "保存配置失败")
		return
	}

	config.SetLLMConfig(config.LLMConfig{
		Provider:    req.Provider,
		APIURL:      req.APIURL,
		APIKey:      req.APIKey,
		Model:       req.Model,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
	})

	response.Success(c, nil)
}
