package session

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	Migrate() error
	CreateSession(ctx context.Context, session *Session) error
	GetSessionByID(ctx context.Context, sessionID string) (*Session, error)
	DeleteSession(ctx context.Context, sessionID string) error
	TouchSession(ctx context.Context, sessionID string, t time.Time) error
}

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{
		db: db,
	}
}

func (r *GormRepository) Migrate() error {
	return r.db.AutoMigrate(&Session{})
}

func (r *GormRepository) CreateSession(ctx context.Context, session *Session) error {
	return gorm.G[Session](r.db).Create(ctx, session)
}

func (r *GormRepository) GetSessionByID(ctx context.Context, sessionID string) (*Session, error) {
	session, err := gorm.G[Session](r.db).Where("session_id = ?", sessionID).First(ctx)
	return &session, err
}

func (r *GormRepository) DeleteSession(ctx context.Context, sessionID string) error {
	_, err := gorm.G[Session](r.db).Where("session_id = ?", sessionID).Delete(ctx)
	return err
}

func (r *GormRepository) TouchSession(ctx context.Context, sessionID string, t time.Time) error {
	_, err := gorm.G[Session](r.db).Where("session_id = ?", sessionID).Update(ctx, "last_seen_at", t)
	return err
}
