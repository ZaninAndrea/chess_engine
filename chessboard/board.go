package chessboard

import "fmt"

//TODO: are the support bitboards necessary?

// Board contains the position of all pieces on the chessboard
type Board struct {
	bbWhiteKing     Bitboard
	bbWhiteQueen    Bitboard
	bbWhiteRook     Bitboard
	bbWhiteBishop   Bitboard
	bbWhiteKnight   Bitboard
	bbWhitePawn     Bitboard
	bbBlackKing     Bitboard
	bbBlackQueen    Bitboard
	bbBlackRook     Bitboard
	bbBlackBishop   Bitboard
	bbBlackKnight   Bitboard
	bbBlackPawn     Bitboard
	whiteSquares    Bitboard
	blackSquares    Bitboard
	emptySquares    Bitboard
	whiteKingSquare square
	blackKingSquare square
}

// Piece returns the piece in a given square of the board
func (b *Board) Piece(s square) Piece {
	bbSquare := s.Bitboard()
	switch {
	case b.bbWhiteKing&bbSquare != 0:
		return WhiteKing
	case b.bbWhiteQueen&bbSquare != 0:
		return WhiteQueen
	case b.bbWhiteRook&bbSquare != 0:
		return WhiteRook
	case b.bbWhiteBishop&bbSquare != 0:
		return WhiteBishop
	case b.bbWhiteKnight&bbSquare != 0:
		return WhiteKnight
	case b.bbWhitePawn&bbSquare != 0:
		return WhitePawn
	case b.bbBlackKing&bbSquare != 0:
		return BlackKing
	case b.bbBlackQueen&bbSquare != 0:
		return BlackQueen
	case b.bbBlackRook&bbSquare != 0:
		return BlackRook
	case b.bbBlackBishop&bbSquare != 0:
		return BlackBishop
	case b.bbBlackKnight&bbSquare != 0:
		return BlackKnight
	case b.bbBlackPawn&bbSquare != 0:
		return BlackPawn
	default:
		return NoPiece
	}
}

func (b Board) String() string {
	s := ""
	for r := 7; r >= 0; r-- {
		s += fmt.Sprintf("%d ", r+1)
		for f := 0; f < 8; f++ {
			s += b.Piece(square(f+r*8)).String() + " "
		}
		s += "\n"
	}
	s += "  A B C D E F G H"

	return s
}

// Move updates the boarding moving the piece in the starting square to the target square
// it also captures the square in the target square if needed
func (b *Board) Move(move Move) {
	var piece Piece
	if move.promotion != NoPiece {
		piece = move.promotion
	} else {
		piece = b.Piece(move.from)
	}

	fromBB := move.from.Bitboard()
	toBB := move.to.Bitboard()
	othersBB := ^(fromBB | toBB)

	// Remove the piece in the start and target squares
	b.bbWhiteKing &= othersBB
	b.bbWhiteQueen &= (othersBB)
	b.bbWhiteRook &= (othersBB)
	b.bbWhiteBishop &= (othersBB)
	b.bbWhiteKnight &= (othersBB)
	b.bbWhitePawn &= (othersBB)
	b.bbBlackKing &= (othersBB)
	b.bbBlackQueen &= (othersBB)
	b.bbBlackRook &= (othersBB)
	b.bbBlackBishop &= (othersBB)
	b.bbBlackKnight &= (othersBB)
	b.bbBlackPawn &= (othersBB)

	// Place piece in target square
	switch piece {
	case WhiteKing:
		b.bbWhiteKing |= toBB
	case WhiteQueen:
		b.bbWhiteQueen |= toBB
	case WhiteRook:
		b.bbWhiteRook |= toBB
	case WhiteBishop:
		b.bbWhiteBishop |= toBB
	case WhiteKnight:
		b.bbWhiteKnight |= toBB
	case WhitePawn:
		b.bbWhitePawn |= toBB
	case BlackKing:
		b.bbBlackKing |= toBB
	case BlackQueen:
		b.bbBlackQueen |= toBB
	case BlackRook:
		b.bbBlackRook |= toBB
	case BlackBishop:
		b.bbBlackBishop |= toBB
	case BlackKnight:
		b.bbBlackKnight |= toBB
	case BlackPawn:
		b.bbBlackPawn |= toBB
	default:
		panic("From position is empty")
	}

	// Update king position
	if piece == WhiteKing {
		b.whiteKingSquare = move.to
	} else if piece == BlackKing {
		b.blackKingSquare = move.to
	}

	// Update summary bitboards
	b.emptySquares = (b.emptySquares | fromBB) & (^toBB)
	if piece.Color() == WhiteColor {
		b.whiteSquares = (b.whiteSquares | toBB) & (^fromBB)
	} else {
		b.blackSquares = (b.blackSquares | toBB) & (^fromBB)
	}

	// TODO: manage en passant captures
}
