package request

type CourseQuery struct {
	Keyword string `form:"keyword"`
}

type CreateCourseReq struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Semester    string `json:"semester" binding:"required"`
}

type UpdateCourseReq struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Semester    string `json:"semester"`
	Status      *int   `json:"status"`
}

type JoinCourseReq struct {
	Code string `json:"code" binding:"required"`
}
