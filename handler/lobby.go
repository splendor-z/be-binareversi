package handler

import (
	"be-binareversi/websocket"

	"github.com/gin-gonic/gin"
)

func LobbyWebSocket(c *gin.Context) {
	websocket.HandleLobby(c.Writer, c.Request)
}
