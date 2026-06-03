package repository

import (
	"training_eval_system/internal/model"

	"gorm.io/gorm"
)

type EvalRepo struct {
	db *gorm.DB
}

func NewEvalRepo(db *gorm.DB) *EvalRepo {
	return &EvalRepo{db: db}
}

func (r *EvalRepo) ListIndicators() ([]model.EvalIndicator, error) {
	var list []model.EvalIndicator
	err := r.db.Where("status = 1").Order("sort_order ASC").Find(&list).Error
	return list, err
}

func (r *EvalRepo) GetAllIndicators() ([]model.EvalIndicator, error) {
	var list []model.EvalIndicator
	err := r.db.Order("sort_order ASC").Find(&list).Error
	return list, err
}

func (r *EvalRepo) FindIndicatorByID(id uint) (*model.EvalIndicator, error) {
	var ind model.EvalIndicator
	err := r.db.First(&ind, id).Error
	if err != nil {
		return nil, err
	}
	return &ind, nil
}

func (r *EvalRepo) CreateIndicator(ind *model.EvalIndicator) error {
	return r.db.Create(ind).Error
}

func (r *EvalRepo) UpdateIndicator(id uint, data map[string]interface{}) error {
	return r.db.Model(&model.EvalIndicator{}).Where("id = ?", id).Updates(data).Error
}

func (r *EvalRepo) DeleteIndicator(id uint) error {
	return r.db.Delete(&model.EvalIndicator{}, id).Error
}

func (r *EvalRepo) GetTaskIndicators(taskID uint) ([]model.TaskIndicator, error) {
	var list []model.TaskIndicator
	err := r.db.Model(&model.TaskIndicator{}).
		Select("task_indicators.*, COALESCE(ei.name, '') as indicator_name").
		Joins("LEFT JOIN eval_indicators ei ON ei.id = task_indicators.indicator_id").
		Where("task_indicators.task_id = ?", taskID).
		Find(&list).Error
	return list, err
}

func (r *EvalRepo) SaveTaskIndicators(taskID uint, items []model.TaskIndicator) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("task_id = ?", taskID).Delete(&model.TaskIndicator{}).Error; err != nil {
			return err
		}
		if len(items) > 0 {
			if err := tx.Create(&items).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *EvalRepo) GetCourseIndicators(courseID uint) ([]model.CourseIndicator, error) {
	var list []model.CourseIndicator
	err := r.db.Model(&model.CourseIndicator{}).
		Select("course_indicators.*, COALESCE(ei.name, '') as indicator_name").
		Joins("LEFT JOIN eval_indicators ei ON ei.id = course_indicators.indicator_id").
		Where("course_indicators.course_id = ?", courseID).
		Find(&list).Error
	return list, err
}

func (r *EvalRepo) SaveCourseIndicators(courseID uint, items []model.CourseIndicator) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("course_id = ?", courseID).Delete(&model.CourseIndicator{}).Error; err != nil {
			return err
		}
		if len(items) > 0 {
			if err := tx.Create(&items).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *EvalRepo) GetScores(submissionID uint) ([]model.EvalScore, error) {
	var scores []model.EvalScore
	err := r.db.Model(&model.EvalScore{}).
		Select("eval_scores.*, COALESCE(ei.name, '') as indicator_name").
		Joins("LEFT JOIN eval_indicators ei ON ei.id = eval_scores.indicator_id").
		Where("eval_scores.submission_id = ?", submissionID).
		Find(&scores).Error
	return scores, err
}

func (r *EvalRepo) UpsertScore(score *model.EvalScore) error {
	return r.db.Where("submission_id = ? AND indicator_id = ?", score.SubmissionID, score.IndicatorID).
		Assign(score).FirstOrCreate(score).Error
}

func (r *EvalRepo) SaveTeacherScores(submissionID uint, scores []model.EvalScore) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, s := range scores {
			err := tx.Where("submission_id = ? AND indicator_id = ?", submissionID, s.IndicatorID).
				Assign(map[string]interface{}{
					"teacher_score":   s.TeacherScore,
					"teacher_comment": s.TeacherComment,
				}).FirstOrCreate(&model.EvalScore{}).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *EvalRepo) GetVerificationBySubmission(submissionID uint) (*model.VerificationResult, error) {
	var vr model.VerificationResult
	err := r.db.Where("submission_id = ?", submissionID).First(&vr).Error
	if err != nil {
		return nil, err
	}
	return &vr, nil
}

func (r *EvalRepo) CreateVerification(vr *model.VerificationResult) error {
	return r.db.Create(vr).Error
}

func (r *EvalRepo) UpsertVerification(vr *model.VerificationResult) error {
	return r.db.Where("submission_id = ?", vr.SubmissionID).
		Assign(map[string]interface{}{
			"completeness":      vr.Completeness,
			"logic_issues":      vr.LogicIssues,
			"requirement_match": vr.RequirementMatch,
			"overall_pass":      vr.OverallPass,
			"raw_llm_response":  vr.RawLLMResponse,
			"verified_at":       vr.VerifiedAt,
		}).FirstOrCreate(vr).Error
}

func (r *EvalRepo) UpdateScoreFinal(scoreID uint, final float64) error {
	return r.db.Model(&model.EvalScore{}).Where("id = ?", scoreID).
		Update("final_score", final).Error
}
