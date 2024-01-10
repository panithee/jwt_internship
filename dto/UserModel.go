package dto

type Users struct {
	ID       uint
	Email    string `gorm:"unique;not null"`
	Password string
	Post     []Posts `gorm:"foreignKey:UserId"`
}
