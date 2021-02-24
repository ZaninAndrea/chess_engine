package chessboard

import "testing"

// TestSquareBitboard checks that SquareBitboard generates the correct numbers
func TestSquareBitboard(t *testing.T) {
	// Bitboard for A1
	got := (Square{File(0), Rank(0)}).Bitboard()
	if got != 1 {
		t.Errorf("Bitboard for A1 should be 1, %d was returned instead", got)
	}

	// Bitboard for D3
	got = (Square{File(3), Rank(2)}).Bitboard()
	if got != 524_288 {
		t.Errorf("Bitboard for D3 should be 524.288, %d was returned instead", got)
	}

	// Bitboard for H8
	got = (Square{File(7), Rank(7)}).Bitboard()
	if got != 9_223_372_036_854_775_808 {
		t.Errorf("Bitboard for H8 should be 9.223.372.036.854.775.808, %d was returned instead", got)
	}

}

// TestSquareBitboard checks that Color returns the correct color
func TestColor(t *testing.T) {
	// Color for A1
	got := (Square{File(0), Rank(0)}).Color()
	if got != BlackColor {
		t.Errorf("Bitboard for A1 should be Black, %s was returned instead", got)
	}

	// Color for D3
	got = Square{File(3), Rank(2)}.Color()
	if got != WhiteColor {
		t.Errorf("Bitboard for D3 should be White, %s was returned instead", got)
	}

	// Color for H8
	got = Square{File(7), Rank(7)}.Color()
	if got != BlackColor {
		t.Errorf("Bitboard for H8 should be Black, %s was returned instead", got)
	}
}

func TestSquareFromIndex(t *testing.T) {
	got := SquareFromIndex(0)
	if got.String() != "A1" {
		t.Errorf("Square from index 0 should be A1, %s was returned instead", got)
	}

	got = SquareFromIndex(4)
	if got.String() != "E1" {
		t.Errorf("Square from index 4 should be A5, %s was returned instead", got)
	}

	got = SquareFromIndex(10)
	if got.String() != "C2" {
		t.Errorf("Square from index 10 should be C2, %s was returned instead", got)
	}

	got = SquareFromIndex(63)
	if got.String() != "H8" {
		t.Errorf("Square from index 63 should be H8, %s was returned instead", got)
	}
}