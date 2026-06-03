package request

type SubmissionQuery struct {
	TaskID uint   `form:"task_id"`
	Status string `form:"status"`
}

type CreateSubmissionReq struct {
	TaskID uint            `json:"task_id" binding:"required"`
	Files  []FileInfoReq `json:"files" binding:"required"`
}

type FileInfoReq struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Type string `json:"type"`
	Size int64  `json:"size"`
}
