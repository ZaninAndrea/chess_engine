package chessboard

// LegalMoves returns the legal moves in the current position and caches them
func (game *Game) LegalMoves() []*Move {
	if game.Position.legalMoves != nil {
		return game.Position.legalMoves
	}
	// The moves list starts empty but we set the capacity to 35
	// to avoid expanding the slice's capacity many times
	moves := make([]*Move, 0, 35)

	var ownPieces Bitboard
	if game.Position.turn == WhiteColor {
		ownPieces = game.Position.board.whiteSquares
	} else {
		ownPieces = game.Position.board.blackSquares
	}

	// ADD KING MOVES
	var kingSquareIndex square
	if game.Position.turn == WhiteColor {
		kingSquareIndex = game.Position.board.whiteKingSquare
	} else {
		kingSquareIndex = game.Position.board.blackKingSquare
	}

	// Get precomputed king moves for that square and remove self-captures
	kingMovesBB := game.PrecomputedData.KingMoves[kingSquareIndex]
	kingMovesBB &^= ownPieces

	// Iterating target squares in the bitboard and add moves to the list
	fromSquare := kingSquareIndex
	for kingMovesBB != 0 {
		toSquare := square(kingMovesBB.LeastSignificantBit())
		moves = append(moves, &Move{from: fromSquare, to: toSquare})

		kingMovesBB = kingMovesBB & (kingMovesBB - 1)
	}

	game.Position.legalMoves = moves
	return moves
}
