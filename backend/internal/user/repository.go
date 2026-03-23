package user

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	Migrate() error
	CreateUser(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uint) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
}

type GormRepository struct {
	db *gorm.DB
}

func NewGormRepository(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Migrate() error {
	return r.db.AutoMigrate(&User{})
}

func (r *GormRepository) CreateUser(ctx context.Context, user *User) error {
	return gorm.G[User](r.db, clause.OnConflict{
		Columns:   []clause.Column{{Name: "username"}},
		DoNothing: true,
	}).Create(ctx, user)
}

func (r *GormRepository) GetByID(ctx context.Context, id uint) (*User, error) {
	user, err := gorm.G[User](r.db).Where("id = ?", id).First(ctx)
	return &user, err
}

func (r *GormRepository) GetByUsername(ctx context.Context, username string) (*User, error) {
	user, err := gorm.G[User](r.db).Where("username = ?", username).First(ctx)
	return &user, err
}
