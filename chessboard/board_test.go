package chessboard

import (
	"testing"
)

func TestPiece(t *testing.T) {
	b := Board{bbWhitePawn: A5.Bitboard()}
	got := b.Piece(A5)

	if got != WhitePawn {
		t.Errorf("Piece in A5 should be ♙, %s was returned instead", got)
	}

	b = Board{bbWhiteRook: B8.Bitboard()}
	got = b.Piece(B8)
	if got != WhiteRook {
		t.Errorf("Piece in B8 should be ♖, %s was returned instead", got)
	}

	b = Board{}
	got = b.Piece(A1)
	if got != NoPiece {
		t.Errorf("Piece in A1 should be -, %s was returned instead", got)
	}
}

func TestMove(t *testing.T) {
	b := Board{
		bbWhiteKing: C2.Bitboard(),
		bbBlackPawn: C3.Bitboard(),
		bbBlackKing: E3.Bitboard(),
	}

	b.Move(&Move{from: C2, to: C3})

	if b.Piece(C2) != NoPiece {
		t.Error("Piece should be removed from starting square")
	} else if b.Piece(C3) != WhiteKing {
		t.Error("Piece should be placed in target square")
	} else if b.bbBlackPawn != 0 {
		t.Error("Piece in target square should be captured")
	}
}
