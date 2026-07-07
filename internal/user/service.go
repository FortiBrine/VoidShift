package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/FortiBrine/VoidShift/internal/config"
	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 12

var ErrNotFound = errors.New("user not found")

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Load(ctx context.Context, cfg config.Config) error {
	return s.loadAdminUser(ctx, cfg)
}

func (s *Service) loadAdminUser(ctx context.Context, cfg config.Config) error {
	if cfg.AdminUsername == "" {
		return nil
	}

	var passwordHash string
	switch {
	case cfg.AdminPasswordHash != "":
		if _, err := bcrypt.Cost([]byte(cfg.AdminPasswordHash)); err != nil {
			return fmt.Errorf("invalid ADMIN_PASSWORD_HASH: %w", err)
		}
		passwordHash = cfg.AdminPasswordHash
	case cfg.AdminPassword != "":
		hashed, err := bcrypt.GenerateFromPassword([]byte(cfg.AdminPassword), bcryptCost)
		if err != nil {
			return fmt.Errorf("hashing admin password: %w", err)
		}
		passwordHash = string(hashed)
	default:
		return nil
	}

	if err := s.CreateUser(ctx, &User{
		Username:     cfg.AdminUsername,
		PasswordHash: passwordHash,
		Admin:        true,
	}); err != nil {
		return fmt.Errorf("creating admin user: %w", err)
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
