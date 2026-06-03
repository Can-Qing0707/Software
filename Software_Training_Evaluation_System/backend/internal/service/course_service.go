package service

import (
	"crypto/rand"
	"errors"
	"math/big"

	"gorm.io/gorm"

	"training_eval_system/internal/dto/request"
	"training_eval_system/internal/model"
	"training_eval_system/internal/repository"
)

type CourseService struct {
	courseRepo *repository.CourseRepo
}

func NewCourseService(courseRepo *repository.CourseRepo) *CourseService {
	return &CourseService{courseRepo: courseRepo}
}

const codeChars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
const codeLength = 6

func generateCourseCode() string {
	code := make([]byte, codeLength)
	for i := range code {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(codeChars))))
		code[i] = codeChars[n.Int64()]
	}
	return string(code)
}

func (s *CourseService) List(req *request.CourseQuery) ([]model.Course, error) {
	return s.courseRepo.List(req.Keyword)
}

func (s *CourseService) ListByTeacher(teacherID uint, req *request.CourseQuery) ([]model.Course, error) {
	return s.courseRepo.ListByTeacher(teacherID, req.Keyword)
}

func (s *CourseService) GetByID(id uint) (*model.Course, error) {
	return s.courseRepo.FindByID(id)
}

func (s *CourseService) Create(req *request.CreateCourseReq, teacherID uint) error {
	code := generateCourseCode()
	course := &model.Course{
		Name:        req.Name,
		Description: req.Description,
		TeacherID:   teacherID,
		Semester:    req.Semester,
		Code:        code,
		Status:      1,
	}
	return s.courseRepo.Create(course)
}

func (s *CourseService) Update(id uint, req *request.UpdateCourseReq) error {
	data := make(map[string]interface{})
	if req.Name != "" {
		data["name"] = req.Name
	}
	if req.Description != "" {
		data["description"] = req.Description
	}
	if req.Semester != "" {
		data["semester"] = req.Semester
	}
	if req.Status != nil {
		data["status"] = *req.Status
	}
	return s.courseRepo.Update(id, data)
}

func (s *CourseService) Delete(id uint) error {
	return s.courseRepo.Delete(id)
}

func (s *CourseService) JoinByCode(code string, studentID uint) error {
	course, err := s.courseRepo.FindByCode(code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("课程代码无效")
		}
		return err
	}
	if s.courseRepo.IsEnrolled(course.ID, studentID) {
		return errors.New("已加入该课程")
	}
	return s.courseRepo.EnrollStudent(course.ID, studentID)
}

func (s *CourseService) GetMyCourses(studentID uint) ([]model.Course, error) {
	return s.courseRepo.ListByStudent(studentID)
}

func (s *CourseService) LeaveCourse(courseID, studentID uint) error {
	if !s.courseRepo.IsEnrolled(courseID, studentID) {
		return errors.New("未加入该课程")
	}
	return s.courseRepo.UnenrollStudent(courseID, studentID)
}
