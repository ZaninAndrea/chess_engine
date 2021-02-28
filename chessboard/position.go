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
	hash            uint64
}

func (pos Position) String() string {
	s := pos.board.String() + "\n"
	s += fmt.Sprintf("Turn: %s, CastleRights: %s, EnPassantSquare: %s, HalfMove: %d, Move: %d", pos.turn, pos.castleRights, pos.enPassantSquare, pos.halfMoveClock, pos.moveCount)

	return s
}

// Hash returns the zobrist hash for the position
func (pos *Position) Hash() uint64 {
	return pos.hash
}

// Move returns a new position applying the move, the operation is NOT in place
func (pos Position) Move(move *Move) Position {
	pos.hash ^= pos.board.Move(move)
	pos.turn = pos.turn.Other()
	pos.hash ^= zobristHashBlackTurn

	pos.moveCount++
	if move.ShouldResetHalfMoveClock() {
		pos.halfMoveClock = 0
	} else {
		pos.halfMoveClock++
	}

	pos.legalMoves = nil

	// update castle rights
	if move.from == E1 {
		if pos.castleRights.WhiteKingSide {
			pos.hash ^= zobristHashWhiteKingCastle
			pos.castleRights.WhiteKingSide = false

		}
		if pos.castleRights.WhiteQueenSide {
			pos.hash ^= zobristHashWhiteQueenCastle
			pos.castleRights.WhiteQueenSide = false
		}
	}
	if move.from == E8 {
		if pos.castleRights.BlackKingSide {
			pos.hash ^= zobristHashBlackKingCastle
			pos.castleRights.BlackKingSide = false

		}
		if pos.castleRights.BlackQueenSide {
			pos.hash ^= zobristHashBlackQueenCastle
			pos.castleRights.BlackQueenSide = false
		}
	}
	if (move.from == A1 || move.to == A1) && pos.castleRights.WhiteQueenSide {
		pos.hash ^= zobristHashWhiteQueenCastle
		pos.castleRights.WhiteQueenSide = false
	}
	if (move.from == A8 || move.to == A8) && pos.castleRights.BlackQueenSide {
		pos.hash ^= zobristHashBlackQueenCastle
		pos.castleRights.BlackQueenSide = false
	}
	if (move.from == H1 || move.to == H1) && pos.castleRights.WhiteKingSide {
		pos.hash ^= zobristHashWhiteKingCastle
		pos.castleRights.WhiteKingSide = false
	}
	if (move.from == H8 || move.to == H8) && pos.castleRights.BlackKingSide {
		pos.hash ^= zobristHashBlackKingCastle
		pos.castleRights.BlackKingSide = false
	}

	// Remove previous enpassant square from the hash
	if pos.enPassantSquare != NoSquare {
		pos.hash ^= zobristHashEnPassant[pos.enPassantSquare%8]
	}
	// update en passant square
	if move.IsDoublePawnPush() {
		pos.enPassantSquare = (move.to + move.from) / 2
		pos.hash ^= zobristHashEnPassant[pos.enPassantSquare%8]
	} else {
		pos.enPassantSquare = NoSquare
	}

	return pos
}
