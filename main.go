package main

import (
	"be-binareversi/router"
	"log"
	"time"

	"be-binareversi/db"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDatabase()

	go func() {
		for {
			time.Sleep(5 * time.Minute) // 5分おきにチェック
			if err := db.DeleteOldRooms(); err != nil {
				log.Println("Failed to delete old rooms:", err)
			} else {
				log.Println("Old rooms cleanup completed.")
			}
		}
	}()

	r := gin.Default()
	router.Setup(r)
	r.Run(":8080")
}
