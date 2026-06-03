package model

import "time"

type Report struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	SubmissionID *uint     `gorm:"default:null;index;comment:关联提交ID(NULL=统计报表)" json:"submission_id,omitempty"`
	CourseID     *uint     `gorm:"default:null;index;comment:关联课程ID(统计报表)" json:"course_id,omitempty"`
	TaskID       *uint     `gorm:"default:null;index;comment:关联任务ID(统计报表)" json:"task_id,omitempty"`
	Type         string    `gorm:"size:32;not null;index;comment:报告类型" json:"type"`
	Format       string    `gorm:"size:16;not null;comment:导出格式" json:"format"`
	FileURL      string    `gorm:"size:512;not null;comment:文件存储路径" json:"file_url"`
	Title        string    `gorm:"size:256;default:null;comment:报告标题" json:"title,omitempty"`
	GeneratedBy  uint      `gorm:"not null;comment:生成人ID" json:"generated_by"`
	GeneratedAt  time.Time `gorm:"autoCreateTime;comment:生成时间" json:"generated_at"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (Report) TableName() string {
	return "reports"
}
