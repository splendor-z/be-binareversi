package websocket

import (
	"be-binareversi/model"
	"net/http"

	"github.com/gorilla/websocket"
)

var lobbyClients = map[*websocket.Conn]bool{}
var lobbyBroadcast = make(chan interface{})

func HandleLobby(w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	lobbyClients[conn] = true

	roomList := make([]*model.Room, 0, len(model.Rooms))
	for _, room := range model.Rooms {
		roomList = append(roomList, room)
	}
	conn.WriteJSON(map[string]interface{}{
		"type":  "room_list",
		"rooms": roomList,
	})

	for {
		var msg map[string]string
		if err := conn.ReadJSON(&msg); err != nil {
			delete(lobbyClients, conn)
			break
		}

		switch msg["type"] {
		case "create_room":
			roomID := msg["roomID"]
			player := msg["player"]
			room := &model.Room{ID: roomID, Player1: player, IsFull: false}
			model.Rooms[roomID] = room
			lobbyBroadcast <- map[string]interface{}{"type": "room_created", "room": room}

		case "join_room":
			roomID := msg["roomID"]
			player := msg["player"]
			if room, ok := model.Rooms[roomID]; ok && !room.IsFull {
				room.Player2 = &player
				room.IsFull = true
				lobbyBroadcast <- map[string]interface{}{"type": "room_updated", "room": room}
			}
		}
	}
}

func init() {
	go func() {
		for msg := range lobbyBroadcast {
			for conn := range lobbyClients {
				if err := conn.WriteJSON(msg); err != nil {
					conn.Close()
					delete(lobbyClients, conn)
				}
			}
		}
	}()
}
