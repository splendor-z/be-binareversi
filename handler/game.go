package handler

import (
	"be-binareversi/websocket"

	"github.com/gin-gonic/gin"
)

func GameWebSocket(c *gin.Context) {
	roomID := c.Param("roomID")
	playerID := c.Param("playerID")

	// WebSocketハンドラーにplayerIDを渡す
	websocket.HandleGame(roomID, playerID, c.Writer, c.Request)
}
