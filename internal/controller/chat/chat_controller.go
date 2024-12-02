package chat

import "github.com/gin-gonic/gin"

type ChatController interface {
	HandleWebsocket(c *gin.Context)
}