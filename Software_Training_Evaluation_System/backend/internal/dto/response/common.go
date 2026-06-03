package response

type PageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

type LoginResp struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

type ScoreResp struct {
	IndicatorID    uint    `json:"indicator_id"`
	IndicatorName  string  `json:"indicator_name"`
	LLMScore       float64 `json:"llm_score,omitempty"`
	LLMComment     string  `json:"llm_comment,omitempty"`
	TeacherScore   float64 `json:"teacher_score,omitempty"`
	TeacherComment string  `json:"teacher_comment,omitempty"`
	FinalScore     float64 `json:"final_score,omitempty"`
	Weight         float64 `json:"weight,omitempty"`
}

type DashboardStats struct {
	Courses     int `json:"courses"`
	Tasks       int `json:"tasks"`
	Submissions int `json:"submissions"`
	Evaluated   int `json:"evaluated"`
}
