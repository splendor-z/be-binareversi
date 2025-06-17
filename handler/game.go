package handler

import (
	"be-binareversi/websocket"

	"github.com/gin-gonic/gin"
)

func GameWebSocket(c *gin.Context) {
	roomID := c.Param("roomID")
	websocket.HandleGame(roomID, c.Writer, c.Request)
}
