package chessboard

// Engine is the common interface for all the implemented AIs
type Engine interface {
	BestMove(remainingTime int) *Move
}

func materialDifference(pos *Position) int {
	pawnDifference := pos.board.bbWhitePawn.PopCount() - pos.board.bbBlackPawn.PopCount()
	knightDifference := pos.board.bbWhiteKnight.PopCount() - pos.board.bbBlackKnight.PopCount()
	bishopDifference := pos.board.bbWhiteBishop.PopCount() - pos.board.bbBlackBishop.PopCount()
	rookDifference := pos.board.bbWhiteRook.PopCount() - pos.board.bbBlackRook.PopCount()
	queenDifference := pos.board.bbWhiteQueen.PopCount() - pos.board.bbBlackQueen.PopCount()

	return (pawnDifference * 256) +
		(knightDifference * 832) +
		(bishopDifference * 896) +
		(rookDifference * 1280) +
		(queenDifference * 2496)
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
