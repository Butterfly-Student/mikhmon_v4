//go:build modern

package repository

import "mikhmon_v4/internal/clean/domain/entity"

type UserRepository interface {
	FindByUsername(username string) (*entity.User, error)
	Create(user *entity.User) error
}

type RouterRepository interface {
	List() ([]entity.Router, error)
	Create(router *entity.Router) error
}
