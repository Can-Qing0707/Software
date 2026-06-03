package request

type ReportQuery struct {
	CourseID uint `form:"course_id"`
	TaskID   uint `form:"task_id"`
}

type GenerateReportReq struct {
	CourseID uint `json:"course_id"`
	TaskID   uint `json:"task_id"`
	Type     string `json:"type" binding:"required,oneof=individual class course"`
	Format   string `json:"format" binding:"required,oneof=pdf excel"`
}
