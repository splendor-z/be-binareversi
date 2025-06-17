package main

import (
	"be-binareversi/router"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	router.Setup(r)
	r.Run(":8080")
}
