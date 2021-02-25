package chessboard

// Bitboard contains boolean informations about each square on the board.
// Conventionally A1 is the 0th bit and H8 is the 63rd bit
type Bitboard uint64

// KernighanPopCount counts the 1 bits in the bitboard with Khernighan's algorithm,
// this method is efficient for sparsely populated bitboards
func (b Bitboard) KernighanPopCount() int {
	if b == 0 {
		return 0
	}

	bb := b & (b - 1) // clear the least significant bit
	return bb.KernighanPopCount() + 1
}

// ParallelPopCount counts the 1 bits in the bitboard with a divide and conquer algorithm,
// this method is efficient for densely populated bitboards
func (b Bitboard) ParallelPopCount() int {
	var maskDuos uint64 = 0x5555555555555555    // in binary 0101010101...
	var maskNibbles uint64 = 0x3333333333333333 // in binary 00110011...
	var maskBytes uint64 = 0x0f0f0f0f0f0f0f0f   // in binary 0000111100001111...
	var factor uint64 = 0x0101010101010101      // the sum of 256 to the power of 0,1,2,3...

	bb := uint64(b) - (uint64(b)>>1)&maskDuos           // count bits in each duo
	bb = (bb & maskNibbles) + ((bb >> 2) & maskNibbles) // count bits in each nibble
	bb = (bb + (bb >> 4)) & maskBytes                   // count bits in each byte

	// use a multiplication to sum all the bytes count, the result
	// can be read from the 8 most significant bits,
	// on processors with fast multiplication this is faster
	// than continuing with the previous pattern
	return int((bb * factor) >> 56)
}

// LeastSignificantBit computes the index of the Least Significant Bit assuming that the bitboard is not empty
func (b Bitboard) LeastSignificantBit() int {
	// We can leverage the two-complement representation to rapidly generate
	// a bitboard with only the bits preceding the Least Significant Bit set as 1
	return ((b & -b) - 1).ParallelPopCount()
}

func (b Bitboard) String() string {
	s := ""

	for r := 7; r >= 0; r-- {
		for f := 0; f < 8; f++ {
			squareBB := (square(f + r*8)).Bitboard()
			if b&squareBB != 0 {
				s += "1 "
			} else {
				s += ". "
			}
		}
		s += "\n"
	}

	return s[:len(s)-1]
}
