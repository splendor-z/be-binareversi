package handler

import (
	"be-binareversi/db"
	"be-binareversi/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type RegisterRequest struct {
	Name string `json:"name"`
}

type RegisterResponse struct {
	UserID string `json:"userID"`
	Name   string `json:"name"`
}

func RegisterPlayer(c *gin.Context) {
	var req RegisterRequest
	if err := c.BindJSON(&req); err != nil || req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid name"})
		return
	}

	id := uuid.New().String()

	player := &model.Player{
		ID:         id,
		Name:       req.Name,
		LastUsedAt: time.Now(),
	}

	if err := db.CreatePlayer(player); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register player"})
		return
	}

	c.JSON(http.StatusOK, RegisterResponse{UserID: id, Name: req.Name})
}
