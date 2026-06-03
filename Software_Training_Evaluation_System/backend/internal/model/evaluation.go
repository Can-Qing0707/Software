package model

import "time"

type EvalIndicator struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name          string    `gorm:"size:128;not null;comment:指标名称" json:"name"`
	Description   string    `gorm:"size:512;default:null;comment:指标说明/评价标准" json:"description,omitempty"`
	DefaultWeight float64   `gorm:"type:decimal(5,2);not null;default:0;comment:默认权重(0-100)" json:"default_weight"`
	SortOrder     int       `gorm:"not null;default:0;comment:排序号" json:"sort_order"`
	Status        int       `gorm:"not null;default:1;index;comment:状态 1-启用 0-禁用" json:"status"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (EvalIndicator) TableName() string {
	return "eval_indicators"
}

type TaskIndicator struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	TaskID      uint      `gorm:"not null;uniqueIndex:uk_task_indicator;comment:关联任务ID" json:"task_id"`
	IndicatorID uint      `gorm:"not null;uniqueIndex:uk_task_indicator;comment:关联指标ID" json:"indicator_id"`
	Weight      float64   `gorm:"type:decimal(5,2);not null;default:0;comment:权重(0-100)" json:"weight"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	IndicatorName string  `gorm:"->" json:"name,omitempty"`
}

func (TaskIndicator) TableName() string {
	return "task_indicators"
}

type CourseIndicator struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CourseID      uint      `gorm:"not null;uniqueIndex:uk_course_indicator" json:"course_id"`
	IndicatorID   uint      `gorm:"not null;uniqueIndex:uk_course_indicator" json:"indicator_id"`
	Weight        float64   `gorm:"type:decimal(5,2);not null;default:0" json:"weight"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	IndicatorName string    `gorm:"->" json:"name,omitempty"`
}

func (CourseIndicator) TableName() string {
	return "course_indicators"
}

type EvalScore struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	SubmissionID   uint      `gorm:"not null;uniqueIndex:uk_submission_indicator;comment:关联提交ID" json:"submission_id"`
	IndicatorID    uint      `gorm:"not null;uniqueIndex:uk_submission_indicator;comment:关联指标ID" json:"indicator_id"`
	LLMScore       *float64  `gorm:"type:decimal(5,2);default:null;comment:LLM评分(0-100)" json:"llm_score,omitempty"`
	LLMComment     string    `gorm:"type:text;default:null;comment:LLM评语" json:"llm_comment,omitempty"`
	TeacherScore   *float64  `gorm:"type:decimal(5,2);default:null;comment:教师评分(0-100)" json:"teacher_score,omitempty"`
	TeacherComment string    `gorm:"type:text;default:null;comment:教师评语" json:"teacher_comment,omitempty"`
	FinalScore     *float64  `gorm:"type:decimal(5,2);default:null;comment:最终得分(加权计算)" json:"final_score,omitempty"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	IndicatorName  string    `gorm:"->" json:"indicator_name,omitempty"`
}

func (EvalScore) TableName() string {
	return "eval_scores"
}

type VerificationResult struct {
	ID               uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	SubmissionID     uint      `gorm:"not null;uniqueIndex;comment:关联提交ID" json:"submission_id"`
	Completeness     string    `gorm:"type:text;comment:步骤完整性核查结果(JSON)" json:"completeness,omitempty"`
	LogicIssues      string    `gorm:"type:text;comment:逻辑漏洞核查结果(JSON)" json:"logic_issues,omitempty"`
	RequirementMatch string    `gorm:"type:text;comment:要求匹配度核查结果(JSON)" json:"requirement_match,omitempty"`
	OverallPass      *int      `gorm:"type:tinyint(1);default:null;comment:是否整体通过 1-通过 0-不通过" json:"overall_pass,omitempty"`
	RawLLMResponse   string    `gorm:"type:longtext;default:null;comment:LLM原始返回" json:"-"`
	VerifiedAt       *time.Time `gorm:"default:null;comment:核查时间" json:"verified_at,omitempty"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (VerificationResult) TableName() string {
	return "verification_results"
}
