package main

import "fmt"

// Rank represents a rank (1 to 8) of the chessboard as 0 to 7
type Rank uint8

// File represents a file (A to H) of the chessboard as 0 to 7
type File uint8

// Square represents a square on the board
type Square struct {
	file File
	rank Rank
}

// SquareBitboard returns a bitboard with only the passed square on
func SquareBitboard(s Square) Bitboard {
	return Bitboard(1 << (int(s.rank)*8 + int(s.file)))
}

func (s Square) String() string {
	return fmt.Sprintf("%c%d", 65+int(s.file), s.rank+1)
}
