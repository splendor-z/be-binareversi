package reversi

import (
	"errors"
	"fmt"
)

const (
	White = 0 // 白い石
	Black = 1 // 黒い石
	Empty = 7 // 空きマス
)

// Point は座標 (x, y) を表す構造体
type Point struct {
	X int
	Y int
}

// Game はオセロゲームの状態を管理する構造体
type Game struct {
	RoomID string    // ゲームのルーム識別子
	Board  [8][8]int // 盤面の状態
	Turn   int       // 現在の手番（1=Black, 0=White）
}

var directions = []Point{
	{-1, -1}, {-1, 0}, {-1, 1},
	{0, -1}, {0, 1},
	{1, -1}, {1, 0}, {1, 1},
}

// 新しいオセロゲームを初期化して返す
// @param roomID ゲームを識別するためのID
// @return 初期化済みの *Game インスタンス
func NewGame(roomID string) *Game {
	g := &Game{
		RoomID: roomID,
		Turn:   Black,
	}
	g.initBoard()
	return g
}

// 盤面を初期状態に設定する
func (g *Game) initBoard() {
	for i := range g.Board {
		for j := range g.Board[i] {
			g.Board[i][j] = Empty
		}
	}
	g.Board[3][3], g.Board[4][4] = White, White
	g.Board[3][4], g.Board[4][3] = Black, Black
}

// 現在の盤面を返す
// @return [8][8]int 現在の盤面
func (g *Game) GetBoard() [8][8]int {
	return g.Board
}

// 現在の手番プレイヤーを返す
// @return int Black(1) または White(0)
func (g *Game) GetTurn() int {
	return g.Turn
}

// 手番を相手に交代する
func (g *Game) PassTurn() {
	if g.Turn == Black {
		g.Turn = White
	} else {
		g.Turn = Black
	}
}

// 指定プレイヤーが置ける合法手をすべて返す
// @param player プレイヤーの色（Black=1, White=0）
// @return []Point 合法手の座標リスト
func (g *Game) GetValidMoves(player int) []Point {
	var moves []Point
	for x := 0; x < 8; x++ {
		for y := 0; y < 8; y++ {
			if g.Board[x][y] == Empty && g.canPlace(player, x, y) {
				moves = append(moves, Point{x, y})
			}
		}
	}
	return moves
}

// 盤面を外部から上書きする
// @param newBoard 新しい盤面の状態
func (g *Game) SetBoard(newBoard [8][8]int) {
	g.Board = newBoard
}

// 指定座標に石を置き、盤面を更新する
// @param player 手番プレイヤー（Black=1, White=0）
// @param x X座標
// @param y Y座標
// @return [8][8]int 更新後の盤面
// @return error 不正な手であればエラー
func (g *Game) PlaceDisc(player, x, y int) ([8][8]int, error) {
	if x < 0 || x >= 8 || y < 0 || y >= 8 {
		return g.Board, errors.New("move out of board bounds")
	}
	if player != g.Turn {
		return g.Board, errors.New("not your turn")
	}
	if !g.canPlace(player, x, y) {
		return g.Board, errors.New("invalid move")
	}

	g.Board[x][y] = player
	g.flipDiscs(player, x, y)
	g.PassTurn()
	return g.Board, nil
}

// 現在の勝者を返す（ゲーム終了前でも呼び出し可能）
// @return int 勝者（Black=1, White=0, 引き分け=-1）
func (g *Game) GetWinner() int {
	blackCount, whiteCount := 0, 0
	for _, row := range g.Board {
		for _, cell := range row {
			if cell == Black {
				blackCount++
			} else if cell == White {
				whiteCount++
			}
		}
	}

	if blackCount > whiteCount {
		return Black
	} else if whiteCount > blackCount {
		return White
	}
	return -1
}

// ゲームが終了しているかどうかを返す
// @return bool 両者が合法手を持たなければ true
func (g *Game) IsGameOver() bool {
	return len(g.GetValidMoves(Black)) == 0 && len(g.GetValidMoves(White)) == 0
}

// 盤面を標準出力に表示（デバッグ用）
func (g *Game) PrintBoard() {
	for i := 0; i < 8; i++ {
		fmt.Print("|")
		for j := 0; j < 8; j++ {
			fmt.Printf(" %d", g.Board[i][j])
		}
		fmt.Println(" |")
	}
	fmt.Println()
}

// 指定座標に石を置けるかを判定
// @param player プレイヤーの色
// @param x X座標
// @param y Y座標
// @return bool 置けるなら true
func (g *Game) canPlace(player, x, y int) bool {
	if g.Board[x][y] != Empty {
		return false
	}
	for _, dir := range directions {
		if g.countFlippable(player, x, y, dir) > 0 {
			return true
		}
	}
	return false
}

// 指定方向に裏返せる石の数をカウント
// @param player プレイヤーの色
// @param x X座標
// @param y Y座標
// @param dir チェックする方向
// @return int 裏返せる相手の石の数
func (g *Game) countFlippable(player, x, y int, dir Point) int {
	opponent := 1 - player
	count := 0
	nx, ny := x+dir.X, y+dir.Y

	for nx >= 0 && nx < 8 && ny >= 0 && ny < 8 {
		if g.Board[nx][ny] == opponent {
			count++
		} else if g.Board[nx][ny] == player {
			return count
		} else {
			break
		}
		nx += dir.X
		ny += dir.Y
	}
	return 0
}

// 石を置いた後、裏返し処理を実行する
// @param player プレイヤーの色
// @param x X座標
// @param y Y座標
func (g *Game) flipDiscs(player, x, y int) {
	for _, dir := range directions {
		if g.countFlippable(player, x, y, dir) > 0 {
			nx, ny := x+dir.X, y+dir.Y
			for g.Board[nx][ny] == 1-player {
				g.Board[nx][ny] = player
				nx += dir.X
				ny += dir.Y
			}
		}
	}
}

// 合法手マップを取得（9で合法手を示す）
// @param player プレイヤーの色
// @return [8][8]int 合法手マップ（合法手=9, その他=0）
func (g *Game) GetValidMovesMap(player int) [8][8]int {
	var movesMap [8][8]int
	for _, move := range g.GetValidMoves(player) {
		movesMap[move.X][move.Y] = 9
	}
	return movesMap
}

// 盤面と合法手を同時に表示（デバッグ用）
// @param movesMap 合法手マップ（GetValidMovesMapの結果）
func (g *Game) PrintBoardWithMovesMap(movesMap [8][8]int) {
	for i := 0; i < 8; i++ {
		fmt.Print("|")
		for j := 0; j < 8; j++ {
			if movesMap[i][j] == 9 {
				// 合法手
				fmt.Printf(" 9")
			} else {
				// 盤面の石
				fmt.Printf(" %d", g.Board[i][j])
			}
		}
		fmt.Println(" |")
	}
	fmt.Println()
}
