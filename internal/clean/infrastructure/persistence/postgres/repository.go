//go:build modern

package postgres

import (
	"mikhmon_v4/internal/clean/domain/entity"
	"mikhmon_v4/internal/clean/domain/repository"

	"gorm.io/gorm"
)

type UserRepo struct{ db *gorm.DB }

type RouterRepo struct{ db *gorm.DB }

func NewUserRepo(db *gorm.DB) repository.UserRepository { return &UserRepo{db: db} }

func NewRouterRepo(db *gorm.DB) repository.RouterRepository { return &RouterRepo{db: db} }

func (r *UserRepo) FindByUsername(username string) (*entity.User, error) {
	var user entity.User
	if err := r.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *RouterRepo) List() ([]entity.Router, error) {
	var routers []entity.Router
	if err := r.db.Order("name asc").Find(&routers).Error; err != nil {
		return nil, err
	}
	return routers, nil
}

func (r *RouterRepo) Create(router *entity.Router) error {
	return r.db.Create(router).Error
}
