package chat

import (
	"context"
	"fmt"
	"sync"

	"ChatApp/internal/model"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type ChatRepositoryImpl struct {
	clients map[string]*websocket.Conn
	mu sync.Mutex
	redis *redis.Client
	ctx context.Context
	db *gorm.DB
}

func NewChatRepository(redisClient *redis.Client, db *gorm.DB) ChatRepository {
	return &ChatRepositoryImpl{
		clients: make(map[string]*websocket.Conn),
		redis: redisClient,
		ctx: context.Background(),
		db: db,
	}
}

func (c *ChatRepositoryImpl) GetAllMessages(receiverID string) ([]model.Chat, error) {
	var messages []model.Chat
	err := c.db.Preload("FromUser").Preload("ToUser").Where("to_id = ?", receiverID).Find(&messages).Error
	return messages, err
}

func (c *ChatRepositoryImpl) SaveMessageToDB(message model.Chat) error {
	newMessage := model.Chat{
		FromID: message.FromID,
		ToID: message.ToID,
		Message: message.Message,
	}

	return c.db.Create(&newMessage).Error
}


func (c *ChatRepositoryImpl) AddClient(ctx context.Context, userID string, conn *websocket.Conn) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.clients[userID] = conn

	go c.DeliverPendingMessage(userID)
}

func (c *ChatRepositoryImpl) StorePendingMessage(receiverID string, message []byte) error {
	return c.redis.LPush(c.ctx, fmt.Sprintf("message:%s", receiverID), message).Err()
}

func (c *ChatRepositoryImpl) DeliverPendingMessage(receiverID string) {
	conn, exist := c.clients[receiverID]
	if !exist {
		return
	}

	key := fmt.Sprintf("message:%s", receiverID)
	for {
		message, err := c.redis.RPop(c.ctx, key).Result()
		if err == redis.Nil {
			break
		} else if err != nil {
			continue
		}

		conn.WriteMessage(websocket.TextMessage, []byte(message))
	}
}

func (c *ChatRepositoryImpl) RemoveClient(ctx context.Context, userID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.clients, userID)
}

func (c *ChatRepositoryImpl) SendMessage(ctx context.Context, receiverID string, message []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	conn, exists := c.clients[receiverID]
	if exists {
		return conn.WriteMessage(websocket.TextMessage, message)
	}

	return c.StorePendingMessage(receiverID, message)
}
