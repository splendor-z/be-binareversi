package websocket

import (
	"be-binareversi/db"
	"be-binareversi/libs/reversi"

	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var gameClients = make(map[string]map[*websocket.Conn]string)
var gameInstances = make(map[string]*reversi.Game)

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

	// クライアント登録
	if gameClients[roomID] == nil {
		gameClients[roomID] = make(map[*websocket.Conn]string)
	}
	gameClients[roomID][conn] = playerID

	// ゲームインスタンス初期化
	if _, ok := gameInstances[roomID]; !ok {
		gameInstances[roomID] = reversi.NewGame(roomID)
	}
	game := gameInstances[roomID]

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
				"board":    game.GetBoard(),
				"turn":     game.GetTurn(),
			})

		case "move":
			x := int(msg["x"].(float64))
			y := int(msg["y"].(float64))
			player := int(msg["player"].(float64))

			// ターンチェック
			if game.GetTurn() != player {
				conn.WriteJSON(map[string]string{"error": "not your turn"})
				continue
			}

			board, err := game.PlaceDisc(player, x, y)
			if err != nil {
				conn.WriteJSON(map[string]string{"error": err.Error()})
				continue
			}

			broadcastToRoom(roomID, map[string]interface{}{
				"type":  "board_update",
				"board": board,
				"turn":  game.GetTurn(),
			})

			// ゲーム終了判定
			if game.IsGameOver() {
				broadcastToRoom(roomID, map[string]interface{}{
					"type":   "game_over",
					"winner": game.GetWinner(),
				})
			}

		case "get_valid_moves":
			player := int(msg["player"].(float64))
			moves := game.GetValidMovesMap(player)
			conn.WriteJSON(map[string]interface{}{
				"type":      "valid_moves",
				"moves_map": moves,
			})
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
