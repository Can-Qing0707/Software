package repository

import (
	"training_eval_system/internal/model"

	"gorm.io/gorm"
)

type CourseRepo struct {
	db *gorm.DB
}

func NewCourseRepo(db *gorm.DB) *CourseRepo {
	return &CourseRepo{db: db}
}

func (r *CourseRepo) FindByID(id uint) (*model.Course, error) {
	var course model.Course
	err := r.db.First(&course, id).Error
	if err != nil {
		return nil, err
	}
	return &course, nil
}

func (r *CourseRepo) FindByCode(code string) (*model.Course, error) {
	var course model.Course
	err := r.db.Where("code = ?", code).First(&course).Error
	if err != nil {
		return nil, err
	}
	return &course, nil
}

func (r *CourseRepo) List(keyword string) ([]model.Course, error) {
	var courses []model.Course
	query := r.db.Model(&model.Course{}).
		Select("courses.*, COALESCE(u.real_name, '') as teacher_name, (SELECT COUNT(*) FROM course_enrollments ce WHERE ce.course_id = courses.id) as student_count").
		Joins("LEFT JOIN users u ON u.id = courses.teacher_id")
	if keyword != "" {
		query = query.Where("courses.name LIKE ?", "%"+keyword+"%")
	}
	err := query.Order("courses.created_at DESC").Scan(&courses).Error
	return courses, err
}

func (r *CourseRepo) ListByStudent(studentID uint) ([]model.Course, error) {
	var courses []model.Course
	err := r.db.Model(&model.Course{}).
		Select("courses.*, COALESCE(u.real_name, '') as teacher_name, (SELECT COUNT(*) FROM course_enrollments ce2 WHERE ce2.course_id = courses.id) as student_count").
		Joins("LEFT JOIN users u ON u.id = courses.teacher_id").
		Joins("INNER JOIN course_enrollments ce ON ce.course_id = courses.id AND ce.student_id = ?", studentID).
		Order("courses.created_at DESC").Scan(&courses).Error
	return courses, err
}

func (r *CourseRepo) ListByTeacher(teacherID uint, keyword string) ([]model.Course, error) {
	var courses []model.Course
	query := r.db.Model(&model.Course{}).
		Select("courses.*, COALESCE(u.real_name, '') as teacher_name, (SELECT COUNT(*) FROM course_enrollments ce WHERE ce.course_id = courses.id) as student_count").
		Joins("LEFT JOIN users u ON u.id = courses.teacher_id").
		Where("courses.teacher_id = ?", teacherID)
	if keyword != "" {
		query = query.Where("courses.name LIKE ?", "%"+keyword+"%")
	}
	err := query.Order("courses.created_at DESC").Scan(&courses).Error
	return courses, err
}

func (r *CourseRepo) Create(course *model.Course) error {
	return r.db.Create(course).Error
}

func (r *CourseRepo) Update(id uint, data map[string]interface{}) error {
	return r.db.Model(&model.Course{}).Where("id = ?", id).Updates(data).Error
}

func (r *CourseRepo) Delete(id uint) error {
	return r.db.Delete(&model.Course{}, id).Error
}

func (r *CourseRepo) EnrollStudent(courseID, studentID uint) error {
	enrollment := &model.CourseEnrollment{
		CourseID:  courseID,
		StudentID: studentID,
	}
	return r.db.Create(enrollment).Error
}

func (r *CourseRepo) IsEnrolled(courseID, studentID uint) bool {
	var count int64
	r.db.Model(&model.CourseEnrollment{}).
		Where("course_id = ? AND student_id = ?", courseID, studentID).
		Count(&count)
	return count > 0
}

func (r *CourseRepo) UnenrollStudent(courseID, studentID uint) error {
	return r.db.Where("course_id = ? AND student_id = ?", courseID, studentID).
		Delete(&model.CourseEnrollment{}).Error
}
