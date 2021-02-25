package chessboard

// Color holds the color of a piece
type Color int

// Possible colors of a piece or square on the board
const (
	WhiteColor Color = 1
	BlackColor Color = -1
	NoColor    Color = 0
)

func (c Color) String() string {
	switch c {
	case WhiteColor:
		return "White"
	case BlackColor:
		return "Black"
	case NoColor:
		return "NoColor"
	default:
		panic("Unrecognized color")
	}
}

// Other returns the other color
func (c Color) Other() Color {
	switch c {
	case WhiteColor:
		return BlackColor
	case BlackColor:
		return WhiteColor
	case NoColor:
		return NoColor
	default:
		panic("Unrecognized color")
	}
}

// Piece contains information about the piece type; e.g. white king, black rook, ...
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
