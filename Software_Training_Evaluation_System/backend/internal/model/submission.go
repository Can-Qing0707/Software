package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type FileItem struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Type string `json:"type"`
	Size int64  `json:"size"`
}

type FileList []FileItem

func (f FileList) Value() (driver.Value, error) {
	return json.Marshal(f)
}

func (f *FileList) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, f)
}

type Submission struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	TaskID      uint      `gorm:"not null;index;comment:关联任务ID" json:"task_id"`
	StudentID   uint      `gorm:"not null;index;comment:提交学生ID" json:"student_id"`
	FilesJSON   FileList  `gorm:"type:json;not null;comment:文件列表JSON" json:"files"`
	ContentText string    `gorm:"type:longtext;default:null;comment:LLM提取的文本内容" json:"content_text,omitempty"`
	Status      string    `gorm:"size:32;not null;default:uploaded;index;comment:提交状态" json:"status"`
	SubmitTime  time.Time `gorm:"autoCreateTime;comment:提交时间" json:"submit_time"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	StudentName string    `gorm:"->" json:"student_name,omitempty"`
	TaskTitle   string    `gorm:"->" json:"task_title,omitempty"`
	TotalScore   *float64 `gorm:"->" json:"total_score,omitempty"`
	LlmTotal     *float64 `gorm:"->" json:"llm_total,omitempty"`
	TeacherTotal *float64 `gorm:"->" json:"teacher_total,omitempty"`
}

func (Submission) TableName() string {
	return "submissions"
}
