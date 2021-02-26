package chessboard

import (
	"testing"
)

func countAllMoves(game *Game, depth int) int {
	if depth == 0 {
		return 1
	}

	total := 0
	moves := game.LegalMoves()

	for _, move := range moves {
		game.Move(move)
		total += countAllMoves(game, depth-1)
		game.UndoMove()
	}

	return total
}

func BenchmarkMoveGeneration(b *testing.B) {
	for i := 0; i < b.N; i++ {
		game := NewGame()
		total := countAllMoves(&game, 6)
		b.Log(total)
	}
}
