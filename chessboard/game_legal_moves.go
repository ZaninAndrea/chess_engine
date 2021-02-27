package chessboard

// LegalMoves returns the legal moves in the current position and caches them
func (game *Game) LegalMoves() []*Move {
	// Return cached value if possible
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
			// Captures should reset the half move clock
			if pseudolegalMoves[i].to.Bitboard()&game.position.board.emptySquares == 0 {
				pseudolegalMoves[i].flags |= ResetHalfMoveClockFlag
			}

			legalMoves = append(legalMoves, pseudolegalMoves[i])
		}
	}

	game.position.legalMoves = legalMoves
	return legalMoves
}

// Takes a pseudolegal move and checks whether is is also legal (will king be checked?)
func checkMoveLegality(move *Move, game *Game) bool {
	simulationBoard := game.position.board
	simulationBoard.Move(move)

	var kingSquare square
	if game.position.turn == WhiteColor {
		kingSquare = simulationBoard.whiteKingSquare
	} else {
		kingSquare = simulationBoard.blackKingSquare
	}

	return !simulationBoard.IsUnderAttack(game, kingSquare)
}

// Bitboards for the squares that must be empty in order to castle
const InBetweenWhiteKingCastle = Bitboard(96)
const InBetweenWhiteQueenCastle = Bitboard(14)
const InBetweenBlackKingCastle = Bitboard(6917529027641081856)
const InBetweenBlackQueenCastle = Bitboard(1008806316530991104)

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
		toSquare := square(kingMovesBB.LeastSignificant1Bit())
		*moves = append(*moves, &Move{from: fromSquare, to: toSquare})

		kingMovesBB.ClearLeastSignificant1Bit()
	}

	// Check castling conditions: still have rights to castle, squares between rook and king are free,
	// not under check, king will not move in or through an attacked position
	if game.position.turn == WhiteColor {
		if game.position.castleRights.WhiteKingSide &&
			(InBetweenWhiteKingCastle&game.position.board.emptySquares == InBetweenWhiteKingCastle) &&
			!game.position.board.IsUnderAttack(game, E1) &&
			!game.position.board.IsUnderAttack(game, F1) &&
			!game.position.board.IsUnderAttack(game, G1) {
			*moves = append(*moves, &Move{from: E1, to: G1, flags: WhiteKingCastleFlag})
		}

		if game.position.castleRights.WhiteQueenSide &&
			(InBetweenWhiteQueenCastle&game.position.board.emptySquares == InBetweenWhiteQueenCastle) &&
			!game.position.board.IsUnderAttack(game, E1) &&
			!game.position.board.IsUnderAttack(game, D1) &&
			!game.position.board.IsUnderAttack(game, C1) {
			*moves = append(*moves, &Move{from: E1, to: C1, flags: WhiteQueenCastleFlag})
		}
	} else {
		if game.position.castleRights.BlackKingSide &&
			(InBetweenBlackKingCastle&game.position.board.emptySquares == InBetweenBlackKingCastle) &&
			!game.position.board.IsUnderAttack(game, E8) &&
			!game.position.board.IsUnderAttack(game, F8) &&
			!game.position.board.IsUnderAttack(game, G8) {
			*moves = append(*moves, &Move{from: E8, to: G8, flags: BlackKingCastleFlag})
		}

		if game.position.castleRights.BlackQueenSide &&
			(InBetweenBlackQueenCastle&game.position.board.emptySquares == InBetweenBlackQueenCastle) &&
			!game.position.board.IsUnderAttack(game, E8) &&
			!game.position.board.IsUnderAttack(game, D8) &&
			!game.position.board.IsUnderAttack(game, C8) {
			*moves = append(*moves, &Move{from: E8, to: C8, flags: BlackQueenCastleFlag})
		}

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
		fromSquare := square(knights.LeastSignificant1Bit())
		knights.ClearLeastSignificant1Bit()

		// Get precomputed king moves for that square and remove self-captures
		knightMovesBB := game.precomputedData.KnightMoves[fromSquare]
		knightMovesBB &^= *ownPieces

		for knightMovesBB != 0 {
			toSquare := square(knightMovesBB.LeastSignificant1Bit())
			knightMovesBB.ClearLeastSignificant1Bit()

			*moves = append(*moves, &Move{from: fromSquare, to: toSquare})
		}
	}
}

func computePawnMoves(game *Game, moves *[]*Move) {
	if game.position.turn == WhiteColor {
		pawns := game.position.board.bbWhitePawn

		// iterate all white pawns on the board
		for pawns != 0 {
			fromSquare := square(pawns.LeastSignificant1Bit())
			pawns.ClearLeastSignificant1Bit()

			// Move forward
			forwardSquare := square(fromSquare + 8)
			if game.position.board.emptySquares.IsSquareOccupied(forwardSquare) {
				appendPawnMove(fromSquare, forwardSquare, moves)

				// Move forward by two squares
				if fromSquare < A3 {
					forwardSquare := square(fromSquare + 16)
					if game.position.board.emptySquares.IsSquareOccupied(forwardSquare) {
						*moves = append(*moves, &Move{
							from:  fromSquare,
							to:    forwardSquare,
							flags: ResetHalfMoveClockFlag | DoublePawnPushFlag})
					}
				}
			}

			// Capture to the left
			if int(fromSquare)%8 != 0 {
				captureSquare := square(fromSquare + 7)
				if game.position.board.blackSquares.IsSquareOccupied(captureSquare) {
					appendPawnMove(fromSquare, captureSquare, moves)
				} else if game.position.enPassantSquare == captureSquare {
					*moves = append(*moves, &Move{
						from:  fromSquare,
						to:    captureSquare,
						flags: ResetHalfMoveClockFlag | WhiteEnPassantFlag})
				}
			}

			// Capture to the right
			if int(fromSquare)%8 != 7 {
				captureSquare := square(fromSquare + 9)
				if game.position.board.blackSquares.IsSquareOccupied(captureSquare) {
					appendPawnMove(fromSquare, captureSquare, moves)
				} else if game.position.enPassantSquare == captureSquare {
					*moves = append(*moves, &Move{
						from:  fromSquare,
						to:    captureSquare,
						flags: ResetHalfMoveClockFlag | WhiteEnPassantFlag})
				}
			}
		}
	} else {
		pawns := game.position.board.bbBlackPawn

		// iterate all black pawns on the board
		for pawns != 0 {
			fromSquare := square(pawns.LeastSignificant1Bit())
			pawns.ClearLeastSignificant1Bit()

			// Move forward
			forwardSquare := square(fromSquare - 8)
			if game.position.board.emptySquares.IsSquareOccupied(forwardSquare) {
				appendPawnMove(fromSquare, forwardSquare, moves)

				// Move forward by two squares
				if fromSquare > H6 {
					forwardSquare := square(fromSquare - 16)
					if game.position.board.emptySquares.IsSquareOccupied(forwardSquare) {
						*moves = append(*moves, &Move{
							from:  fromSquare,
							to:    forwardSquare,
							flags: ResetHalfMoveClockFlag | DoublePawnPushFlag})
					}
				}
			}

			// Capture to the left
			if int(fromSquare)%8 != 0 {
				captureSquare := square(fromSquare - 9)
				if game.position.board.whiteSquares.IsSquareOccupied(captureSquare) {
					appendPawnMove(fromSquare, captureSquare, moves)
				} else if game.position.enPassantSquare == captureSquare {
					*moves = append(*moves, &Move{
						from:  fromSquare,
						to:    captureSquare,
						flags: ResetHalfMoveClockFlag | BlackEnPassantFlag})
				}
			}

			// Capture to the right
			if int(fromSquare)%8 != 7 {
				captureSquare := square(fromSquare - 7)
				if game.position.board.whiteSquares.IsSquareOccupied(captureSquare) {
					appendPawnMove(fromSquare, captureSquare, moves)
				} else if game.position.enPassantSquare == captureSquare {
					*moves = append(*moves, &Move{
						from:  fromSquare,
						to:    captureSquare,
						flags: ResetHalfMoveClockFlag | BlackEnPassantFlag})
				}
			}
		}
	}
}

func appendPawnMove(from square, to square, moves *[]*Move) {
	if to > H7 {
		*moves = append(*moves, &Move{from: from, to: to, promotion: WhiteBishop, flags: ResetHalfMoveClockFlag})
		*moves = append(*moves, &Move{from: from, to: to, promotion: WhiteKnight, flags: ResetHalfMoveClockFlag})
		*moves = append(*moves, &Move{from: from, to: to, promotion: WhiteRook, flags: ResetHalfMoveClockFlag})
		*moves = append(*moves, &Move{from: from, to: to, promotion: WhiteQueen, flags: ResetHalfMoveClockFlag})
	} else if to < A2 {
		*moves = append(*moves, &Move{from: from, to: to, promotion: BlackBishop, flags: ResetHalfMoveClockFlag})
		*moves = append(*moves, &Move{from: from, to: to, promotion: BlackKnight, flags: ResetHalfMoveClockFlag})
		*moves = append(*moves, &Move{from: from, to: to, promotion: BlackRook, flags: ResetHalfMoveClockFlag})
		*moves = append(*moves, &Move{from: from, to: to, promotion: BlackQueen, flags: ResetHalfMoveClockFlag})
	} else {
		*moves = append(*moves, &Move{from: from, to: to, flags: ResetHalfMoveClockFlag})
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
		fromSquare := square(rooks.LeastSignificant1Bit())
		rooks.ClearLeastSignificant1Bit()

		blockers := (^game.position.board.emptySquares) & game.precomputedData.RookMasks[fromSquare]

		key := (uint64(blockers) * game.precomputedData.RookMagics[fromSquare]) >> (64 - game.precomputedData.RookIndexBits[fromSquare])

		// Return the preinitialized attack set bitboard from the table
		rookMovesBB := game.precomputedData.RookMoves[fromSquare][key]

		// Remove self-captures
		rookMovesBB &^= *ownPieces

		for rookMovesBB != 0 {
			toSquare := square(rookMovesBB.LeastSignificant1Bit())
			rookMovesBB.ClearLeastSignificant1Bit()

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
		fromSquare := square(bishops.LeastSignificant1Bit())
		bishops.ClearLeastSignificant1Bit()

		blockers := (^game.position.board.emptySquares) & game.precomputedData.BishopMasks[fromSquare]

		key := (uint64(blockers) * game.precomputedData.BishopMagics[fromSquare]) >> (64 - game.precomputedData.BishopIndexBits[fromSquare])

		// Return the preinitialized attack set bitboard from the table
		bishopMovesBB := game.precomputedData.BishopMoves[fromSquare][key]

		// Remove self-captures
		bishopMovesBB &^= *ownPieces

		for bishopMovesBB != 0 {
			toSquare := square(bishopMovesBB.LeastSignificant1Bit())
			bishopMovesBB.ClearLeastSignificant1Bit()

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
		fromSquare := square(queens.LeastSignificant1Bit())
		queens.ClearLeastSignificant1Bit()

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
			toSquare := square(queenMovesBB.LeastSignificant1Bit())
			queenMovesBB.ClearLeastSignificant1Bit()

			*moves = append(*moves, &Move{from: fromSquare, to: toSquare})
		}
	}
}
