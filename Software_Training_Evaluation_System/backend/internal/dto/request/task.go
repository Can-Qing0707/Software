package request

import "time"

type TaskQuery struct {
	CourseID uint   `form:"course_id"`
	Keyword  string `form:"keyword"`
}

type CreateTaskReq struct {
	CourseID    uint      `json:"course_id" binding:"required"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	Deadline    *time.Time `json:"deadline"`
}

type UpdateTaskReq struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Deadline    *time.Time `json:"deadline"`
	Status      *int       `json:"status"`
}
