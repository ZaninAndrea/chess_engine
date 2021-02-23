package main

import "testing"

// TestSquareBitboard checks that SquareBitboard generates the correct numbers
func TestSquareBitboard(t *testing.T) {
	// Bitboard for A1
	got := SquareBitboard(Square{File(0), Rank(0)})
	if got != 1 {
		t.Errorf("Bitboard for A1 should be 1, %d was returned instead", got)
	}

	// Bitboard for D3
	got = SquareBitboard(Square{File(3), Rank(2)})
	if got != 524_288 {
		t.Errorf("Bitboard for D3 should be 524.288, %d was returned instead", got)
	}

	// Bitboard for H8
	got = SquareBitboard(Square{File(7), Rank(7)})
	if got != 9_223_372_036_854_775_808 {
		t.Errorf("Bitboard for H8 should be 9.223.372.036.854.775.808, %d was returned instead", got)
	}

}
