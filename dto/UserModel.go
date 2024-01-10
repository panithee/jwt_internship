package dto

type Users struct {
	ID       uint
	Email    string
	Password string
	Post     []Posts `gorm:"foreignKey:UserId"`
}
