package model

import "time"

type Task struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CourseID      uint      `gorm:"not null;index;comment:所属课程ID" json:"course_id"`
	Title         string    `gorm:"size:256;not null;comment:任务标题" json:"title"`
	Description   string    `gorm:"type:text;comment:任务描述/要求" json:"description,omitempty"`
	AttachmentURL string    `gorm:"size:512;default:null;comment:任务附件URL" json:"attachment_url,omitempty"`
	Deadline      *time.Time `gorm:"default:null;comment:截止时间" json:"deadline,omitempty"`
	Status        int       `gorm:"not null;default:1;index;comment:状态 1-发布 0-草稿" json:"status"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	CourseName    string    `gorm:"->" json:"course_name,omitempty"`
}

func (Task) TableName() string {
	return "tasks"
}
