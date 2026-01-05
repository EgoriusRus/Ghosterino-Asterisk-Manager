package domain

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

// UserRole тип роли пользователя
type UserRole string

const (
	UserRoleAdmin UserRole = "admin"
	UserRoleUser  UserRole = "user"
)

// User представляет пользователя системы
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"uniqueIndex;not null" json:"username"`
	PasswordHash string    `gorm:"not null" json:"-"`
	Role         UserRole  `gorm:"default:user" json:"role"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// TableName указывает имя таблицы в БД
func (User) TableName() string {
	return "sipadmin.users"
}

// SetPassword устанавливает хеш пароля
func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

// CheckPassword проверяет пароль
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// UserResponse DTO для ответа (без пароля)
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Role      UserRole  `json:"role"`
	CreatedAt time.Time `json:"createdAt"`
}

// ToResponse конвертирует User в UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
	}
}
