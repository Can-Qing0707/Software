package service

import (
	"training_eval_system/internal/dto/request"
	"training_eval_system/internal/model"
	"training_eval_system/internal/repository"
	"training_eval_system/pkg/utils"
)

type UserService struct {
	userRepo *repository.UserRepo
}

func NewUserService(userRepo *repository.UserRepo) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetProfile(userID uint) (*model.User, error) {
	return s.userRepo.FindByID(userID)
}

func (s *UserService) List(req *request.UserQuery) ([]model.User, error) {
	return s.userRepo.List(req.Role, req.Keyword)
}

func (s *UserService) Create(req *request.CreateUserReq) error {
	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}
	user := &model.User{
		Username: req.Username,
		Password: hash,
		RealName: req.RealName,
		Role:     req.Role,
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   1,
	}
	return s.userRepo.Create(user)
}

func (s *UserService) Update(id uint, req *request.UpdateUserReq) error {
	data := make(map[string]interface{})
	if req.RealName != "" {
		data["real_name"] = req.RealName
	}
	if req.Password != "" {
		hash, err := utils.HashPassword(req.Password)
		if err != nil {
			return err
		}
		data["password"] = hash
	}
	if req.Role != "" {
		data["role"] = req.Role
	}
	if req.Email != "" {
		data["email"] = req.Email
	}
	if req.Phone != "" {
		data["phone"] = req.Phone
	}
	if req.Status != nil {
		data["status"] = *req.Status
	}
	return s.userRepo.Update(id, data)
}

func (s *UserService) Delete(id uint) error {
	return s.userRepo.Delete(id)
}
