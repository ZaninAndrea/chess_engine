package main

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

	b = SquareBitboard(Square{File(4), Rank(2)})
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

	b = SquareBitboard(Square{File(4), Rank(2)})
	got = b.ParallelPopCount()
	if got != 1 {
		t.Errorf("ParallelCount for Square bitboard should be 1, %d was returned instead", got)
	}
}
