package service

import (
	"errors"
	"log"

	"gorm.io/gorm"

	"training_eval_system/config"
	"training_eval_system/internal/dto/request"
	"training_eval_system/internal/dto/response"
	"training_eval_system/internal/model"
	"training_eval_system/internal/repository"
	"training_eval_system/pkg/jwt"
	"training_eval_system/pkg/utils"
)

type AuthService struct {
	userRepo *repository.UserRepo
}

func NewAuthService(userRepo *repository.UserRepo) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Login(req *request.LoginReq) (*response.LoginResp, error) {
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		log.Printf("Login FindByUsername error: %v", err)
		return nil, errors.New("系统繁忙，请稍后再试")
	}
	if user.Status != 1 {
		return nil, errors.New("账号已被禁用")
	}
	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("用户名或密码错误")
	}
	token, err := jwt.GenerateToken(
		config.AppConfig.JWT.Secret,
		config.AppConfig.JWT.ExpireHours,
		user.ID, user.Username, user.Role,
	)
	if err != nil {
		log.Printf("Login GenerateToken error: %v", err)
		return nil, errors.New("系统繁忙，请稍后再试")
	}
	return &response.LoginResp{
		Token: token,
		User: map[string]interface{}{
			"id":        user.ID,
			"username":  user.Username,
			"real_name": user.RealName,
			"role":      user.Role,
			"email":     user.Email,
			"phone":     user.Phone,
		},
	}, nil
}

func (s *AuthService) Register(req *request.RegisterReq) error {
	existing, err := s.userRepo.FindByUsername(req.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Printf("Register FindByUsername error: %v", err)
		return errors.New("系统繁忙，请稍后再试")
	}
	if existing != nil {
		return errors.New("用户名已存在")
	}
	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Printf("Register HashPassword error: %v", err)
		return errors.New("系统繁忙，请稍后再试")
	}
	user := &model.User{
		Username: req.Username,
		Password: hash,
		RealName: req.RealName,
		Role:     req.Role,
		Status:   1,
	}
	if err := s.userRepo.Create(user); err != nil {
		log.Printf("Register Create error: %v", err)
		return errors.New("系统繁忙，请稍后再试")
	}
	return nil
}
