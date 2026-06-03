package service

import (
	"training_eval_system/internal/dto/request"
	"training_eval_system/internal/model"
	"training_eval_system/internal/repository"
)

type TaskService struct {
	taskRepo *repository.TaskRepo
}

func NewTaskService(taskRepo *repository.TaskRepo) *TaskService {
	return &TaskService{taskRepo: taskRepo}
}

func (s *TaskService) List(req *request.TaskQuery) ([]model.Task, error) {
	return s.taskRepo.List(req.CourseID, req.Keyword)
}

func (s *TaskService) ListByTeacher(teacherID uint, req *request.TaskQuery) ([]model.Task, error) {
	return s.taskRepo.ListByTeacher(teacherID, req.CourseID, req.Keyword)
}

func (s *TaskService) ListByStudent(studentID uint, req *request.TaskQuery) ([]model.Task, error) {
	return s.taskRepo.ListByStudent(studentID, req.CourseID, req.Keyword)
}

func (s *TaskService) GetByID(id uint) (*model.Task, error) {
	return s.taskRepo.FindByID(id)
}

func (s *TaskService) Create(req *request.CreateTaskReq) error {
	task := &model.Task{
		CourseID:    req.CourseID,
		Title:       req.Title,
		Description: req.Description,
		Deadline:    req.Deadline,
		Status:      1,
	}
	return s.taskRepo.Create(task)
}

func (s *TaskService) Update(id uint, req *request.UpdateTaskReq) error {
	data := make(map[string]interface{})
	if req.Title != "" {
		data["title"] = req.Title
	}
	if req.Description != "" {
		data["description"] = req.Description
	}
	if req.Deadline != nil {
		data["deadline"] = req.Deadline
	}
	if req.Status != nil {
		data["status"] = *req.Status
	}
	return s.taskRepo.Update(id, data)
}

func (s *TaskService) Delete(id uint) error {
	return s.taskRepo.Delete(id)
}
