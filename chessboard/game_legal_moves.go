package chessboard

// LegalMoves returns the legal moves in the current position and caches them
func (game *Game) LegalMoves() []*Move {
	if game.position.legalMoves != nil {
		return game.position.legalMoves
	}
	// The moves list starts empty but we set the capacity to 35
	// to avoid expanding the slice's capacity many times
	moves := make([]*Move, 0, 35)

	var ownPieces Bitboard
	if game.position.turn == WhiteColor {
		ownPieces = game.position.board.whiteSquares
	} else {
		ownPieces = game.position.board.blackSquares
	}

	computeKingMoves(game, &moves, &ownPieces)
	computeKnightMoves(game, &moves, &ownPieces)
	computeRookMoves(game, &moves, &ownPieces)
	computeBishopMoves(game, &moves, &ownPieces)
	computeQueenMoves(game, &moves, &ownPieces)
	computePawnMoves(game, &moves)

	game.position.legalMoves = moves
	return moves
}

func computeKingMoves(game *Game, moves *[]*Move, ownPieces *Bitboard) {
	var kingSquareIndex square
	if game.position.turn == WhiteColor {
		kingSquareIndex = game.position.board.whiteKingSquare
	} else {
		kingSquareIndex = game.position.board.blackKingSquare
	}

	// Get precomputed king moves for that square and remove self-captures
	kingMovesBB := game.precomputedData.KingMoves[kingSquareIndex]
	kingMovesBB &^= *ownPieces

	// Iterating target squares in the bitboard and add moves to the list
	fromSquare := kingSquareIndex
	for kingMovesBB != 0 {
		toSquare := square(kingMovesBB.LeastSignificantBit())
		*moves = append(*moves, &Move{from: fromSquare, to: toSquare})

		kingMovesBB.ClearLeastSignificantBit()
	}
}

func computeKnightMoves(game *Game, moves *[]*Move, ownPieces *Bitboard) {
	var knights Bitboard
	if game.position.turn == WhiteColor {
		knights = game.position.board.bbWhiteKnight
	} else {
		knights = game.position.board.bbBlackKnight
	}

	for knights != 0 {
		fromSquare := square(knights.LeastSignificantBit())
		knights.ClearLeastSignificantBit()

		// Get precomputed king moves for that square and remove self-captures
		knightMovesBB := game.precomputedData.KnightMoves[fromSquare]
		knightMovesBB &^= *ownPieces

		for knightMovesBB != 0 {
			toSquare := square(knightMovesBB.LeastSignificantBit())
			knightMovesBB.ClearLeastSignificantBit()

			*moves = append(*moves, &Move{from: fromSquare, to: toSquare})
		}
	}
}

func computePawnMoves(game *Game, moves *[]*Move) {
	if game.position.turn == WhiteColor {
		pawns := game.position.board.bbWhitePawn

		// iterate all white pawns on the board
		for pawns != 0 {
			fromSquare := square(pawns.LeastSignificantBit())
			pawns.ClearLeastSignificantBit()

			// Move forward
			forwardSquare := square(fromSquare + 8)
			if game.position.board.emptySquares.IsSquareOccupied(forwardSquare) {
				appendPawnMove(fromSquare, forwardSquare, moves)

				// Move forward by two squares
				if fromSquare < A3 {
					forwardSquare := square(fromSquare + 16)
					if game.position.board.emptySquares.IsSquareOccupied(forwardSquare) {
						appendPawnMove(fromSquare, forwardSquare, moves)
					}
				}
			}

			// Capture to the left
			if int(fromSquare)%8 != 0 {
				captureSquare := square(fromSquare + 7)
				if game.position.board.blackSquares.IsSquareOccupied(captureSquare) {
					appendPawnMove(fromSquare, captureSquare, moves)
				}
			}

			// Capture to the right
			if int(fromSquare)%8 != 7 {
				captureSquare := square(fromSquare + 9)
				if game.position.board.blackSquares.IsSquareOccupied(captureSquare) {
					appendPawnMove(fromSquare, captureSquare, moves)
				}
			}
		}
	} else {
		pawns := game.position.board.bbBlackPawn

		// iterate all black pawns on the board
		for pawns != 0 {
			fromSquare := square(pawns.LeastSignificantBit())
			pawns.ClearLeastSignificantBit()

			// Move forward
			forwardSquare := square(fromSquare - 8)
			if game.position.board.emptySquares.IsSquareOccupied(forwardSquare) {
				appendPawnMove(fromSquare, forwardSquare, moves)

				// Move forward by two squares
				if fromSquare > H6 {
					forwardSquare := square(fromSquare - 16)
					if game.position.board.emptySquares.IsSquareOccupied(forwardSquare) {
						appendPawnMove(fromSquare, forwardSquare, moves)
					}
				}
			}

			// Capture to the left
			if int(fromSquare)%8 != 0 {
				captureSquare := square(fromSquare - 9)
				if game.position.board.whiteSquares.IsSquareOccupied(captureSquare) {
					appendPawnMove(fromSquare, captureSquare, moves)
				}
			}

			// Capture to the right
			if int(fromSquare)%8 != 7 {
				captureSquare := square(fromSquare - 7)
				if game.position.board.whiteSquares.IsSquareOccupied(captureSquare) {
					appendPawnMove(fromSquare, captureSquare, moves)
				}
			}
		}
	}
}

func appendPawnMove(from square, to square, moves *[]*Move) {
	if to > H7 {
		*moves = append(*moves, &Move{from: from, to: to, promotion: WhiteBishop})
		*moves = append(*moves, &Move{from: from, to: to, promotion: WhiteKnight})
		*moves = append(*moves, &Move{from: from, to: to, promotion: WhiteRook})
		*moves = append(*moves, &Move{from: from, to: to, promotion: WhiteQueen})
	} else if to < A2 {
		*moves = append(*moves, &Move{from: from, to: to, promotion: BlackBishop})
		*moves = append(*moves, &Move{from: from, to: to, promotion: BlackKnight})
		*moves = append(*moves, &Move{from: from, to: to, promotion: BlackRook})
		*moves = append(*moves, &Move{from: from, to: to, promotion: BlackQueen})
	} else {
		*moves = append(*moves, &Move{from: from, to: to})
	}
}

func computeRookMoves(game *Game, moves *[]*Move, ownPieces *Bitboard) {
	var rooks Bitboard
	if game.position.turn == WhiteColor {
		rooks = game.position.board.bbWhiteRook
	} else {
		rooks = game.position.board.bbBlackRook
	}

	for rooks != 0 {
		fromSquare := square(rooks.LeastSignificantBit())
		rooks.ClearLeastSignificantBit()

		blockers := (^game.position.board.emptySquares) & game.precomputedData.RookMasks[fromSquare]

		key := (uint64(blockers) * game.precomputedData.RookMagics[fromSquare]) >> (64 - game.precomputedData.RookIndexBits[fromSquare])

		// Return the preinitialized attack set bitboard from the table
		rookMovesBB := game.precomputedData.RookMoves[fromSquare][key]

		// Remove self-captures
		rookMovesBB &^= *ownPieces

		for rookMovesBB != 0 {
			toSquare := square(rookMovesBB.LeastSignificantBit())
			rookMovesBB.ClearLeastSignificantBit()

			*moves = append(*moves, &Move{from: fromSquare, to: toSquare})
		}
	}
}

func computeBishopMoves(game *Game, moves *[]*Move, ownPieces *Bitboard) {
	var bishops Bitboard
	if game.position.turn == WhiteColor {
		bishops = game.position.board.bbWhiteBishop
	} else {
		bishops = game.position.board.bbBlackBishop
	}

	for bishops != 0 {
		fromSquare := square(bishops.LeastSignificantBit())
		bishops.ClearLeastSignificantBit()

		blockers := (^game.position.board.emptySquares) & game.precomputedData.BishopMasks[fromSquare]

		key := (uint64(blockers) * game.precomputedData.BishopMagics[fromSquare]) >> (64 - game.precomputedData.BishopIndexBits[fromSquare])

		// Return the preinitialized attack set bitboard from the table
		bishopMovesBB := game.precomputedData.BishopMoves[fromSquare][key]

		// Remove self-captures
		bishopMovesBB &^= *ownPieces

		for bishopMovesBB != 0 {
			toSquare := square(bishopMovesBB.LeastSignificantBit())
			bishopMovesBB.ClearLeastSignificantBit()

			*moves = append(*moves, &Move{from: fromSquare, to: toSquare})
		}
	}
}

func computeQueenMoves(game *Game, moves *[]*Move, ownPieces *Bitboard) {
	var queens Bitboard
	if game.position.turn == WhiteColor {
		queens = game.position.board.bbWhiteQueen
	} else {
		queens = game.position.board.bbBlackQueen
	}

	// Compute the moves of each queen by considering both rook and bishop moves
	for queens != 0 {
		fromSquare := square(queens.LeastSignificantBit())
		queens.ClearLeastSignificantBit()

		blockersBishop := (^game.position.board.emptySquares) & game.precomputedData.BishopMasks[fromSquare]
		blockersRook := (^game.position.board.emptySquares) & game.precomputedData.RookMasks[fromSquare]

		keyBishop := (uint64(blockersBishop) * game.precomputedData.BishopMagics[fromSquare]) >> (64 - game.precomputedData.BishopIndexBits[fromSquare])
		keyRook := (uint64(blockersRook) * game.precomputedData.RookMagics[fromSquare]) >> (64 - game.precomputedData.RookIndexBits[fromSquare])

		// Return the preinitialized attack set bitboard from the table
		bishopMovesBB := game.precomputedData.BishopMoves[fromSquare][keyBishop]
		rookMovesBB := game.precomputedData.RookMoves[fromSquare][keyRook]

		queenMovesBB := bishopMovesBB | rookMovesBB

		// Remove self-captures
		queenMovesBB &^= *ownPieces

		for queenMovesBB != 0 {
			toSquare := square(queenMovesBB.LeastSignificantBit())
			queenMovesBB.ClearLeastSignificantBit()

			*moves = append(*moves, &Move{from: fromSquare, to: toSquare})
		}
	}
}
