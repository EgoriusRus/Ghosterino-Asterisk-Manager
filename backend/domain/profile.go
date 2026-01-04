package domain

import (
	"time"
)

// Profile представляет профиль сотрудника с SIP-настройками
type Profile struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	Name           string    `json:"name"`
	Email          string    `json:"email"`
	Device         *string   `gorm:"type:macaddr" json:"device"`
	LocationID     *uint     `json:"locationId"`
	InternalNumber int       `gorm:"uniqueIndex;not null" json:"internalNumber"`
	ExternalNumber string    `json:"externalNumber"`
	RingGroup      *int      `json:"ringGroup"`
	PickupGroup    *int      `json:"pickupGroup"`
	IsActive       bool      `gorm:"default:true" json:"isActive"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`

	// Связи
	Location *Location `gorm:"foreignKey:LocationID" json:"location,omitempty"`
}

// TableName указывает имя таблицы в БД
func (Profile) TableName() string {
	return "sipadmin.profiles"
}

// IsLocalOnly проверяет, только локальные звонки
func (p *Profile) IsLocalOnly() bool {
	return p.RingGroup != nil && *p.RingGroup == 0
}

// GetExternalNumberClean возвращает внешний номер без дефисов
func (p *Profile) GetExternalNumberClean() string {
	return cleanPhoneNumber(p.ExternalNumber)
}

func cleanPhoneNumber(phone string) string {
	result := ""
	for _, char := range phone {
		if char >= '0' && char <= '9' {
			result += string(char)
		}
	}
	return result
}
