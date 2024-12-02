package chat

import (
	"ChatApp/internal/model"
	"ChatApp/internal/service/chat"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool{
		return true
	},
}

type ChatControllerImpl struct {
	ChatService *chat.ChatServiceImpl
}

func NewChatController(router *gin.RouterGroup, chatService *chat.ChatServiceImpl) {
	controller := &ChatControllerImpl{
		ChatService: chatService,
	}

	router.GET("/ws", controller.HandleWebsocket)
}

func (cc *ChatControllerImpl) HandleWebsocket(c *gin.Context) {
	userID := c.Query("userID")

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userID is required"})
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
		return
	}
	defer conn.Close()

	cc.ChatService.ChatRepository.AddClient(c, userID, conn)
	defer cc.ChatService.ChatRepository.RemoveClient(c, userID)

	for {
		_, message, err := conn.ReadMessage()

		if err != nil {
			break
		}

		var wsMessage model.WebSocketMessage
		if err := json.Unmarshal(message, &wsMessage); err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("Invalid message format"))
        	continue
		}
		
		if wsMessage.ReceiverID == "" {
			conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("receiverID is required")))
			continue
		}

		err = cc.ChatService.HandleDirectMessage(userID, wsMessage.ReceiverID, string(wsMessage.Content))
		if err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("Failed to send message"))
		}
	}
	
}