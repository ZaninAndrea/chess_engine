package chessboard

import "testing"

// TestKernighanPopCount checks that KernighanPopCount computes the correct counts
func TestKernighanPopCount(t *testing.T) {
	b := Bitboard(0b01110101)
	got := b.KernighanPopCount()

	if got != 5 {
		t.Errorf("KernighanCount for 0b01110101 should be 5, %d was returned instead", got)
	}

	b = Bitboard(0)
	got = b.KernighanPopCount()
	if got != 0 {
		t.Errorf("KernighanCount for 0 should be 0, %d was returned instead", got)
	}

	b = E3.Bitboard()
	got = b.KernighanPopCount()
	if got != 1 {
		t.Errorf("KernighanCount for Square bitboard should be 1, %d was returned instead", got)
	}
}

// TestParallelPopCount checks that ParallelPopCount computes the correct counts
func TestParallelPopCount(t *testing.T) {
	b := Bitboard(0b01110101)
	got := b.ParallelPopCount()

	if got != 5 {
		t.Errorf("ParallelCount for 0b01110101 should be 5, %d was returned instead", got)
	}

	b = Bitboard(0)
	got = b.ParallelPopCount()
	if got != 0 {
		t.Errorf("ParallelCount for 0 should be 0, %d was returned instead", got)
	}

	b = E3.Bitboard()
	got = b.ParallelPopCount()
	if got != 1 {
		t.Errorf("ParallelCount for Square bitboard should be 1, %d was returned instead", got)
	}
}

func TestLeastSignificantBit(t *testing.T) {
	b := Bitboard(0b01110101)
	got := b.LeastSignificantBit()

	if got != 0 {
		t.Errorf("LeastSignificantBit for 0b01110101 should be 0, %d was returned instead", got)
	}

	b = Bitboard(0b1110000)
	got = b.LeastSignificantBit()
	if got != 4 {
		t.Errorf("LeastSignificantBit for 0b1110000 should be 4, %d was returned instead", got)
	}

	b = E3.Bitboard()
	got = b.LeastSignificantBit()
	if got != 20 {
		t.Errorf("LeastSignificantBit for E3 bitboard should be 20, %d was returned instead", got)
	}
}
