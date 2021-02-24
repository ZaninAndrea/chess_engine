package chessboard

import "fmt"

// Rank represents a rank (1 to 8) of the chessboard as 0 to 7
type Rank uint8

// File represents a file (A to H) of the chessboard as 0 to 7
type File uint8

// Square represents a square on the board
type Square struct {
	File File
	Rank Rank
}

// Bitboard returns a bitboard with only the passed square on
func (s Square) Bitboard() Bitboard {
	return Bitboard(1 << (int(s.Rank)*8 + int(s.File)))
}

// SquareFromIndex returns a square from the passed index
func SquareFromIndex(index int) Square {
	file := File(index % 8)
	rank := Rank(index / 8)

	return Square{File: file, Rank: rank}
}

// Index returns which digit in the bitboard (from least meaningful to most meaningful)
// encodes this square
func (s Square) Index() int {
	return int(s.Rank)*8 + int(s.File)
}

func (s Square) String() string {
	return fmt.Sprintf("%c%d", 65+int(s.File), s.Rank+1)
}

// Color returns the color of the square
func (s Square) Color() Color {
	if (int(s.File)+int(s.Rank))%2 == 0 {
		return BlackColor
	}

	return WhiteColor
}
