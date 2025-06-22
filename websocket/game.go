package websocket

import (
	"be-binareversi/db"
	"be-binareversi/libs/bitop"
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
var playerPassCounts = make(map[string]map[string]int)
var lastPassPlayer = make(map[string]string)
var playerOperatorCounts = make(map[string]map[string]map[string]int)

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
			var boardToSend [8][8]int
			if game.GetTurn() == playerColor {
				boardToSend = game.GetBoardWithValidMoves(playerColor)
			} else {
				boardToSend = game.GetBoard()
			}

			conn.WriteJSON(map[string]interface{}{
				"type":        "game_start",
				"playerID":    playerID,
				"yourColor":   playerColor,
				"board":       boardToSend,
				"currentTurn": (game.GetTurnCount() + 1) / 2,
				"isYourTurn":  (game.GetTurn() == playerColor),
			})

		case "move":
			game.IncrementTurnCount()
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
				var boardToSend [8][8]int
				if game.GetTurn() == color {
					boardToSend = game.GetBoardWithValidMoves(color)
				} else {
					boardToSend = board
				}

				c.WriteJSON(map[string]interface{}{
					"type":        "board_update",
					"board":       boardToSend,
					"currentTurn": (game.GetTurnCount() + 1) / 2,
					"isYourTurn":  (game.GetTurn() == color),
				})
			}

			if game.IsGameOver() {
				broadcastToRoom(roomID, map[string]interface{}{
					"type":   "game_over",
					"winner": game.GetWinner(),
				})
			}

		case "operation":
			game.IncrementTurnCount()
			rowRaw, rowOk := msg["row"].(float64)
			valueRaw, valueOk := msg["value"].(float64)
			operator, opOk := msg["operator"].(string)

			if playerOperatorCounts[roomID] == nil {
				playerOperatorCounts[roomID] = make(map[string]map[string]int)
			}
			if playerOperatorCounts[roomID][playerID] == nil {
				playerOperatorCounts[roomID][playerID] = map[string]int{"+": 0, "*": 0}
			}

			if playerOperatorCounts[roomID][playerID][operator] >= 2 {
				conn.WriteJSON(map[string]string{"error": "Operator " + operator + " used too many times (max 2)."})
				continue
			}
			playerOperatorCounts[roomID][playerID][operator]++

			if !rowOk || !valueOk || !opOk {
				conn.WriteJSON(map[string]string{"error": "missing or invalid operation parameters"})
				continue
			}
			rowIndex := int(rowRaw)
			value := int(valueRaw)

			if rowIndex < 0 || rowIndex >= 8 {
				conn.WriteJSON(map[string]string{"error": "row index out of bounds"})
				continue
			}

			// 対象の行を取得し演算
			row := game.GetBoard()[rowIndex]
			newRow, err := bitop.ApplyBitOperation(row, value, operator)
			if err != nil {
				conn.WriteJSON(map[string]string{"error": err.Error()})
				continue
			}

			// 盤面の更新
			newBoard := game.GetBoard()
			newBoard[rowIndex] = newRow
			game.SetBoard(newBoard)
			game.PassTurn()

			// 全クライアントに board_update を送信
			for c, pid := range gameClients[roomID] {
				color := playerColors[roomID][pid]
				var boardToSend [8][8]int
				if game.GetTurn() == color {
					boardToSend = game.GetBoardWithValidMoves(color)
				} else {
					boardToSend = game.GetBoard()
				}

				c.WriteJSON(map[string]interface{}{
					"type":        "board_update",
					"board":       boardToSend,
					"currentTurn": (game.GetTurnCount() + 1) / 2,
					"isYourTurn":  (game.GetTurn() == color),
				})
			}

		case "surrender":
			// 通知: surrender したプレイヤーが敗北
			var winner int
			if playerColor == reversi.Black {
				winner = reversi.White
			} else {
				winner = reversi.Black
			}

			broadcastToRoom(roomID, map[string]interface{}{
				"type":   "game_over",
				"winner": winner,
			})

		case "pass":
			if playerPassCounts[roomID] == nil {
				playerPassCounts[roomID] = make(map[string]int)
			}
			playerPassCounts[roomID][playerID]++

			if playerPassCounts[roomID][playerID] > 3 {
				conn.WriteJSON(map[string]string{
					"error": "You have exceeded the maximum number of passes (3).",
				})
				playerPassCounts[roomID][playerID] = 3 // 上限固定
				continue
			}

			// 連続パス判定
			if lastPassPlayer[roomID] != "" && lastPassPlayer[roomID] != playerID {
				// 2人連続でパスされた → 勝者判定
				board := game.GetBoard()
				blackCount, whiteCount := 0, 0
				for _, row := range board {
					for _, cell := range row {
						if cell == reversi.Black {
							blackCount++
						} else if cell == reversi.White {
							whiteCount++
						}
					}
				}
				winner := -1
				if blackCount > whiteCount {
					winner = reversi.Black
				} else if whiteCount > blackCount {
					winner = reversi.White
				}
				broadcastToRoom(roomID, map[string]interface{}{
					"type":   "game_over",
					"winner": winner,
				})
			} else {
				game.IncrementTurnCount()
				// 手番変更、通知
				game.PassTurn()
				lastPassPlayer[roomID] = playerID

				for c, pid := range gameClients[roomID] {
					color := playerColors[roomID][pid]
					var boardToSend [8][8]int
					if game.GetTurn() == color {
						boardToSend = game.GetBoardWithValidMoves(color)
					} else {
						boardToSend = game.GetBoard()
					}

					c.WriteJSON(map[string]interface{}{
						"type":        "board_update",
						"board":       boardToSend,
						"currentTurn": (game.GetTurnCount() + 1) / 2,
						"isYourTurn":  (game.GetTurn() == color),
					})
				}
			}

		case "get_valid_moves":
			moves := game.GetValidMovesMap(playerColor)
			conn.WriteJSON(map[string]interface{}{
				"type":      "valid_moves",
				"moves_map": moves,
			})

		case "get_status":
			plusCount := 0
			mulCount := 0
			passCount := 0

			if playerOperatorCounts[roomID] != nil && playerOperatorCounts[roomID][playerID] != nil {
				plusCount = playerOperatorCounts[roomID][playerID]["+"]
				mulCount = playerOperatorCounts[roomID][playerID]["*"]
			}
			if playerPassCounts[roomID] != nil {
				passCount = playerPassCounts[roomID][playerID]
			}

			conn.WriteJSON(map[string]interface{}{
				"type":           "status_info",
				"remaining_plus": 2 - plusCount,
				"remaining_mul":  2 - mulCount,
				"remaining_pass": 3 - passCount,
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
