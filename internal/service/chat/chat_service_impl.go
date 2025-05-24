package chat

import (
	"ChatApp/internal/http/response"
	"ChatApp/internal/model"
	"ChatApp/internal/repository/chat"

	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type ChatServiceImpl struct {
	ChatRepository chat.ChatRepository
	ctx context.Context
}

func NewChatService(chatRepository chat.ChatRepository) *ChatServiceImpl {
	return &ChatServiceImpl{ChatRepository: chatRepository, ctx: context.Background()}
}

func (cs *ChatServiceImpl) GetAllMessages(ctx context.Context, receiverID string) ([]response.ChatResponse, error) {
	chats, err := cs.ChatRepository.GetAllMessages(receiverID)
	if err != nil {
		return nil, err
	}

	var chatResponses []response.ChatResponse
	for _, chat := range chats {
		chatResponses = append(chatResponses, response.ChatResponse{
			ID: chat.ID,
			FromID: chat.FromID,
			ToID: chat.ToID,
			Message: chat.Message,
			IsRead: chat.IsRead,
			FromUser: response.UserResponse{
				ID: chat.FromUser.ID,
				Username: chat.FromUser.Username,
				Email: chat.FromUser.Email,
			},
			ToUser: response.UserResponse{
				ID: chat.ToUser.ID,
				Username: chat.ToUser.Username,
				Email: chat.ToUser.Email,
			},
		})
	}

	return chatResponses, nil
}

func (cs *ChatServiceImpl) HandleDirectMessage(senderID, receiverID, content string) error {

	chat := model.Chat{
		FromID: uuid.MustParse(senderID),
		ToID: uuid.MustParse(receiverID),
		Message: content,
	}

	if err := cs.ChatRepository.SaveMessageToDB(chat); err != nil {
		return err
	}

	message := model.Message{
		Sender: senderID,
		Content: content,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	messageBytes, _ := json.Marshal(message)

	return cs.ChatRepository.SendMessage(cs.ctx, receiverID, messageBytes)
}