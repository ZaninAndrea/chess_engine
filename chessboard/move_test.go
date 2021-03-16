package chessboard

import "testing"

func TestNewMove(t *testing.T) {
	m := NewMove(A1, A2, NoPiece, ResetHalfMoveClockFlag)
	if m.From() != A1 {
		t.Errorf("Move from should be a1, %s was returned instead", m.From())
	}
	if m.To() != A2 {
		t.Errorf("Move from should be a1, %s was returned instead", m.To())
	}
	if m.Promotion() != NoPiece {
		t.Errorf("Move promotion should be -, %s was returned instead", m.Promotion())
	}

	m = NewMove(H7, H8, WhiteBishop, NoFlag)
	if m.From() != H7 {
		t.Errorf("Move from should be h7, %s was returned instead", m.From())
	}
	if m.To() != H8 {
		t.Errorf("Move from should be h8, %s was returned instead", m.To())
	}
	if m.Promotion() != WhiteBishop {
		t.Errorf("Move promotion should be -, %s was returned instead", m.Promotion())
	}
}
