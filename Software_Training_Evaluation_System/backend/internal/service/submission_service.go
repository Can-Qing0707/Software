package service

import (
	"encoding/json"
	"errors"
	"time"

	"training_eval_system/internal/dto/request"
	"training_eval_system/internal/model"
	"training_eval_system/internal/repository"
)

type SubmissionService struct {
	subRepo  *repository.SubmissionRepo
	evalRepo *repository.EvalRepo
	taskRepo *repository.TaskRepo
}

func NewSubmissionService(subRepo *repository.SubmissionRepo, evalRepo *repository.EvalRepo, taskRepo *repository.TaskRepo) *SubmissionService {
	return &SubmissionService{subRepo: subRepo, evalRepo: evalRepo, taskRepo: taskRepo}
}

func (s *SubmissionService) List(role string, userID uint, req *request.SubmissionQuery) ([]model.Submission, error) {
	switch role {
	case "student":
		return s.subRepo.ListByStudent(userID, req.TaskID, req.Status)
	case "teacher":
		return s.subRepo.ListByTeacher(userID, req.TaskID, req.Status)
	default:
		return s.subRepo.List(req.TaskID, req.Status)
	}
}

func (s *SubmissionService) GetByID(id uint) (map[string]interface{}, error) {
	sub, err := s.subRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	result := map[string]interface{}{
		"id":           sub.ID,
		"task_id":      sub.TaskID,
		"student_id":   sub.StudentID,
		"files":        sub.FilesJSON,
		"content_text": sub.ContentText,
		"status":       sub.Status,
		"submit_time":  sub.SubmitTime,
		"created_at":   sub.CreatedAt,
		"updated_at":   sub.UpdatedAt,
		"student_name": sub.StudentName,
		"task_title":   sub.TaskTitle,
	}

	vr, _ := s.evalRepo.GetVerificationBySubmission(id)
	if vr != nil {
		type completenessStruct struct {
			StepsTotal        int      `json:"steps_total"`
			StepsCompleted    int      `json:"steps_completed"`
			MissingSteps      []string `json:"missing_steps"`
			CompletenessRatio float64  `json:"completeness_ratio"`
		}
		type requirementStruct struct {
			MatchedRequirements   []string `json:"matched_requirements"`
			UnmatchedRequirements []string `json:"unmatched_requirements"`
			MatchRatio            float64  `json:"match_ratio"`
		}
		type logicStruct struct {
			HasLogicIssues bool     `json:"has_logic_issues"`
			Issues         []string `json:"issues"`
			Summary        string   `json:"summary"`
		}
		vData := map[string]interface{}{
			"overall_pass": vr.OverallPass,
		}
		if vr.Completeness != "" {
			var c completenessStruct
			if err := json.Unmarshal([]byte(vr.Completeness), &c); err == nil {
				vData["completeness"] = c
			}
		}
		if vr.LogicIssues != "" {
			var l logicStruct
			if err := json.Unmarshal([]byte(vr.LogicIssues), &l); err == nil {
				vData["logic_issues"] = l
			}
		}
		if vr.RequirementMatch != "" {
			var r requirementStruct
			if err := json.Unmarshal([]byte(vr.RequirementMatch), &r); err == nil {
				vData["requirement_match"] = r
			}
		}
		result["verification"] = vData
	}

	return result, nil
}

func (s *SubmissionService) Create(req *request.CreateSubmissionReq, studentID uint) error {
	existing, _ := s.subRepo.GetByTaskAndStudent(req.TaskID, studentID)
	if existing != nil {
		return errors.New("该任务已提交过成果")
	}
	sub := &model.Submission{
		TaskID:     req.TaskID,
		StudentID:  studentID,
		FilesJSON:  make(model.FileList, 0),
		Status:     "uploaded",
		SubmitTime: time.Now(),
	}
	for _, f := range req.Files {
		sub.FilesJSON = append(sub.FilesJSON, model.FileItem{
			Name: f.Name,
			URL:  f.URL,
			Type: f.Type,
			Size: f.Size,
		})
	}
	return s.subRepo.Create(sub)
}

func (s *SubmissionService) Resubmit(req *request.CreateSubmissionReq, studentID uint) error {
	existing, err := s.subRepo.GetByTaskAndStudent(req.TaskID, studentID)
	if err != nil {
		return errors.New("未找到已提交记录，请先提交成果")
	}
	files := make(model.FileList, 0)
	for _, f := range req.Files {
		files = append(files, model.FileItem{
			Name: f.Name,
			URL:  f.URL,
			Type: f.Type,
			Size: f.Size,
		})
	}
	data := map[string]interface{}{
		"files_json":   files,
		"status":       "uploaded",
		"submit_time":  time.Now(),
		"content_text": "",
	}
	return s.subRepo.Update(existing.ID, data)
}

func (s *SubmissionService) Update(id uint, data map[string]interface{}) error {
	return s.subRepo.Update(id, data)
}
