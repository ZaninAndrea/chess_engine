package chessboard

import "fmt"

type MoveFlags int

const (
	ResetHalfMoveClockFlag = 1 << iota
	WhiteKingCastleFlag
	WhiteQueenCastleFlag
	BlackKingCastleFlag
	BlackQueenCastleFlag
	EnPassantFlag
)

// CastlesFlag is a mask to check whether any of the castles flags is set
const CastlesMask = WhiteKingCastleFlag | WhiteQueenCastleFlag | BlackKingCastleFlag | BlackQueenCastleFlag

// Move contains the informations about a move
type Move struct {
	from      square
	to        square
	promotion Piece
	flags     MoveFlags
}

func (m Move) String() string {
	if m.promotion != NoPiece {
		return fmt.Sprintf("%s-%s=%s", m.from, m.to, m.promotion)
	}

	return fmt.Sprintf("%s-%s", m.from, m.to)
}

func (m *Move) ShouldResetHalfMoveClock() bool {
	return m.flags&ResetHalfMoveClockFlag != 0
}
func (m *Move) IsCastle() bool {
	return m.flags&CastlesMask != 0
}
func (m *Move) IsEnPassant() bool {
	return m.flags&EnPassantFlag != 0
}
