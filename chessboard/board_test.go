package chessboard

import (
	"testing"
)

func TestPiece(t *testing.T) {
	a5 := Square{File(0), Rank(4)}
	b := Board{bbWhitePawn: a5.Bitboard()}
	got := b.Piece(a5)

	if got != WhitePawn {
		t.Errorf("Piece in A5 should be ♙, %s was returned instead", got)
	}

	b8 := Square{File(1), Rank(7)}
	b = Board{bbWhiteRook: b8.Bitboard()}
	got = b.Piece(b8)
	if got != WhiteRook {
		t.Errorf("Piece in B8 should be ♖, %s was returned instead", got)
	}

	b = Board{}
	got = b.Piece(Square{})
	if got != NoPiece {
		t.Errorf("Piece in B8 should be -, %s was returned instead", got)
	}
}

func TestMove(t *testing.T) {
	c2 := Square{File(2), Rank(1)}
	c3 := Square{File(2), Rank(2)}
	e3 := Square{File(4), Rank(2)}

	b := Board{
		bbWhiteKing: c2.Bitboard(),
		bbBlackPawn: c3.Bitboard(),
		bbBlackKing: e3.Bitboard(),
	}

	b.Move(Move{from: c2, to: c3})

	if b.Piece(c2) != NoPiece {
		t.Error("Piece should be removed from starting square")
	} else if b.Piece(c3) != WhiteKing {
		t.Error("Piece should be placed in target square")
	} else if b.bbBlackPawn != 0 {
		t.Error("Piece in target square should be captured")
	}
}
