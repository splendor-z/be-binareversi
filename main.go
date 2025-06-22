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
			time.Sleep(10 * time.Minute) // 10分おきにチェック
			if err := db.DeleteOldRooms(360); err != nil {
				log.Println("Failed to delete old rooms:", err)
			} else {
				log.Println("Old rooms cleanup completed.")
			}
			if err := db.DeleteInactivePlayers(720); err != nil {
				log.Println("Failed to delete inactive players:", err)
			} else {
				log.Println("Inactive players cleanup completed.")
			}
		}
	}()

	r := gin.Default()
	router.Setup(r)
	r.Run(":8080")
}
