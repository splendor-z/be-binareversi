package router

import (
	"time"

	"be-binareversi/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:    []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:   []string{"Content-Length"},
		MaxAge:          12 * time.Hour,
	}))

	r.POST("/api/register", handler.RegisterPlayer)
	r.GET("/ws/lobby", handler.LobbyWebSocket)
	r.GET("/ws/game/:roomID", handler.GameWebSocket)
}
