package websocket

import (
	"be-binareversi/db"
	"be-binareversi/libs/reversi"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var gameClients = make(map[string]map[*websocket.Conn]string)
var gameInstances = make(map[string]*reversi.Game)
var playerColors = make(map[string]map[string]int)

func HandleGame(roomID string, playerID string, w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in HandleGame: %v", r)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}()

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

	if playerID != room.Player1 && (room.Player2 == nil || playerID != *room.Player2) {
		conn.WriteControl(websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, "unauthorized player"),
			time.Now().Add(time.Second),
		)
		return
	}

	if gameClients[roomID] == nil {
		gameClients[roomID] = make(map[*websocket.Conn]string)
	}
	gameClients[roomID][conn] = playerID

	if _, ok := gameInstances[roomID]; !ok {
		gameInstances[roomID] = reversi.NewGame(roomID)
	}
	game := gameInstances[roomID]

	if _, ok := playerColors[roomID]; !ok {
		playerColors[roomID] = make(map[string]int)
		playerColors[roomID][room.Player1] = reversi.Black
		if room.Player2 != nil {
			playerColors[roomID][*room.Player2] = reversi.White
		}
	}
	playerColor := playerColors[roomID][playerID]

	for {
		_, reader, err := conn.NextReader()
		if err != nil {
			delete(gameClients[roomID], conn)
			break
		}

		var msg map[string]interface{}
		if err := json.NewDecoder(reader).Decode(&msg); err != nil {
			conn.WriteJSON(map[string]string{"error": "invalid JSON"})
			continue
		}

		typeVal, ok := msg["type"].(string)
		if !ok {
			conn.WriteJSON(map[string]string{"error": "missing or invalid type"})
			continue
		}

		switch typeVal {
		case "join":
			conn.WriteJSON(map[string]interface{}{
				"type":       "game_start",
				"playerID":   playerID,
				"yourColor":  playerColor,
				"board":      game.GetBoard(),
				"isYourTurn": (game.GetTurn() == playerColor),
			})

		case "move":
			xRaw, xOk := msg["x"].(float64)
			yRaw, yOk := msg["y"].(float64)
			if !xOk || !yOk {
				conn.WriteJSON(map[string]string{"error": "invalid x or y"})
				continue
			}
			x, y := int(xRaw), int(yRaw)

			board, err := game.PlaceDisc(playerColor, x, y)
			if err != nil {
				conn.WriteJSON(map[string]string{"error": err.Error()})
				continue
			}

			for c, pid := range gameClients[roomID] {
				color := playerColors[roomID][pid]
				c.WriteJSON(map[string]interface{}{
					"type":       "board_update",
					"board":      board,
					"isYourTurn": (game.GetTurn() == color),
				})
			}

			if game.IsGameOver() {
				broadcastToRoom(roomID, map[string]interface{}{
					"type":   "game_over",
					"winner": game.GetWinner(),
				})
			}

		case "get_valid_moves":
			moves := game.GetValidMovesMap(playerColor)
			conn.WriteJSON(map[string]interface{}{
				"type":      "valid_moves",
				"moves_map": moves,
			})

		case "exit_room":
			db.DeleteRoom(roomID) //ルームの削除
			conn.WriteJSON(map[string]interface{}{
				"type":     "exited_room",
				"roomID":   roomID,
				"playerID": playerID,
			})

		default:
			conn.WriteJSON(map[string]string{"error": "unknown message type"})
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
