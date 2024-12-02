package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Chat struct {
	gorm.Model
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	FromID uuid.UUID `gorm:"type:uuid;not null"`
	ToID uuid.UUID `gorm:"type:uuid;not null"`
	Message string `gorm:"type:text;not null"`
	IsRead bool `gorm:"type:boolean;default:false"`
	FromUser User `gorm:"foreignKey:FromID;references:ID"`
	ToUser User `gorm:"foreignKey:ToID;references:ID"`
}

type Message struct {
	Sender string
	Content string
	Timestamp string
}

type WebSocketMessage struct {
	ReceiverID string `json: "receiverID"`
	Content string `json: "content"`
}