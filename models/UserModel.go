package models

type Users struct {
	ID       uint    `gorm:"primaryKey"`
	Email    string  `gorm:"unique;not null"`
	Password string  `gorm:"not null"`
	Post     []Posts `gorm:"foreignKey:UserId"`
}
