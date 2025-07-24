package chat

import (
	"ChatApp/internal/model"
	"context"

	"github.com/gorilla/websocket"
)

type ChatRepository interface {
	AddClient(ctx context.Context, userID string, conn *websocket.Conn)
	StorePendingMessage(receiverID string, message []byte) error
	DeliverPendingMessage(receiverID string)
	RemoveClient(ctx context.Context, userID string)
	SendMessage(ctx context.Context, receiverID string, message []byte) error
	SaveMessageToDB(message model.Chat) error
	GetAllMessages(receiverID string) ([]model.Chat, error)
}