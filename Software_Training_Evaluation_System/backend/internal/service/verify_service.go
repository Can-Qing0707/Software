package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"training_eval_system/internal/model"
	"training_eval_system/internal/repository"
)

type VerifyService struct {
	subRepo     *repository.SubmissionRepo
	taskRepo    *repository.TaskRepo
	evalRepo    *repository.EvalRepo
	llmService  *LLMService
	fileService *FileService
}

func NewVerifyService(subRepo *repository.SubmissionRepo, taskRepo *repository.TaskRepo, evalRepo *repository.EvalRepo, llmService *LLMService, fileService *FileService) *VerifyService {
	return &VerifyService{subRepo: subRepo, taskRepo: taskRepo, evalRepo: evalRepo, llmService: llmService, fileService: fileService}
}

type llmVerificationResult struct {
	Completeness     completenessResult `json:"completeness"`
	LogicIssues      logicResult        `json:"logic_issues"`
	RequirementMatch requirementResult  `json:"requirement_match"`
}

type completenessResult struct {
	StepsTotal        int      `json:"steps_total"`
	StepsCompleted    int      `json:"steps_completed"`
	MissingSteps      []string `json:"missing_steps"`
	CompletenessRatio float64  `json:"completeness_ratio"`
}

type logicResult struct {
	HasLogicIssues bool     `json:"has_logic_issues"`
	Issues         []string `json:"issues"`
	Summary        string   `json:"summary"`
}

type requirementResult struct {
	MatchedRequirements   []string `json:"matched_requirements"`
	UnmatchedRequirements []string `json:"unmatched_requirements"`
	MatchRatio            float64  `json:"match_ratio"`
}

func (s *VerifyService) Verify(submissionID uint) (*model.VerificationResult, error) {
	sub, err := s.subRepo.FindByID(submissionID)
	if err != nil {
		return nil, errors.New("提交记录不存在")
	}

	if sub.Status == "parsing" {
		return nil, errors.New("该提交正在核查中，请稍后再试")
	}

	task, err := s.taskRepo.FindByID(sub.TaskID)
	if err != nil {
		return nil, errors.New("关联任务不存在")
	}

	if !s.subRepo.TryClaimForParsing(submissionID) {
		return nil, errors.New("该提交正在核查中或已完成核查，请稍后再试")
	}

	contentText := sub.ContentText
	if contentText == "" {
		contentText = s.fileService.ExtractContent(sub.FilesJSON)
		if contentText == "" {
			contentText = "学生未提交文本内容，仅有文件附件。"
		} else {
			_ = s.subRepo.Update(submissionID, map[string]interface{}{
				"content_text": contentText,
				"status":       "parsed",
			})
		}
	}

	requirements := task.Description
	if requirements == "" {
		requirements = task.Title
	}

	llmResponse, err := s.llmService.AnalyzeContent(contentText, requirements)
	if err != nil {
		s.subRepo.Update(submissionID, map[string]interface{}{"status": sub.Status})
		return nil, fmt.Errorf("LLM分析失败: %w", err)
	}

	parsed, err := parseLLMVerificationResponse(llmResponse)
	if err != nil {
		_ = s.subRepo.Update(submissionID, map[string]interface{}{"status": sub.Status})
		return nil, fmt.Errorf("解析LLM返回结果失败: %w", err)
	}

	completenessJSON, _ := json.Marshal(parsed.Completeness)
	logicJSON, _ := json.Marshal(parsed.LogicIssues)
	requirementJSON, _ := json.Marshal(parsed.RequirementMatch)

	overallPass := 0
	if parsed.Completeness.CompletenessRatio >= 0.7 &&
		parsed.RequirementMatch.MatchRatio >= 0.7 &&
		!parsed.LogicIssues.HasLogicIssues {
		overallPass = 1
	}

	now := time.Now()
	vr := &model.VerificationResult{
		SubmissionID:     submissionID,
		Completeness:     string(completenessJSON),
		LogicIssues:      string(logicJSON),
		RequirementMatch: string(requirementJSON),
		OverallPass:      &overallPass,
		RawLLMResponse:   llmResponse,
		VerifiedAt:       &now,
	}

	if err := s.evalRepo.UpsertVerification(vr); err != nil {
		_ = s.subRepo.Update(submissionID, map[string]interface{}{"status": sub.Status})
		return nil, fmt.Errorf("保存核查结果失败: %w", err)
	}

	_ = s.subRepo.Update(submissionID, map[string]interface{}{"status": "verified"})
	return vr, nil
}

func (s *VerifyService) GetVerification(submissionID uint) (*model.VerificationResult, error) {
	vr, err := s.evalRepo.GetVerificationBySubmission(submissionID)
	if err != nil {
		return nil, nil
	}
	return vr, nil
}

func parseLLMVerificationResponse(raw string) (*llmVerificationResult, error) {
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

	var result llmVerificationResult
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		start := strings.Index(text, "{")
		end := strings.LastIndex(text, "}")
		if start != -1 && end > start {
			if err2 := json.Unmarshal([]byte(text[start:end+1]), &result); err2 != nil {
				return nil, err
			}
			return &result, nil
		}
		return nil, err
	}
	return &result, nil
}
