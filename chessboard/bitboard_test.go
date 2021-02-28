package chessboard

import (
	"math/bits"
	"testing"
)

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

func BenchmarkKernighanPopCount64Bits(b *testing.B) {
	bb := ^Bitboard(0)

	for i := 0; i < b.N; i++ {
		bb.KernighanPopCount()
	}
}
func BenchmarkKernighanPopCount1Bits(b *testing.B) {
	bb := Bitboard(1)

	for i := 0; i < b.N; i++ {
		bb.KernighanPopCount()
	}
}

// TestPopCount checks that PopCount computes the correct counts
func TestPopCount(t *testing.T) {
	t.Run("BB=0b01110101", func(t *testing.T) {
		b := Bitboard(0b01110101)
		got := b.PopCount()

		if got != 5 {
			t.Errorf("ParallelCount for 0b01110101 should be 5, %d was returned instead", got)
		}
	})

	t.Run("BB=0", func(t *testing.T) {
		b := Bitboard(0)
		got := b.PopCount()
		if got != 0 {
			t.Errorf("ParallelCount for 0 should be 0, %d was returned instead", got)
		}
	})

	t.Run("BB=E3", func(t *testing.T) {
		b := E3.Bitboard()
		got := b.PopCount()
		if got != 1 {
			t.Errorf("ParallelCount for Square bitboard should be 1, %d was returned instead", got)
		}
	})
}

func BenchmarkPopCount64Bits(b *testing.B) {
	bb := ^Bitboard(0)

	for i := 0; i < b.N; i++ {
		bb.PopCount()
	}
}
func BenchmarkPopCount1Bits(b *testing.B) {
	bb := Bitboard(1)

	for i := 0; i < b.N; i++ {
		bb.PopCount()
	}
}

// TestPopCount checks that PopCount computes the correct counts
func TestPopCountNoMultiply(t *testing.T) {
	t.Run("BB=0b01110101", func(t *testing.T) {
		b := Bitboard(0b01110101)
		got := b.PopCountNoMultiply()

		if got != 5 {
			t.Errorf("PopCountNoMultiply for 0b01110101 should be 5, %d was returned instead", got)
		}
	})

	t.Run("BB=0", func(t *testing.T) {
		b := Bitboard(0)
		got := b.PopCountNoMultiply()
		if got != 0 {
			t.Errorf("PopCountNoMultiply for 0 should be 0, %d was returned instead", got)
		}
	})

	t.Run("BB=E3", func(t *testing.T) {
		b := E3.Bitboard()
		got := b.PopCountNoMultiply()
		if got != 1 {
			t.Errorf("PopCountNoMultiply for Square bitboard should be 1, %d was returned instead", got)
		}
	})
}

func BenchmarkPopCountNoMultiply64Bits(b *testing.B) {
	bb := ^Bitboard(0)

	for i := 0; i < b.N; i++ {
		bb.PopCount()
	}
}
func BenchmarkPopCountNoMultiply1Bits(b *testing.B) {
	bb := Bitboard(1)

	for i := 0; i < b.N; i++ {
		bb.PopCount()
	}
}

func TestLeastSignificant1Bit(t *testing.T) {
	b := Bitboard(0b01110101)
	got := b.LeastSignificant1Bit()

	if got != 0 {
		t.Errorf("LeastSignificant1Bit for 0b01110101 should be 0, %d was returned instead", got)
	}

	b = Bitboard(0b1110000)
	got = b.LeastSignificant1Bit()
	if got != 4 {
		t.Errorf("LeastSignificant1Bit for 0b1110000 should be 4, %d was returned instead", got)
	}

	b = E3.Bitboard()
	got = b.LeastSignificant1Bit()
	if got != 20 {
		t.Errorf("LeastSignificant1Bit for E3 bitboard should be 20, %d was returned instead", got)
	}
}

func BenchmarkLeastSignificant1BitPopCount(b *testing.B) {
	bb := Bitboard(0x10010)

	for i := 0; i < b.N; i++ {
		bb.LeastSignificant1Bit()
	}
}

func BenchmarkLeastSignificant1BitDebruijn(b *testing.B) {
	bb := Bitboard(0x10010)

	for i := 0; i < b.N; i++ {
		bits.TrailingZeros64(uint64(bb))
	}
}

func TestClearLeastSignificant1Bit(t *testing.T) {
	b := Bitboard(0b01110101)
	b.ClearLeastSignificant1Bit()

	if b != Bitboard(0b01110100) {
		t.Errorf("ClearLeastSignificant1Bit for 0b01110101 should leave 0b01110100, %d was returned instead", b)
	}

	b = Bitboard(0b1110000)
	b.ClearLeastSignificant1Bit()
	if b != Bitboard(0b1100000) {
		t.Errorf("ClearLeastSignificant1Bit for 0b1110000 should leave 0b1100000, %d was returned instead", b)
	}

	b = E3.Bitboard()
	b.ClearLeastSignificant1Bit()
	if b != 0 {
		t.Errorf("ClearLeastSignificant1Bit for Square should leave 0, %d was returned instead", b)
	}
}

func TestIsSquareOccupied(t *testing.T) {
	b := Bitboard(0)

	if b.IsSquareOccupied(A1) {
		t.Errorf("IsSquareOccupied(A1) for empty bitboard should return false")
	}

	b = E3.Bitboard()
	if !b.IsSquareOccupied(E3) {
		t.Errorf("IsSquareOccupied(E3) for E3 bitboard should return true")
	}

	b = F7.Bitboard()
	if b.IsSquareOccupied(G1) {
		t.Errorf("IsSquareOccupied(G1) for F7 bitboard should return false")
	}
}

func BenchmarkIsSquareOccupied(b *testing.B) {
	bb := F7.Bitboard()

	for i := 0; i < b.N; i++ {
		bb.IsSquareOccupied(G1)
	}
}
