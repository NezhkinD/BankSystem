package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"unique" db:"username" json:"username"`
	Email    string `gorm:"unique" db:"email" json:"email"`
	Password string `db:"password" json:"password"`
}
