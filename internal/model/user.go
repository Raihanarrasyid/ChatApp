package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Username string `gorm:"type:varchar(255);not null"`
	Password string `gorm:"type:varchar(255);not null"`
	Email string `gorm:"type:varchar(255);not null"`
	Chats []Chat `gorm:"foreignKey:FromID"`
}