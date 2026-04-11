package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/FortiBrine/VoidShift/internal/session"
	"github.com/FortiBrine/VoidShift/internal/shared"
	"github.com/FortiBrine/VoidShift/internal/user"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	sessionService *session.Service
	userService    *user.Service
}

type LoginResult struct {
	SessionID string
	ExpiresAt time.Time
}

func (s *Service) Login(
	ctx context.Context,
	username, password, userAgent, ip string,
) (*LoginResult, error) {
	u, err := s.userService.GetByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, shared.ErrInvalidCredentials
		}
		return nil, fmt.Errorf("get user by username: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, shared.ErrInvalidCredentials
	}

	sessionID, expiresAt, err := s.sessionService.CreateUserSession(ctx, u.ID, userAgent, ip)
	if err != nil {
		return nil, fmt.Errorf("create user session: %w", err)
	}

	return &LoginResult{
		SessionID: sessionID,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *Service) LogoutSession(ctx context.Context, sessionID string) error {
	if err := s.sessionService.LogoutSession(ctx, sessionID); err != nil {
		return fmt.Errorf("logout session: %w", err)
	}
	return nil
}
