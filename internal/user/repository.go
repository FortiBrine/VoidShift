package user

import (
	"context"
	"database/sql"

	"github.com/FortiBrine/VoidShift/internal/store"
)

type Repository interface {
	CreateUser(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uint) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
}

type SqlcRepository struct {
	q *store.Queries
}

func NewSqlcRepository(database *sql.DB) *SqlcRepository {
	return &SqlcRepository{q: store.New(database)}
}

func (r *SqlcRepository) CreateUser(ctx context.Context, user *User) error {
	admin := int64(0)
	if user.Admin {
		admin = 1
	}
	return r.q.CreateUser(ctx, store.CreateUserParams{
		Username:     user.Username,
		PasswordHash: user.PasswordHash,
		Admin:        admin,
	})
}

func (r *SqlcRepository) GetByID(ctx context.Context, id uint) (*User, error) {
	row, err := r.q.GetUserByID(ctx, int64(id))
	if err != nil {
		return nil, err
	}
	return fromDB(row), nil
}

func (r *SqlcRepository) GetByUsername(ctx context.Context, username string) (*User, error) {
	row, err := r.q.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return fromDB(row), nil
}

func fromDB(row store.User) *User {
	return &User{
		ID:           uint(row.ID),
		Username:     row.Username,
		PasswordHash: row.PasswordHash,
		Admin:        row.Admin != 0,
	}
}
