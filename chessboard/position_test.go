package chessboard

import "testing"

func TestFENGeneration(t *testing.T) {
	got := NewGame().position.FEN()
	if got != "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1" {
		t.Errorf("FEN at starting position should be rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1, %s was returned instead", got)
	}
}
