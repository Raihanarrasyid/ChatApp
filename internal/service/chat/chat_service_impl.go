package chat

import (
	"ChatApp/internal/model"
	"ChatApp/internal/repository/chat"
	"context"
	"encoding/json"
	"time"
)

type ChatServiceImpl struct {
	ChatRepository chat.ChatRepository
	ctx context.Context
}

func NewChatService(chatRepository chat.ChatRepository) *ChatServiceImpl {
	return &ChatServiceImpl{ChatRepository: chatRepository, ctx: context.Background()}
}

func (cs *ChatServiceImpl) HandleDirectMessage(senderID, receiverID, content string) error {
	message := model.Message{
		Sender: senderID,
		Content: content,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	messageBytes, _ := json.Marshal(message)

	return cs.ChatRepository.SendMessage(cs.ctx, receiverID, messageBytes)
}