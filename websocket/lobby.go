package websocket

import (
	"be-binareversi/db"
	"be-binareversi/model"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// WebSocketクライアント管理
var lobbyClients = map[*websocket.Conn]bool{}
var lobbyBroadcast = make(chan interface{})
var roomMu sync.RWMutex

// レスポンス用構造体（ID）
type RoomResponse struct {
	ID        string    `json:"id"`
	Player1   string    `json:"player1"`           // id
	Player2   string    `json:"player2,omitempty"` // id
	IsFull    bool      `json:"isFull"`
	CreatedAt time.Time `json:"createdAt"`
}

func HandleLobby(w http.ResponseWriter, r *http.Request) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer func() {
		conn.Close()
		delete(lobbyClients, conn)
	}()

	lobbyClients[conn] = true
	for {
		var msg map[string]string
		if err := conn.ReadJSON(&msg); err != nil {
			break
		}

		switch msg["type"] {
		case "room_init":
			roomMu.RLock()
			roomList := []*RoomResponse{}
			rooms, _ := db.GetAllRooms()
			for _, room := range rooms {
				player1, _ := db.GetPlayerByID(room.Player1)
				var player2Name string
				if room.Player2 != nil {
					player2, _ := db.GetPlayerByID(*room.Player2)
					if player2 != nil {
						player2Name = player2.Name
					}
				}
				roomList = append(roomList, &RoomResponse{
					ID:        room.ID,
					Player1:   player1.Name,
					Player2:   player2Name,
					IsFull:    room.IsFull,
					CreatedAt: room.CreatedAt,
				})
			}
			roomMu.RUnlock()
			// roomListの長さを出力
			println("roomList length:", len(roomList))
			conn.WriteJSON(map[string]interface{}{
				"type":  "room_list",
				"rooms": roomList,
			})
		case "create_room":
			playerID := msg["playerID"]

			player, err := db.GetPlayerByID(playerID)
			if err != nil || player == nil {
				conn.WriteJSON(map[string]string{"error": "invalid playerID"})
				continue
			}

			roomMu.Lock()
			roomID := uuid.New().String()
			room := &model.Room{ID: roomID, Player1: playerID, IsFull: false}
			model.Rooms[roomID] = room
			db.CreateRoom(room)
			roomMu.Unlock()

			resp := RoomResponse{
				ID:        roomID,
				Player1:   player.Name,
				Player2:   "",
				IsFull:    false,
				CreatedAt: room.CreatedAt,
			}
			lobbyBroadcast <- map[string]interface{}{"type": "room_created", "room": resp}

		case "join_room":
			roomID := msg["roomID"]
			playerID := msg["playerID"]

			player, err := db.GetPlayerByID(playerID)
			if err != nil || player == nil {
				conn.WriteJSON(map[string]string{"error": "invalid playerID"})
				continue
			}

			roomMu.Lock()
			room, ok := model.Rooms[roomID]
			if ok && !room.IsFull {
				room.Player2 = &playerID
				room.IsFull = true
				db.UpdateRoom(room)
				roomMu.Unlock()

				player1, _ := db.GetPlayerByID(room.Player1)
				player2 := player.Name
				player1Name := "unknown"
				if player1 != nil {
					player1Name = player1.Name
				}

				resp := RoomResponse{
					ID:        room.ID,
					Player1:   player1Name,
					Player2:   player2,
					IsFull:    true,
					CreatedAt: room.CreatedAt,
				}
				lobbyBroadcast <- map[string]interface{}{"type": "room_updated", "room": resp}
			} else {
				roomMu.Unlock()
				conn.WriteJSON(map[string]string{"error": "room not found or already full"})
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
