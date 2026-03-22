package session

import "time"

type Session struct {
	ID         uint   `gorm:"primary_key"`
	SessionID  string `gorm:"uniqueIndex:idx_session_id"`
	UserID     uint
	UserAgent  string
	IP         string
	ExpiresAt  time.Time
	CreatedAt  time.Time
	LastSeenAt time.Time
}
