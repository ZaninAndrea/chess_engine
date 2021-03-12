package chessboard

import "testing"

// TestSquareBitboard checks that SquareBitboard generates the correct numbers
func TestSquareBitboard(t *testing.T) {
	// Bitboard for A1
	got := A1.Bitboard()
	if got != 1 {
		t.Errorf("Bitboard for A1 should be 1, %d was returned instead", got)
	}

	// Bitboard for D3
	got = D3.Bitboard()
	if got != 524_288 {
		t.Errorf("Bitboard for D3 should be 524.288, %d was returned instead", got)
	}

	// Bitboard for H8
	got = H8.Bitboard()
	if got != 9_223_372_036_854_775_808 {
		t.Errorf("Bitboard for H8 should be 9.223.372.036.854.775.808, %d was returned instead", got)
	}

}

// TestSquareBitboard checks that Color returns the correct color
func TestColor(t *testing.T) {
	// Color for A1
	got := A1.Color()
	if got != BlackColor {
		t.Errorf("Bitboard for A1 should be Black, %s was returned instead", got)
	}

	// Color for D3
	got = D3.Color()
	if got != WhiteColor {
		t.Errorf("Bitboard for D3 should be White, %s was returned instead", got)
	}

	// Color for H8
	got = H8.Color()
	if got != BlackColor {
		t.Errorf("Bitboard for H8 should be Black, %s was returned instead", got)
	}
}

func TestSquareFromIndex(t *testing.T) {
	got := square(0)
	if got.String() != "a1" {
		t.Errorf("Square from index 0 should be a1, %s was returned instead", got)
	}

	got = square(4)
	if got.String() != "e1" {
		t.Errorf("Square from index 4 should be a5, %s was returned instead", got)
	}

	got = square(10)
	if got.String() != "c2" {
		t.Errorf("Square from index 10 should be c2, %s was returned instead", got)
	}

	got = square(63)
	if got.String() != "h8" {
		t.Errorf("Square from index 63 should be h8, %s was returned instead", got)
	}
}
