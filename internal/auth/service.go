package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/FortiBrine/VoidShift/internal/user"
	"golang.org/x/crypto/bcrypt"
)

var ErrUserNotFound = errors.New("user not found")
var ErrWrongPassword = errors.New("wrong password")

type Service struct {
	userService *user.Service
}

func NewService(userService *user.Service) *Service {
	return &Service{
		userService: userService,
	}
}

type LoginResult struct {
	SessionID string
	ExpiresAt time.Time
}

func (s *Service) Login(
	ctx context.Context,
	username, password, userAgent, ip string,
) error {
	u, err := s.userService.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrUserNotFound
		}
		return fmt.Errorf("get user by username: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return ErrWrongPassword
	}

	return nil
}
