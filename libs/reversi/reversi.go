package reversi

import (
	"errors"
	"fmt"
)

const (
	White = 0
	Black = 1
	Empty = 7
)

type Point struct {
	X int
	Y int
}

type Game struct {
	RoomID string
	Board  [8][8]int
	Turn   int // 1=Black, 0=White
}

var directions = []Point{
	{-1, -1}, {-1, 0}, {-1, 1},
	{0, -1}, {0, 1},
	{1, -1}, {1, 0}, {1, 1},
}

// Game 初期化
func NewGame(roomID string) *Game {
	g := &Game{
		RoomID: roomID,
		Turn:   Black,
	}
	g.initBoard()
	return g
}

func (g *Game) initBoard() {
	for i := range g.Board {
		for j := range g.Board[i] {
			g.Board[i][j] = Empty
		}
	}
	g.Board[3][3], g.Board[4][4] = White, White
	g.Board[3][4], g.Board[4][3] = Black, Black
}

// 現在の盤面取得
func (g *Game) GetBoard() [8][8]int {
	return g.Board
}

// 現在の手番
func (g *Game) GetTurn() int {
	return g.Turn
}

// パス処理
func (g *Game) PassTurn() {
	if g.Turn == Black {
		g.Turn = White
	} else {
		g.Turn = Black
	}
}

// 合法手取得
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

// 盤面上書き
func (g *Game) SetBoard(newBoard [8][8]int) {
	g.Board = newBoard
}

// コマを置く
func (g *Game) PlaceDisc(player, x, y int) ([8][8]int, error) {
	if x < 0 || x >= 8 || y < 0 || y >= 8 {
		return g.Board, errors.New("move out of board bounds")
	}
	if player != g.Turn {
		return g.Board, errors.New("not player's turn")
	}
	if !g.canPlace(player, x, y) {
		return g.Board, errors.New("invalid move")
	}

	g.Board[x][y] = player
	g.flipDiscs(player, x, y)

	g.PassTurn()
	return g.Board, nil
}

// 勝者取得
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
	return -1 // 引き分け
}

// 終了判定
func (g *Game) IsGameOver() bool {
	return len(g.GetValidMoves(Black)) == 0 && len(g.GetValidMoves(White)) == 0
}

// 盤面表示（デバッグ用）
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

func (g *Game) GetValidMovesMap(player int) [8][8]int {
	var movesMap [8][8]int
	for _, move := range g.GetValidMoves(player) {
		movesMap[move.X][move.Y] = 9
	}
	return movesMap
}

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
