package model

import "time"

type SystemConfig struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ConfigKey   string    `gorm:"size:128;not null;uniqueIndex;comment:配置键" json:"config_key"`
	ConfigValue string    `gorm:"type:text;not null;comment:配置值(JSON)" json:"config_value"`
	Description string    `gorm:"size:256;default:null;comment:配置说明" json:"description,omitempty"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (SystemConfig) TableName() string {
	return "system_config"
}
