package chat

import "github.com/gin-gonic/gin"

type ChatController interface {
	HandleWebsocket(c *gin.Context)
	GetAllChats(c *gin.Context)
}