package repository

import (
	"training_eval_system/internal/model"

	"gorm.io/gorm"
)

type TaskRepo struct {
	db *gorm.DB
}

func NewTaskRepo(db *gorm.DB) *TaskRepo {
	return &TaskRepo{db: db}
}

func (r *TaskRepo) FindByID(id uint) (*model.Task, error) {
	var task model.Task
	err := r.db.First(&task, id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepo) List(courseID uint, keyword string) ([]model.Task, error) {
	var tasks []model.Task
	query := r.db.Model(&model.Task{}).Select("tasks.*, COALESCE(c.name, '') as course_name").
		Joins("LEFT JOIN courses c ON c.id = tasks.course_id")
	if courseID > 0 {
		query = query.Where("tasks.course_id = ?", courseID)
	}
	if keyword != "" {
		query = query.Where("tasks.title LIKE ?", "%"+keyword+"%")
	}
	err := query.Order("tasks.created_at DESC").Scan(&tasks).Error
	return tasks, err
}

func (r *TaskRepo) Create(task *model.Task) error {
	return r.db.Create(task).Error
}

func (r *TaskRepo) Update(id uint, data map[string]interface{}) error {
	return r.db.Model(&model.Task{}).Where("id = ?", id).Updates(data).Error
}

func (r *TaskRepo) Delete(id uint) error {
	return r.db.Delete(&model.Task{}, id).Error
}

func (r *TaskRepo) ListByCourseID(courseID uint) ([]model.Task, error) {
	var tasks []model.Task
	err := r.db.Where("course_id = ?", courseID).Find(&tasks).Error
	return tasks, err
}

func (r *TaskRepo) ListByTeacher(teacherID, courseID uint, keyword string) ([]model.Task, error) {
	var tasks []model.Task
	query := r.db.Model(&model.Task{}).Select("tasks.*, COALESCE(c.name, '') as course_name").
		Joins("LEFT JOIN courses c ON c.id = tasks.course_id").
		Where("c.teacher_id = ?", teacherID)
	if courseID > 0 {
		query = query.Where("tasks.course_id = ?", courseID)
	}
	if keyword != "" {
		query = query.Where("tasks.title LIKE ?", "%"+keyword+"%")
	}
	err := query.Order("tasks.created_at DESC").Scan(&tasks).Error
	return tasks, err
}

func (r *TaskRepo) ListByStudent(studentID, courseID uint, keyword string) ([]model.Task, error) {
	var tasks []model.Task
	query := r.db.Model(&model.Task{}).Select("tasks.*, COALESCE(c.name, '') as course_name").
		Joins("LEFT JOIN courses c ON c.id = tasks.course_id").
		Joins("INNER JOIN course_enrollments ce ON ce.course_id = tasks.course_id AND ce.student_id = ?", studentID)
	if courseID > 0 {
		query = query.Where("tasks.course_id = ?", courseID)
	}
	if keyword != "" {
		query = query.Where("tasks.title LIKE ?", "%"+keyword+"%")
	}
	err := query.Order("tasks.created_at DESC").Scan(&tasks).Error
	return tasks, err
}
