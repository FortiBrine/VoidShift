package user

type User struct {
	ID           uint   `gorm:"primary_key"`
	Username     string `gorm:"uniqueIndex:idx_username"`
	PasswordHash string
	Admin        bool
}
