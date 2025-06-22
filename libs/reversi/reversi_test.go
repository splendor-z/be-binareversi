package reversi

import (
	"testing"
)

func Test01_NewGameInitialization(t *testing.T) {
	game := NewGame("room1")
	game.PrintBoard()
	if game.Turn != Black {
		t.Errorf("Expected initial turn to be Black, got %d", game.Turn)
	}
	if game.Board[3][3] != White || game.Board[4][4] != White ||
		game.Board[3][4] != Black || game.Board[4][3] != Black {
		t.Error("Initial positions are incorrect")
	}
}

func Test02_ValidMovesAtStart(t *testing.T) {
	game := NewGame("room2")
	moves := game.GetValidMoves(Black)
	if len(moves) != 4 {
		t.Errorf("Expected 4 valid moves for Black, got %d", len(moves))
	}
}

func Test03_ValidPlacementAndFlipping(t *testing.T) {
	game := NewGame("room3")
	board, err := game.PlaceDisc(Black, 2, 3)
	if err != nil {
		t.Errorf("Expected valid move, got error: %v", err)
	}
	game.PrintBoard()
	if board[3][3] != Black {
		t.Error("Expected (3,3) to be flipped to Black")
	}
	if game.Turn != White {
		t.Error("Turn should switch to White")
	}
}

func Test04_InvalidPlacementOccupied(t *testing.T) {
	game := NewGame("room4")
	_, err := game.PlaceDisc(Black, 3, 3)
	if err == nil {
		t.Error("Expected error when placing on occupied cell")
	}
}

func Test05_InvalidPlacementNotYourTurn(t *testing.T) {
	game := NewGame("room5")
	_, err := game.PlaceDisc(White, 2, 3)
	if err == nil {
		t.Error("Expected error when playing out of turn")
	}
}

func Test06_InvalidPlacementOutOfRange(t *testing.T) {
	game := NewGame("room6")
	_, err := game.PlaceDisc(Black, -1, 9)
	if err == nil {
		t.Error("Expected error when placing out of range")
	}
}

func Test07_PassTurn(t *testing.T) {
	game := NewGame("room7")
	originalTurn := game.Turn
	game.PassTurn()
	if game.Turn == originalTurn {
		t.Error("PassTurn did not change the turn")
	}
}

func Test08_FullBoard_WinnerBlack(t *testing.T) {
	game := NewGame("room8")
	for i := range game.Board {
		for j := range game.Board[i] {
			game.Board[i][j] = Black
		}
	}
	game.PrintBoard()
	if !game.IsGameOver() {
		t.Error("Game should be over when board is full")
	}
	if game.GetWinner() != Black {
		t.Error("Winner should be Black")
	}
}

func Test09_FullBoard_Draw(t *testing.T) {
	game := NewGame("room9")
	for i := range game.Board {
		for j := range game.Board[i] {
			if (i+j)%2 == 0 {
				game.Board[i][j] = Black
			} else {
				game.Board[i][j] = White
			}
		}
	}
	game.PrintBoard()
	if !game.IsGameOver() {
		t.Error("Game should be over when board is full")
	}
	if game.GetWinner() != -1 {
		t.Error("Expected draw")
	}
}

func Test10_GameNotOverMidGame(t *testing.T) {
	game := NewGame("room10")
	// Only first few moves played
	game.PlaceDisc(Black, 2, 3)
	game.PlaceDisc(White, 2, 2)
	game.PrintBoard()
	if game.IsGameOver() {
		t.Error("Game should not be over in mid game")
	}
}

func Test11_PassWhenNoMoves(t *testing.T) {
	game := NewGame("room11")
	// Fill all with black except one empty surrounded by black
	for i := range game.Board {
		for j := range game.Board[i] {
			game.Board[i][j] = Black
		}
	}
	game.Board[0][0] = Empty // surrounded by Black
	game.Turn = White
	moves := game.GetValidMoves(White)
	if len(moves) != 0 {
		t.Error("Expected White to have no valid moves")
	}
	game.PassTurn()
	if game.Turn != Black {
		t.Error("Turn should switch to Black on pass")
	}
}

func Test13_GetValidMovesMap(t *testing.T) {
	game := NewGame("room13")
	movesMap := game.GetValidMovesMap(White)

	// 合法手のマスは9、それ以外は0になっているかチェック
	count := 0
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if movesMap[i][j] == 9 {
				count++
				if game.Board[i][j] != Empty {
					t.Errorf("Valid move marked at (%d,%d) but board is not empty", i, j)
				}
			} else if movesMap[i][j] != 0 {
				t.Errorf("Invalid value %d at (%d,%d) in movesMap", movesMap[i][j], i, j)
			}
		}
	}
	if count == 0 {
		t.Errorf("No valid moves found, but expected some")
	}

	t.Log("Board with valid moves (9):")
	game.PrintBoardWithMovesMap(movesMap)
}
