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

func TestMoveGenerationStart5Ply(t *testing.T) {
	game := NewGame()
	total := countAllMoves(&game, 5)

	if total != 4_865_609 {
		t.Errorf("Total moves up to 5 ply should be 4.865.609, %d was returned instead", total)
	}
}

func TestMoveGenerationOther5Ply(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	game := NewGameFromFEN("rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8")
	total := countAllMoves(&game, 5)

	if total != 89_941_194 {
		t.Errorf("Total moves up to 5 ply should be 89.941.194, %d was returned instead", total)
	}
}

func BenchmarkMoveGeneration6Ply(b *testing.B) {
	for i := 0; i < b.N; i++ {
		game := NewGame()
		total := countAllMoves(&game, 6)
		b.Log(total)
	}
}
