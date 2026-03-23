package user

import (
	"context"
	"fmt"

	"github.com/FortiBrine/VoidShift/internal/config"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Load(ctx context.Context, cfg config.Config) error {
	err := s.repository.Migrate()
	if err != nil {
		return err
	}

	if cfg.AdminUsername != "" && cfg.AdminPassword != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(cfg.AdminPassword), bcrypt.DefaultCost)

		if err != nil {
			return fmt.Errorf("validator while hashing admin password: %w", err)
		}

		err = s.CreateUser(ctx, &User{
			Username:     cfg.AdminUsername,
			PasswordHash: string(hashed),
			Admin:        true,
		})

		if err != nil {
			return fmt.Errorf("validator while creating admin user: %w", err)
		}
	}

	return nil
}

func (s *Service) CreateUser(ctx context.Context, user *User) error {
	return s.repository.CreateUser(ctx, user)
}

func (s *Service) GetByID(ctx context.Context, id uint) (*User, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *Service) GetByUsername(ctx context.Context, username string) (*User, error) {
	return s.repository.GetByUsername(ctx, username)
}
