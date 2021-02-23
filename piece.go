package main

type Piece int

// Names for all the pieces of the board
const (
	NoPiece Piece = iota
	WhiteKing
	WhiteQueen
	WhiteRook
	WhiteBishop
	WhiteKnight
	WhitePawn
	BlackKing
	BlackQueen
	BlackRook
	BlackBishop
	BlackKnight
	BlackPawn
)

// Color holds the color of a piece
type Color int

// Possible colors of a piece
const (
	WhiteColor Color = 1
	BlackColor Color = -1
	NoColor    Color = 0
)

func (p Piece) String() string {
	switch {
	case p == NoPiece:
		return "-"
	case p == WhiteKing:
		return "♚"
	case p == WhiteQueen:
		return "♛"
	case p == WhiteRook:
		return "♜"
	case p == WhiteBishop:
		return "♝"
	case p == WhiteKnight:
		return "♞"
	case p == WhitePawn:
		return "♟︎"
	case p == BlackKing:
		return "♔"
	case p == BlackQueen:
		return "♕"
	case p == BlackRook:
		return "♖"
	case p == BlackBishop:
		return "♗"
	case p == BlackKnight:
		return "♘"
	case p == BlackPawn:
		return "♙"
	default:
		panic("Unknown piece")
	}
}

// Color returns the color of a piece
func (p Piece) Color() Color {
	switch {
	case p == NoPiece:
		return NoColor
	case p == WhiteKing:
		return WhiteColor
	case p == WhiteQueen:
		return WhiteColor
	case p == WhiteRook:
		return WhiteColor
	case p == WhiteBishop:
		return WhiteColor
	case p == WhiteKnight:
		return WhiteColor
	case p == WhitePawn:
		return WhiteColor
	case p == BlackKing:
		return BlackColor
	case p == BlackQueen:
		return BlackColor
	case p == BlackRook:
		return BlackColor
	case p == BlackBishop:
		return BlackColor
	case p == BlackKnight:
		return BlackColor
	case p == BlackPawn:
		return BlackColor
	default:
		panic("Unknown piece")
	}
}
