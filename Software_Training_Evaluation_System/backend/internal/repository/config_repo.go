package repository

import (
	"training_eval_system/internal/model"

	"gorm.io/gorm"
)

type ConfigRepo struct {
	db *gorm.DB
}

func NewConfigRepo(db *gorm.DB) *ConfigRepo {
	return &ConfigRepo{db: db}
}

func (r *ConfigRepo) GetByKey(key string) (*model.SystemConfig, error) {
	var config model.SystemConfig
	err := r.db.Where("config_key = ?", key).First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *ConfigRepo) Upsert(config *model.SystemConfig) error {
	return r.db.Where("config_key = ?", config.ConfigKey).
		Assign(model.SystemConfig{ConfigValue: config.ConfigValue, Description: config.Description}).
		FirstOrCreate(config).Error
}
