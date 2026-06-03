package repository

import (
	"training_eval_system/internal/model"

	"gorm.io/gorm"
)

type SubmissionRepo struct {
	db *gorm.DB
}

func NewSubmissionRepo(db *gorm.DB) *SubmissionRepo {
	return &SubmissionRepo{db: db}
}

const submissionListSelect = `submissions.*, COALESCE(u.real_name, '') as student_name, COALESCE(t.title, '') as task_title,
(SELECT ROUND(SUM(COALESCE(es.llm_score * ti.weight / 100, 0)), 1)
 FROM eval_scores es
 LEFT JOIN task_indicators ti ON ti.indicator_id = es.indicator_id AND ti.task_id = submissions.task_id
 WHERE es.submission_id = submissions.id AND es.llm_score IS NOT NULL) as llm_total,
(SELECT ROUND(SUM(COALESCE(es.teacher_score * ti.weight / 100, 0)), 1)
 FROM eval_scores es
 LEFT JOIN task_indicators ti ON ti.indicator_id = es.indicator_id AND ti.task_id = submissions.task_id
 WHERE es.submission_id = submissions.id AND es.teacher_score IS NOT NULL) as teacher_total,
(SELECT ROUND(SUM(COALESCE(es.final_score, 0)), 1)
 FROM eval_scores es WHERE es.submission_id = submissions.id) as total_score`

func (r *SubmissionRepo) FindByID(id uint) (*model.Submission, error) {
	var sub model.Submission
	err := r.db.First(&sub, id).Error
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func (r *SubmissionRepo) List(taskID uint, status string) ([]model.Submission, error) {
	var list []model.Submission
	query := r.db.Model(&model.Submission{}).
		Select(submissionListSelect).
		Joins("LEFT JOIN users u ON u.id = submissions.student_id").
		Joins("LEFT JOIN tasks t ON t.id = submissions.task_id")
	if taskID > 0 {
		query = query.Where("submissions.task_id = ?", taskID)
	}
	if status != "" {
		query = query.Where("submissions.status IN ?", parseStatusFilter(status))
	}
	err := query.Order("submissions.submit_time DESC").Scan(&list).Error
	return list, err
}

func (r *SubmissionRepo) ListByStudent(studentID, taskID uint, status string) ([]model.Submission, error) {
	var list []model.Submission
	query := r.db.Model(&model.Submission{}).
		Select(submissionListSelect).
		Joins("LEFT JOIN users u ON u.id = submissions.student_id").
		Joins("LEFT JOIN tasks t ON t.id = submissions.task_id").
		Where("submissions.student_id = ?", studentID)
	if taskID > 0 {
		query = query.Where("submissions.task_id = ?", taskID)
	}
	if status != "" {
		query = query.Where("submissions.status IN ?", parseStatusFilter(status))
	}
	err := query.Order("submissions.submit_time DESC").Scan(&list).Error
	return list, err
}

func (r *SubmissionRepo) ListByTeacher(teacherID, taskID uint, status string) ([]model.Submission, error) {
	var list []model.Submission
	query := r.db.Model(&model.Submission{}).
		Select(submissionListSelect).
		Joins("LEFT JOIN users u ON u.id = submissions.student_id").
		Joins("LEFT JOIN tasks t ON t.id = submissions.task_id").
		Joins("INNER JOIN courses c ON c.id = t.course_id AND c.teacher_id = ?", teacherID)
	if taskID > 0 {
		query = query.Where("submissions.task_id = ?", taskID)
	}
	if status != "" {
		query = query.Where("submissions.status IN ?", parseStatusFilter(status))
	}
	err := query.Order("submissions.submit_time DESC").Scan(&list).Error
	return list, err
}

func (r *SubmissionRepo) Create(sub *model.Submission) error {
	return r.db.Create(sub).Error
}

func (r *SubmissionRepo) Update(id uint, data map[string]interface{}) error {
	return r.db.Model(&model.Submission{}).Where("id = ?", id).Updates(data).Error
}

func (r *SubmissionRepo) TryClaimForParsing(id uint) bool {
	result := r.db.Model(&model.Submission{}).
		Where("id = ? AND status IN ?", id, []string{"uploaded", "parsed", "verified", "evaluated"}).
		Update("status", "parsing")
	return result.RowsAffected > 0
}

func (r *SubmissionRepo) GetByTaskAndStudent(taskID, studentID uint) (*model.Submission, error) {
	var sub model.Submission
	err := r.db.Where("task_id = ? AND student_id = ?", taskID, studentID).First(&sub).Error
	if err != nil {
		return nil, err
	}
	return &sub, nil
}

func parseStatusFilter(status string) []string {
	statusMap := map[string][]string{
		"uploaded":  {"uploaded"},
		"parsing":   {"parsing"},
		"parsed":    {"parsed"},
		"verified":  {"verified"},
		"evaluated": {"evaluated"},
	}
	if s, ok := statusMap[status]; ok {
		return s
	}
	return []string{"uploaded", "parsing", "parsed", "verified", "evaluated"}
}
