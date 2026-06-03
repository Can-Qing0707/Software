package repository

import (
	"training_eval_system/internal/model"

	"gorm.io/gorm"
)

type ReportRepo struct {
	db *gorm.DB
}

func NewReportRepo(db *gorm.DB) *ReportRepo {
	return &ReportRepo{db: db}
}

func (r *ReportRepo) List(courseID, taskID uint) ([]model.Report, error) {
	var list []model.Report
	query := r.db.Model(&model.Report{})
	if courseID > 0 {
		query = query.Where("course_id = ?", courseID)
	}
	if taskID > 0 {
		query = query.Where("task_id = ?", taskID)
	}
	err := query.Order("generated_at DESC").Find(&list).Error
	return list, err
}

func (r *ReportRepo) Create(report *model.Report) error {
	return r.db.Create(report).Error
}

func (r *ReportRepo) FindByID(id uint) (*model.Report, error) {
	var report model.Report
	err := r.db.First(&report, id).Error
	if err != nil {
		return nil, err
	}
	return &report, nil
}
