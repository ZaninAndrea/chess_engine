package chessboard

import "fmt"

// CastleRights stored information about which side each player
// is still allowed to castle to
type CastleRights struct {
	WhiteKingSide  bool
	WhiteQueenSide bool
	BlackKingSide  bool
	BlackQueenSide bool
}

func (rights CastleRights) String() string {
	s := ""

	if rights.WhiteKingSide {
		s += "K"
	}
	if rights.WhiteQueenSide {
		s += "Q"
	}
	if rights.BlackKingSide {
		s += "k"
	}
	if rights.BlackQueenSide {
		s += "q"
	}

	return s
}

// Position contains all the information about a give position in a game
// including turn, enpassant information, castle rights, valid moves in this position...
type Position struct {
	board           Board
	turn            Color
	castleRights    CastleRights
	enPassantSquare square
	halfMoveClock   int
	moveCount       int
	inCheck         bool
	legalMoves      []*Move
}

func (pos Position) String() string {
	s := pos.board.String() + "\n"
	s += fmt.Sprintf("Turn: %s, CastleRights: %s, EnPassantSquare: %s, HalfMove: %d, Move: %d", pos.turn, pos.castleRights, pos.enPassantSquare, pos.halfMoveClock, pos.moveCount)

	return s
}

// Move returns a new position applying the move, the operation is NOT in place
func (pos Position) Move(move *Move) Position {
	pos.board.Move(move)
	pos.turn = pos.turn.Other()
	pos.moveCount++

	// TODO: captures should reset the half move clock
	if (pos.board.bbWhitePawn|pos.board.bbBlackPawn)&move.to.Bitboard() != 0 {
		pos.halfMoveClock = 0
	} else {
		pos.halfMoveClock++
	}
	pos.legalMoves = nil

	// update castle rights
	if move.from == E1 {
		pos.castleRights.WhiteKingSide = false
		pos.castleRights.WhiteQueenSide = false
	}
	if move.from == E8 {
		pos.castleRights.BlackKingSide = false
		pos.castleRights.BlackQueenSide = false
	}
	if move.from == A1 || move.to == A1 {
		pos.castleRights.WhiteQueenSide = false
	}
	if move.from == A8 || move.to == A8 {
		pos.castleRights.BlackQueenSide = false
	}
	if move.from == H1 || move.to == H1 {
		pos.castleRights.WhiteKingSide = false
	}
	if move.from == H8 || move.to == H8 {
		pos.castleRights.BlackKingSide = false
	}

	return pos
}
