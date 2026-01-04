package domain

import "time"

// DeviceModel типы поддерживаемых устройств
type DeviceModel string

const (
	DeviceModelYealinkT27G DeviceModel = "Yealink T27G"
	DeviceModelYealinkT23G DeviceModel = "Yealink T23G"
	DeviceModelFanvil      DeviceModel = "Fanvil"
	DeviceModelCisco       DeviceModel = "Cisco"
)

// Device представляет IP-телефон или устройство
type Device struct {
	MAC         string      `gorm:"primaryKey;type:macaddr" json:"mac"`
	DeviceModel DeviceModel `gorm:"not null" json:"deviceModel"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
}

// TableName указывает имя таблицы в БД
func (Device) TableName() string {
	return "sipadmin.devices"
}
