package domain

import "time"

// Location представляет локацию (здание/адрес) с настройками сети
type Location struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Server    string    `gorm:"type:inet;not null" json:"server"`
	Subnet    string    `gorm:"type:cidr;not null" json:"subnet"`
	VoipVLAN  int       `gorm:"not null" json:"voipVlan"`
	VLAN      int       `gorm:"not null" json:"vlan"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// TableName указывает имя таблицы в БД
func (Location) TableName() string {
	return "sipadmin.locations"
}
