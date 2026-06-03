package repository

import (
	"training_eval_system/internal/model"

	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) List(role, keyword string) ([]model.User, error) {
	var users []model.User
	query := r.db.Model(&model.User{})
	if role != "" {
		query = query.Where("role = ?", role)
	}
	if keyword != "" {
		query = query.Where("username LIKE ? OR real_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	err := query.Order("created_at DESC").Find(&users).Error
	return users, err
}

func (r *UserRepo) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepo) Update(id uint, data map[string]interface{}) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Updates(data).Error
}

func (r *UserRepo) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}
