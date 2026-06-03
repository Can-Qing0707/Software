package model

import "time"

type Course struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"size:128;not null;comment:课程名称" json:"name"`
	Description string    `gorm:"type:text;comment:课程描述" json:"description,omitempty"`
	TeacherID   uint      `gorm:"not null;index;comment:授课教师ID" json:"teacher_id"`
	Semester    string    `gorm:"size:32;default:null;comment:学期" json:"semester,omitempty"`
	Code        string    `gorm:"size:32;uniqueIndex;not null;comment:课程代码" json:"code"`
	Status      int       `gorm:"not null;default:1;index;comment:状态 1-进行中 0-已结束" json:"status"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	TeacherName  string    `gorm:"->" json:"teacher_name,omitempty"`
	StudentCount int64     `gorm:"->" json:"student_count,omitempty"`
}

func (Course) TableName() string {
	return "courses"
}

type CourseEnrollment struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CourseID  uint      `gorm:"not null;uniqueIndex:uk_course_student;comment:课程ID" json:"course_id"`
	StudentID uint      `gorm:"not null;uniqueIndex:uk_course_student;comment:学生ID" json:"student_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

func (CourseEnrollment) TableName() string {
	return "course_enrollments"
}
