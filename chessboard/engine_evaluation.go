package chessboard

func materialDifference(pos *Position) int {
	whitePawnCount := pos.board.bbWhitePawn.PopCount() * 256
	blackPawnCount := pos.board.bbBlackPawn.PopCount() * 256
	whiteKnightCount := pos.board.bbWhiteKnight.PopCount() * 832
	blackKnightCount := pos.board.bbBlackKnight.PopCount() * 832
	whiteBishopCount := pos.board.bbWhiteBishop.PopCount() * 896
	blackBishopCount := pos.board.bbBlackBishop.PopCount() * 896
	whiteRookCount := pos.board.bbWhiteRook.PopCount() * 1280
	blackRookCount := pos.board.bbBlackRook.PopCount() * 1280
	whiteQueenCount := pos.board.bbWhiteQueen.PopCount() * 2496
	blackQueenCount := pos.board.bbBlackQueen.PopCount() * 2496

	whiteTotal := whitePawnCount +
		whiteKnightCount +
		whiteBishopCount +
		whiteRookCount +
		whiteQueenCount
	blackTotal := blackPawnCount +
		blackKnightCount +
		blackBishopCount +
		blackRookCount +
		blackQueenCount

	difference := whiteTotal - blackTotal
	total := whiteTotal + blackTotal
	pieceRatio := (difference * 100) / (total + 1)

	return difference +
		pieceRatio
}

// Values from https://www.chessprogramming.org/Simplified_Evaluation_Function converted to 256th of a pawn
var pawnSquareValue = [64]int{0, 0, 0, 0, 0, 0, 0, 0, 12, 25, 25, -51, -51, 25, 25, 12, 12, -12, -25, 0, 0, -25, -12, 12, 0, 0, 0, 51, 51, 0, 0, 0, 12, 12, 25, 64, 64, 25, 12, 12, 25, 25, 51, 76, 76, 51, 25, 25, 128, 128, 128, 128, 128, 128, 128, 128, 0, 0, 0, 0, 0, 0, 0, 0}
var knightSquareValue = [64]int{-128, -102, -76, -76, -76, -76, -102, -128, -102, -51, 0, 12, 12, 0, -51, -102, -76, 12, 25, 38, 38, 25, 12, -76, -76, 0, 38, 51, 51, 38, 0, -76, -76, 12, 38, 51, 51, 38, 12, -76, -76, 0, 25, 38, 38, 25, 0, -76, -102, -51, 0, 0, 0, 0, -51, -102, -128, -102, -76, -76, -76, -76, -102, -128}
var bishopSquareValue = [64]int{-51, -25, -25, -25, -25, -25, -25, -51, -25, 12, 0, 0, 0, 0, 12, -25, -25, 25, 25, 25, 25, 25, 25, -25, -25, 0, 25, 25, 25, 25, 0, -25, -25, 12, 12, 25, 25, 12, 12, -25, -25, 0, 12, 25, 25, 12, 0, -25, -25, 0, 0, 0, 0, 0, 0, -25, -51, -25, -25, -25, -25, -25, -25, -51}
var rookSquareValue = [64]int{0, 0, 0, 12, 12, 0, 0, 0, -12, 0, 0, 0, 0, 0, 0, -12, -12, 0, 0, 0, 0, 0, 0, -12, -12, 0, 0, 0, 0, 0, 0, -12, -12, 0, 0, 0, 0, 0, 0, -12, -12, 0, 0, 0, 0, 0, 0, -12, 12, 25, 25, 25, 25, 25, 25, 12, 0, 0, 0, 0, 0, 0, 0, 0}
var queenSquareValue = [64]int{-51, -25, -25, -12, -12, -25, -25, -51, -25, 0, 12, 0, 0, 0, 0, -25, -25, 12, 12, 12, 12, 12, 0, -25, 0, 0, 12, 12, 12, 12, 0, -12, -12, 0, 12, 12, 12, 12, 0, -12, -25, 0, 12, 12, 12, 12, 0, -25, -25, 0, 0, 0, 0, 0, 0, -25, -51, -25, -25, -12, -12, -25, -25, -51}
var kingMiddleGameSquareValue = [64]int{51, 76, 25, 0, 0, 25, 76, 51, 51, 51, 0, 0, 0, 0, 51, 51, -25, -51, -51, -51, -51, -51, -51, -25, -51, -76, -76, -102, -102, -76, -76, -51, -76, -102, -102, -128, -128, -102, -102, -76, -76, -102, -102, -128, -128, -102, -102, -76, -76, -102, -102, -128, -128, -102, -102, -76, -76, -102, -102, -128, -128, -102, -102, -76}
var kingEndGameSquareValue = [64]int{-128, -76, -76, -76, -76, -76, -76, -128, -76, -76, 0, 0, 0, 0, -76, -76, -76, -25, 51, 76, 76, 51, -25, -76, -76, -25, 76, 102, 102, 76, -25, -76, -76, -25, 76, 102, 102, 76, -25, -76, -76, -25, 51, 76, 76, 51, -25, -76, -76, -51, -25, 0, 0, -25, -51, -76, -128, -102, -76, -51, -51, -76, -102, -128}

func positionDifference(pos *Position) int {
	score := 0

	score += pieceSquareValues(pos.board.bbWhitePawn, &pawnSquareValue, true) -
		pieceSquareValues(pos.board.bbBlackPawn, &pawnSquareValue, false)
	score += pieceSquareValues(pos.board.bbWhiteKnight, &knightSquareValue, true) -
		pieceSquareValues(pos.board.bbBlackKnight, &knightSquareValue, false)
	score += pieceSquareValues(pos.board.bbWhiteBishop, &bishopSquareValue, true) -
		pieceSquareValues(pos.board.bbBlackBishop, &bishopSquareValue, false)
	score += pieceSquareValues(pos.board.bbWhiteRook, &rookSquareValue, true) -
		pieceSquareValues(pos.board.bbBlackRook, &rookSquareValue, false)
	score += pieceSquareValues(pos.board.bbWhiteQueen, &queenSquareValue, true) -
		pieceSquareValues(pos.board.bbBlackQueen, &queenSquareValue, false)
	score += pieceSquareValues(pos.board.bbWhiteKing, &kingMiddleGameSquareValue, true) -
		pieceSquareValues(pos.board.bbBlackKing, &kingMiddleGameSquareValue, false)

	return score
}

func pieceSquareValues(piece Bitboard, values *[64]int, isWhite bool) int {
	score := 0

	if isWhite {
		for piece != 0 {
			sq := piece.LeastSignificant1Bit()
			piece.ClearLeastSignificant1Bit()

			score += (*values)[sq]
		}
	} else {
		for piece != 0 {
			sq := piece.LeastSignificant1Bit()
			piece.ClearLeastSignificant1Bit()

			invertedSquare := SquareFromFileRank(sq%8, 7-(sq/8))
			score += (*values)[invertedSquare]
		}
	}

	return score
}

var centerBitboard = D4.Bitboard() | D5.Bitboard() | E4.Bitboard() | E5.Bitboard()

func centerControl(pos *Position) int {
	score := 0

	if pos.board.bbWhitePawn&centerBitboard != 0 {
		score += 70
	}
	if pos.board.bbBlackPawn&centerBitboard != 0 {
		score -= 70
	}

	return score
}

// Computes penalties for
// - doubled pawns, that is when there are 2 same color pawns in the same file
//   and in the neighbouring files there are no pawns of that same color
// - isolated pawn, that is a file with a single pawn and without pawns in the neighbouring files
func doubledOrIsolatedPawnsPenalties(pos *Position, precomputedData *PrecomputedData) int {
	score := 0

	whitePawns := pos.board.bbWhitePawn
	for whitePawns != 0 {
		sq := whitePawns.LeastSignificant1Bit()
		whitePawns.ClearLeastSignificant1Bit()

		if precomputedData.DoublePawnsSidesMasks[sq]&pos.board.bbWhitePawn == 0 {
			if precomputedData.DoublePawnsForwardMasks[sq]&pos.board.bbWhitePawn != 0 {
				// doubled isolated pawns
				score -= 90
			} else {
				// isolated pawn
				score -= 45
			}
		} else if precomputedData.DoublePawnsForwardMasks[sq]&pos.board.bbWhitePawn != 0 {
			// doubled pawn with neighbouring pawn
			score -= 45
		}

	}

	blackPawns := pos.board.bbBlackPawn
	for blackPawns != 0 {
		sq := blackPawns.LeastSignificant1Bit()
		blackPawns.ClearLeastSignificant1Bit()

		if precomputedData.DoublePawnsSidesMasks[sq]&pos.board.bbBlackPawn == 0 {
			if precomputedData.DoublePawnsForwardMasks[sq]&pos.board.bbBlackPawn != 0 {
				// doubled pawns
				score += 90
			} else {
				// isolated pawn
				score += 45
			}
		} else if precomputedData.DoublePawnsForwardMasks[sq]&pos.board.bbBlackPawn != 0 {
			// doubled pawn with neighbouring pawn
			score += 45
		}
	}

	return score
}

// Compute bonuses for passed pawns
func passedPawnsBonuses(pos *Position, precomputedData *PrecomputedData) int {
	score := 0

	whitePawns := pos.board.bbWhitePawn
	for whitePawns != 0 {
		sq := whitePawns.LeastSignificant1Bit()
		whitePawns.ClearLeastSignificant1Bit()

		if precomputedData.PassedPawnWhiteMasks[sq]&pos.board.bbBlackPawn == 0 {
			score += 120
		}
	}

	blackPawns := pos.board.bbBlackPawn
	for blackPawns != 0 {
		sq := blackPawns.LeastSignificant1Bit()
		blackPawns.ClearLeastSignificant1Bit()

		if precomputedData.PassedPawnBlackMasks[sq]&pos.board.bbWhitePawn == 0 {
			score -= 120
		}
	}

	return score
}

// StaticEvaluation returns an evaluation of the current position from a
// strategic standpoint (e.g. material imbalances, pawn structures, ...) without
// considering any tactical advantages (e.g. ability to capture a piece)
func (eng *BruteForceEngine) StaticEvaluation() int {
	score := 0
	if eng.MaterialDifferenceEval {
		score += materialDifference(eng.game.position)
	}
	if eng.PositionDifferenceEval {
		score += positionDifference(eng.game.position)
	}
	if eng.CenterControlEval {
		score += centerControl(eng.game.position)
	}
	if eng.DoubledIsolatedPawnsEval {
		score += doubledOrIsolatedPawnsPenalties(eng.game.position, &eng.game.precomputedData)
	}
	if eng.PassedPawnsEval {
		score += passedPawnsBonuses(eng.game.position, &eng.game.precomputedData)
	}

	// Stabilizes fluctuations between even and odd depth evaluations
	score += int(eng.game.position.turn) * 15

	return score * int(eng.game.position.turn)
}
