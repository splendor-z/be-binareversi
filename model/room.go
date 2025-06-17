package model

type Room struct {
	ID      string `json:"id"`
	Player1 string `json:"player1"`
	Player2 string `json:"player2,omitempty"`
	IsFull  bool   `json:"isFull"`
}

var Rooms = map[string]*Room{}
