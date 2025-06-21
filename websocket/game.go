package websocket

import (
	"be-binareversi/db"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var gameClients = make(map[string]map[*websocket.Conn]string)
var currentTurn = make(map[string]string)

func HandleGame(roomID string, playerID string, w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	room, err := db.GetRoomByID(roomID)
	if err != nil {
		conn.WriteJSON(map[string]string{"error": "room not found"})
		return
	}

	// プレイヤー1 or プレイヤー2か確認
	if playerID != room.Player1 && (room.Player2 == nil || playerID != *room.Player2) {
		conn.WriteControl(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, "unauthorized player"),
			time.Now().Add(time.Second),
		)
		conn.Close()
		return
	}

	if gameClients[roomID] == nil {
		gameClients[roomID] = make(map[*websocket.Conn]string)
	}
	gameClients[roomID][conn] = playerID

	// 初期ターン状態（黒番: Player1）
	if _, exists := currentTurn[roomID]; !exists {
		currentTurn[roomID] = room.Player1
	}

	for {
		var msg map[string]interface{}
		if err := conn.ReadJSON(&msg); err != nil {
			delete(gameClients[roomID], conn)
			break
		}

		switch msg["type"] {
		case "join":
			broadcastToRoom(roomID, map[string]interface{}{
				"type":     "game_start",
				"playerID": playerID,
			})

		case "move":
			// ターン制御: 今のターンのプレイヤーかどうか
			if currentTurn[roomID] != playerID {
				conn.WriteJSON(map[string]string{"error": "not your turn"})
				continue
			}

			// 石を置く処理（盤面更新など）ここで入れる
			// ...

			// ターン交代
			if room.Player2 != nil && playerID == room.Player1 {
				currentTurn[roomID] = *room.Player2
			} else {
				currentTurn[roomID] = room.Player1
			}

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
