package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"training_eval_system/internal/dto/request"
	"training_eval_system/internal/model"
	"training_eval_system/internal/repository"
)

type EvaluationService struct {
	evalRepo    *repository.EvalRepo
	subRepo     *repository.SubmissionRepo
	taskRepo    *repository.TaskRepo
	llmService  *LLMService
	fileService *FileService
}

func NewEvaluationService(evalRepo *repository.EvalRepo, subRepo *repository.SubmissionRepo, taskRepo *repository.TaskRepo, llmService *LLMService, fileService *FileService) *EvaluationService {
	return &EvaluationService{evalRepo: evalRepo, subRepo: subRepo, taskRepo: taskRepo, llmService: llmService, fileService: fileService}
}

func (s *EvaluationService) ListIndicators() ([]model.EvalIndicator, error) {
	return s.evalRepo.ListIndicators()
}

func (s *EvaluationService) GetAllIndicators() ([]model.EvalIndicator, error) {
	return s.evalRepo.GetAllIndicators()
}

func (s *EvaluationService) CreateIndicator(req *request.CreateIndicatorReq) error {
	ind := &model.EvalIndicator{
		Name:          req.Name,
		Description:   req.Description,
		DefaultWeight: req.DefaultWeight,
		SortOrder:     req.SortOrder,
		Status:        1,
	}
	return s.evalRepo.CreateIndicator(ind)
}

func (s *EvaluationService) UpdateIndicator(id uint, req *request.UpdateIndicatorReq) error {
	data := make(map[string]interface{})
	if req.Name != "" {
		data["name"] = req.Name
	}
	if req.Description != "" {
		data["description"] = req.Description
	}
	data["default_weight"] = req.DefaultWeight
	data["sort_order"] = req.SortOrder
	if req.Status != nil {
		data["status"] = *req.Status
	}
	return s.evalRepo.UpdateIndicator(id, data)
}

func (s *EvaluationService) DeleteIndicator(id uint) error {
	return s.evalRepo.DeleteIndicator(id)
}

func (s *EvaluationService) GetTaskIndicators(taskID uint) ([]model.TaskIndicator, error) {
	return s.evalRepo.GetTaskIndicators(taskID)
}

func (s *EvaluationService) SaveTaskIndicators(taskID uint, req *request.SaveTaskIndicatorsReq) error {
	items := make([]model.TaskIndicator, len(req.Indicators))
	for i, item := range req.Indicators {
		items[i] = model.TaskIndicator{
			TaskID:      taskID,
			IndicatorID: item.IndicatorID,
			Weight:      item.Weight,
		}
	}
	return s.evalRepo.SaveTaskIndicators(taskID, items)
}

func (s *EvaluationService) GetCourseIndicators(courseID uint) ([]model.CourseIndicator, error) {
	return s.evalRepo.GetCourseIndicators(courseID)
}

func (s *EvaluationService) SaveCourseIndicators(courseID uint, req *request.SaveCourseIndicatorsReq) error {
	items := make([]model.CourseIndicator, len(req.Indicators))
	for i, item := range req.Indicators {
		items[i] = model.CourseIndicator{
			CourseID:    courseID,
			IndicatorID: item.IndicatorID,
			Weight:      item.Weight,
		}
	}
	return s.evalRepo.SaveCourseIndicators(courseID, items)
}

func (s *EvaluationService) GetScores(submissionID uint) ([]model.EvalScore, error) {
	return s.evalRepo.GetScores(submissionID)
}

func (s *EvaluationService) SubmitTeacherScore(req *request.TeacherScoreReq) error {
	scores := make([]model.EvalScore, len(req.Scores))
	for i, item := range req.Scores {
		ts := item.TeacherScore
		scores[i] = model.EvalScore{
			SubmissionID:   req.SubmissionID,
			IndicatorID:    item.IndicatorID,
			TeacherScore:   &ts,
			TeacherComment: item.TeacherComment,
		}
	}
	if err := s.evalRepo.SaveTeacherScores(req.SubmissionID, scores); err != nil {
		return err
	}
	if err := s.computeFinalScores(req.SubmissionID); err != nil {
		return err
	}
	_ = s.subRepo.Update(req.SubmissionID, map[string]interface{}{"status": "evaluated"})
	return nil
}

func (s *EvaluationService) computeFinalScores(submissionID uint) error {
	sub, err := s.subRepo.FindByID(submissionID)
	if err != nil {
		return errors.New("提交记录不存在")
	}

	task, err := s.taskRepo.FindByID(sub.TaskID)
	if err != nil {
		return errors.New("关联任务不存在")
	}

	indicators, _ := s.evalRepo.GetTaskIndicators(sub.TaskID)
	weightMap := make(map[uint]float64)
	for _, ind := range indicators {
		weightMap[ind.IndicatorID] = ind.Weight
	}

	if len(weightMap) == 0 {
		courseIndicators, _ := s.evalRepo.GetCourseIndicators(task.CourseID)
		for _, ind := range courseIndicators {
			weightMap[ind.IndicatorID] = ind.Weight
		}
	}

	scores, err := s.evalRepo.GetScores(submissionID)
	if err != nil {
		return err
	}

	for _, sc := range scores {
		scoreVal := sc.TeacherScore
		if scoreVal == nil {
			scoreVal = sc.LLMScore
		}
		if scoreVal == nil {
			continue
		}

		weight := weightMap[sc.IndicatorID]
		if weight <= 0 {
			continue
		}

		final := (*scoreVal) * weight / 100.0
		_ = s.evalRepo.UpdateScoreFinal(sc.ID, final)
	}

	return nil
}

type llmScoreItem struct {
	IndicatorName string  `json:"indicator_name"`
	Score         float64 `json:"score"`
	Comment       string  `json:"comment"`
}

func (s *EvaluationService) ScoreByLLM(submissionID uint) error {
	sub, err := s.subRepo.FindByID(submissionID)
	if err != nil {
		return errors.New("提交记录不存在")
	}

	task, err := s.taskRepo.FindByID(sub.TaskID)
	if err != nil {
		return errors.New("关联任务不存在")
	}

	indicators, err := s.evalRepo.GetTaskIndicators(sub.TaskID)
	if err != nil || len(indicators) == 0 {
		courseIndicators, err2 := s.evalRepo.GetCourseIndicators(task.CourseID)
		if err2 != nil || len(courseIndicators) == 0 {
			return errors.New("该任务和课程均未配置评价指标")
		}
		for _, ci := range courseIndicators {
			indicators = append(indicators, model.TaskIndicator{
				IndicatorID:   ci.IndicatorID,
				Weight:        ci.Weight,
				IndicatorName: ci.IndicatorName,
			})
		}
	}

	contentText := sub.ContentText
	if contentText == "" {
		contentText = s.fileService.ExtractContent(sub.FilesJSON)
	}

	requirements := task.Description
	if requirements == "" {
		requirements = task.Title
	}

	indicatorNames := make([]string, len(indicators))
	for i, ind := range indicators {
		indicatorNames[i] = ind.IndicatorName
	}

	llmResponse, err := s.llmService.ScoreContent(contentText, requirements, indicatorNames)
	if err != nil {
		return fmt.Errorf("LLM评分失败: %w", err)
	}

	scoreItems, err := parseLLMScoreResponse(llmResponse)
	if err != nil {
		return fmt.Errorf("解析LLM评分结果失败: %w", err)
	}

	for _, item := range scoreItems {
		var indicatorID uint
		for _, ind := range indicators {
			if ind.IndicatorName == item.IndicatorName {
				indicatorID = ind.IndicatorID
				break
			}
		}
		if indicatorID == 0 {
			continue
		}

		llmScore := item.Score
		score := &model.EvalScore{
			SubmissionID: submissionID,
			IndicatorID:  indicatorID,
			LLMScore:     &llmScore,
			LLMComment:   item.Comment,
		}
		_ = s.evalRepo.UpsertScore(score)
	}

	if err := s.computeFinalScores(submissionID); err != nil {
		return err
	}

	_ = s.subRepo.Update(submissionID, map[string]interface{}{"status": "evaluated"})
	return nil
}

func parseLLMScoreResponse(raw string) ([]llmScoreItem, error) {
	text := strings.TrimSpace(raw)

	if strings.HasPrefix(text, "```") {
		idx := strings.Index(text, "\n")
		if idx != -1 {
			text = text[idx+1:]
		}
		if idx := strings.LastIndex(text, "```"); idx != -1 {
			text = text[:idx]
		}
		text = strings.TrimSpace(text)
	}

	var result []llmScoreItem
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		start := strings.Index(text, "[")
		end := strings.LastIndex(text, "]")
		if start != -1 && end > start {
			if err2 := json.Unmarshal([]byte(text[start:end+1]), &result); err2 != nil {
				return nil, err
			}
			return result, nil
		}
		return nil, err
	}
	return result, nil
}
