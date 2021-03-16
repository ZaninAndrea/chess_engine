package chessboard

import "fmt"

type MoveFlags uint32

const NoFlag = MoveFlags(0)

const (
	ResetHalfMoveClockFlag = 1 << (iota + 16)
	WhiteKingCastleFlag
	WhiteQueenCastleFlag
	BlackKingCastleFlag
	BlackQueenCastleFlag
	WhiteEnPassantFlag
	BlackEnPassantFlag
	DoublePawnPushFlag
	IsCaptureFlag
)

// CastlesFlag is a mask to check whether any of the castles flags is set
const CastlesMask = WhiteKingCastleFlag | WhiteQueenCastleFlag | BlackKingCastleFlag | BlackQueenCastleFlag
const EnPassantMask = WhiteEnPassantFlag | BlackEnPassantFlag

// Move represents the from, to, promotion and flags information about a move.
// From least significant bit:
// - 6 bit: from
// - 6 bit: to
// - 4 bit: promotion
// - 9 bit: flags
type Move uint32

const fromMask = 0b111111
const toMask = 0b111111000000
const promotionMask = 0b1111000000000000

func NewMove(from square, to square, promotion Piece, flags MoveFlags) *Move {
	m := Move(from)
	m |= Move(to) << 6
	m |= Move(promotion) << 12
	m |= Move(flags)

	return &m
}

func (m Move) From() square {
	return square(m & fromMask)
}

func (m Move) To() square {
	return square((m & toMask) >> 6)
}

func (m Move) Promotion() Piece {
	return Piece((m & promotionMask) >> 12)
}

func (m Move) String() string {
	if m.Promotion() != NoPiece {
		return fmt.Sprintf("%s%s%s", m.From(), m.To(), m.Promotion())
	}

	return fmt.Sprintf("%s%s", m.From(), m.To())
}

func (m *Move) ShouldResetHalfMoveClock() bool {
	return uint32(*m)&uint32(ResetHalfMoveClockFlag) != 0
}
func (m *Move) IsCastle() bool {
	return uint32(*m)&uint32(CastlesMask) != 0
}
func (m *Move) IsEnPassant() bool {
	return uint32(*m)&uint32(EnPassantMask) != 0
}
func (m *Move) IsDoublePawnPush() bool {
	return uint32(*m)&uint32(DoublePawnPushFlag) != 0
}
func (m *Move) IsCapture() bool {
	return uint32(*m)&uint32(IsCaptureFlag) != 0
}

var NullMove = 0
