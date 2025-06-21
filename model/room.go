package model

import "time"

type Room struct {
	ID        string    `json:"id" gorm:"not null;column:id;primaryKey"`
	Player1   string    `json:"player1" gorm:"not null;column:player1"`
	Player2   *string   `json:"player2,omitempty" gorm:"column:player2"`
	IsFull    bool      `json:"isFull" gorm:"column:is_full"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;autoCreateTime"`
}

var Rooms = map[string]*Room{}
