package models

import "time"

type Posts struct {
	ID        uint      `gorm:"primaryKey"`
	Message   string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
	UserId    uint      `gorm:"not null"`
}
