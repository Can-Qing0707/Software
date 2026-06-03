package request

type CreateIndicatorReq struct {
	Name          string  `json:"name" binding:"required"`
	Description   string  `json:"description"`
	DefaultWeight float64 `json:"default_weight" binding:"min=0,max=100"`
	SortOrder     int     `json:"sort_order"`
}

type UpdateIndicatorReq struct {
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	DefaultWeight float64 `json:"default_weight" binding:"min=0,max=100"`
	SortOrder     int     `json:"sort_order"`
	Status        *int    `json:"status"`
}

type TaskIndicatorReq struct {
	IndicatorID uint    `json:"indicator_id" binding:"required"`
	Weight      float64 `json:"weight" binding:"min=0,max=100"`
}

type SaveTaskIndicatorsReq struct {
	Indicators []TaskIndicatorReq `json:"indicators" binding:"required"`
}

type CourseIndicatorItem struct {
	IndicatorID uint    `json:"indicator_id" binding:"required"`
	Weight      float64 `json:"weight" binding:"min=0,max=100"`
}

type SaveCourseIndicatorsReq struct {
	Indicators []CourseIndicatorItem `json:"indicators" binding:"required"`
}

type TeacherScoreReq struct {
	SubmissionID uint               `json:"submission_id" binding:"required"`
	Scores       []TeacherScoreItem `json:"scores" binding:"required"`
}

type TeacherScoreItem struct {
	IndicatorID    uint    `json:"indicator_id" binding:"required"`
	TeacherScore   float64 `json:"teacher_score" binding:"min=0,max=100"`
	TeacherComment string  `json:"teacher_comment"`
}
