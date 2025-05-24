package chat

import (
	_ "ChatApp/internal/http/response"
	"ChatApp/internal/model"
	"ChatApp/internal/service/chat"
	"encoding/json"
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
	router.GET("/", controller.GetAllChats)
}

// @Summary		Get all chats
// @Description	Get all chats
// @Tags			chat
// @Accept			json
// @Produce		json
// @Success		200	{object}	http.Response{value=[]response.ChatResponse}
// @Failure		400	{object}	http.Error
// @Failure		404	{object}	http.Error
// @Failure		500	{object}	http.Error
// @Router			/chat [get]
func (cc *ChatControllerImpl) GetAllChats(c *gin.Context) {
	chats, err := cc.ChatService.GetAllMessages(c, c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if chats == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chats not found"})
		return
	}
	c.JSON(http.StatusOK, chats)
}

//	@Summary		Handle Websocket
//	@Description	Handle Websocket
//	@Tags			chat
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	http.Response
//	@Failure		400	{object}	http.Error
//	@Failure		404	{object}	http.Error
//	@Failure		500	{object}	http.Error
//	@Router			/chat/ws [get]
func (cc *ChatControllerImpl) HandleWebsocket(c *gin.Context) {
	userID := c.GetString("user_id")

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
			conn.WriteMessage(websocket.TextMessage, []byte("receiverID is required"))
			continue
		}

		err = cc.ChatService.HandleDirectMessage(userID, wsMessage.ReceiverID, string(wsMessage.Content))
		if err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("Failed to send message"))
		}
	}
	
}