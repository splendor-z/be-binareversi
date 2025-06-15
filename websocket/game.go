package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var gameClients = make(map[string]map[*websocket.Conn]string)

func HandleGame(roomID string, w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	if gameClients[roomID] == nil {
		gameClients[roomID] = make(map[*websocket.Conn]string)
	}
	gameClients[roomID][conn] = ""

	for {
		var msg map[string]interface{}
		if err := conn.ReadJSON(&msg); err != nil {
			delete(gameClients[roomID], conn)
			break
		}

		switch msg["type"] {
		case "join":
			playerID := msg["playerID"].(string)
			gameClients[roomID][conn] = playerID
			broadcastToRoom(roomID, map[string]interface{}{
				"type":     "game_start",
				"playerID": playerID,
			})
		case "move":
			broadcastToRoom(roomID, msg)
		}
	}
}

func broadcastToRoom(roomID string, msg interface{}) {
	for conn := range gameClients[roomID] {
		if err := conn.WriteJSON(msg); err != nil {
			conn.Close()
			delete(gameClients[roomID], conn)
		}
	}
}
