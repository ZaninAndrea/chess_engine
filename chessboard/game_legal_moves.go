package chessboard

// LegalMoves returns the legal moves in the current position and caches them
func (game *Game) LegalMoves() []*Move {
	if game.position.legalMoves != nil {
		return game.position.legalMoves
	}
	// The moves list starts empty but we set the capacity to 35
	// to avoid expanding the slice's capacity many times
	pseudolegalMoves := make([]*Move, 0, 35)

	var ownPieces Bitboard
	if game.position.turn == WhiteColor {
		ownPieces = game.position.board.whiteSquares
	} else {
		ownPieces = game.position.board.blackSquares
	}

	computeKingMoves(game, &pseudolegalMoves, &ownPieces)
	computeKnightMoves(game, &pseudolegalMoves, &ownPieces)
	computeRookMoves(game, &pseudolegalMoves, &ownPieces)
	computeBishopMoves(game, &pseudolegalMoves, &ownPieces)
	computeQueenMoves(game, &pseudolegalMoves, &ownPieces)
	computePawnMoves(game, &pseudolegalMoves)

	legalMoves := make([]*Move, 0, len(pseudolegalMoves))
	for i := 0; i < len(pseudolegalMoves); i++ {
		if checkMoveLegality(pseudolegalMoves[i], game) {
			legalMoves = append(legalMoves, pseudolegalMoves[i])
		}
	}

	game.position.legalMoves = legalMoves
	return legalMoves
}

// Takes a pseudolegal move and checks whether is is also legal (will king be checked?)
func checkMoveLegality(move *Move, game *Game) bool {
	simulationBoard := game.position.board
	simulationBoard.Move(*move)

	var kingSquare square
	var enemyKnights Bitboard
	var enemyBishopLikes Bitboard
	var enemyRookLikes Bitboard
	var enemyKing Bitboard

	if game.position.turn == WhiteColor {
		kingSquare = simulationBoard.whiteKingSquare
		enemyKnights = simulationBoard.bbBlackKnight
		enemyBishopLikes = simulationBoard.bbBlackBishop | simulationBoard.bbBlackQueen
		enemyRookLikes = simulationBoard.bbBlackRook | simulationBoard.bbBlackQueen
		enemyKing = simulationBoard.bbBlackKing
	} else {
		kingSquare = simulationBoard.blackKingSquare
		enemyKnights = simulationBoard.bbWhiteKnight
		enemyBishopLikes = simulationBoard.bbWhiteBishop | simulationBoard.bbWhiteQueen
		enemyRookLikes = simulationBoard.bbWhiteRook | simulationBoard.bbWhiteQueen
		enemyKing = simulationBoard.bbWhiteKing
	}

	kingCollisions := game.precomputedData.KingMoves[kingSquare] & enemyKing
	if kingCollisions != 0 {
		return false
	}

	// Simulate putting a knight in the square where the allied king is, if the simulated
	// knight attacks an enemy knight then our king is in check by an enemy knight
	knightCollisions := game.precomputedData.KnightMoves[kingSquare] & enemyKnights
	if knightCollisions != 0 {
		return false
	}

	// Simulate rook and queens moving horizontally/vertically
	blockers := (^simulationBoard.emptySquares) & game.precomputedData.RookMasks[kingSquare]
	key := (uint64(blockers) * game.precomputedData.RookMagics[kingSquare]) >> (64 - game.precomputedData.RookIndexBits[kingSquare])
	rookCollisions := game.precomputedData.RookMoves[kingSquare][key] & enemyRookLikes
	if rookCollisions != 0 {
		return false
	}

	// Simulate bishop and queens moving diagonally
	blockers = (^simulationBoard.emptySquares) & game.precomputedData.BishopMasks[kingSquare]
	key = (uint64(blockers) * game.precomputedData.BishopMagics[kingSquare]) >> (64 - game.precomputedData.BishopIndexBits[kingSquare])
	bishopCollisions := game.precomputedData.BishopMoves[kingSquare][key] & enemyBishopLikes
	if bishopCollisions != 0 {
		return false
	}

	if game.position.turn == WhiteColor {
		upLeftSquare := (kingSquare + 7).Bitboard()
		upRightSquare := (kingSquare + 9).Bitboard()

		if ((simulationBoard.bbBlackPawn & upLeftSquare) != 0) ||
			((simulationBoard.bbBlackPawn & upRightSquare) != 0) {
			return false
		}
	} else {
		downLeftSquare := (kingSquare - 9).Bitboard()
		downRightSquare := (kingSquare - 7).Bitboard()

		if ((simulationBoard.bbWhitePawn & downLeftSquare) != 0) ||
			((simulationBoard.bbWhitePawn & downRightSquare) != 0) {
			return false
		}
	}

	return true
}

func computeKingMoves(game *Game, moves *[]*Move, ownPieces *Bitboard) {
	var kingSquare square
	if game.position.turn == WhiteColor {
		kingSquare = game.position.board.whiteKingSquare
	} else {
		kingSquare = game.position.board.blackKingSquare
	}

	// Get precomputed king moves for that square and remove self-captures
	kingMovesBB := game.precomputedData.KingMoves[kingSquare]
	kingMovesBB &^= *ownPieces

	// Iterating target squares in the bitboard and add moves to the list
	fromSquare := kingSquare
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
