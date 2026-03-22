package session

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"gorm.io/gorm"
)

var ErrInvalidSession = errors.New("invalid session")

type Service struct {
	repository      Repository
	sessionDuration time.Duration
}

func NewService(repository Repository, sessionDuration time.Duration) *Service {
	return &Service{repository: repository, sessionDuration: sessionDuration}
}

func (s *Service) Load() error {
	return s.repository.Migrate()
}

func (s *Service) CreateUserSession(
	ctx context.Context,
	userID uint,
	userAgent, ip string,
) (string, time.Time, error) {
	sessionID, err := NewSessionID()
	if err != nil {
		return "", time.Time{}, err
	}

	expiresAt := time.Now().Add(s.sessionDuration)

	err = s.repository.CreateSession(ctx, &Session{
		SessionID: sessionID,
		UserID:    userID,
		UserAgent: userAgent,
		IP:        ip,
		ExpiresAt: expiresAt,
	})

	if err != nil {
		return "", time.Time{}, err
	}

	return sessionID, expiresAt, nil
}

func (s *Service) ValidateSession(
	ctx context.Context,
	sessionID string,
) (uint, error) {
	session, err := s.repository.GetSessionByID(ctx, sessionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, ErrInvalidSession
		}

		return 0, err
	}

	if time.Now().After(session.ExpiresAt) {
		_ = s.repository.DeleteSession(ctx, sessionID)
		return 0, ErrInvalidSession
	}

	_ = s.repository.TouchSession(ctx, sessionID, time.Now())
	return session.UserID, nil
}

func (s *Service) LogoutSession(ctx context.Context, sessionID string) error {
	return s.repository.DeleteSession(ctx, sessionID)
}

func NewSessionID() (string, error) {
	b := make([]byte, 32) // 256 bits
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
