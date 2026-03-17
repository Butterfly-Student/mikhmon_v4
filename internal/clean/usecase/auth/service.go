//go:build modern

package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"mikhmon_v4/internal/clean/domain/entity"
	"mikhmon_v4/internal/clean/domain/repository"
)

type Service struct {
	users repository.UserRepository
}

func NewService(users repository.UserRepository) *Service {
	return &Service{users: users}
}

func (s *Service) Register(username, plainPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.users.Create(&entity.User{Username: username, PasswordHash: string(hash), Role: "admin"})
}

func (s *Service) Login(username, plainPassword string) error {
	user, err := s.users.FindByUsername(username)
	if err != nil {
		return err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(plainPassword)) != nil {
		return errors.New("invalid username or password")
	}
	return nil
}
