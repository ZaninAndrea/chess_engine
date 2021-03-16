package chessboard

import (
	"testing"
)

func countAllMoves(game *Game, depth int) int {
	if depth == 0 {
		return 1
	}

	total := 0
	moves := game.LegalMoves()

	for _, move := range moves {
		game.Move(move)
		total += countAllMoves(game, depth-1)
		game.UndoMove()

	}

	return total
}

func TestMoveGenerationStart5Ply(t *testing.T) {
	game := NewGame()
	total := countAllMoves(&game, 5)

	if total != 4_865_609 {
		t.Errorf("Total moves up to 5 ply should be 4.865.609, %d was returned instead", total)
	}
}

func TestMoveGenerationOther5Ply(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
		return
	}

	game := NewGameFromFEN("rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8")
	total := countAllMoves(&game, 5)

	if total != 89_941_194 {
		t.Errorf("Total moves up to 5 ply should be 89.941.194, %d was returned instead", total)
	}
}

func BenchmarkMoveGeneration6Ply(b *testing.B) {
	for i := 0; i < b.N; i++ {
		game := NewGame()
		total := countAllMoves(&game, 6)
		b.Log(total)
	}
}

func BenchmarkStaticMoveGeneration(b *testing.B) {
	game := NewGame()

	total := 0
	for i := 0; i < b.N; i++ {
		game.position.legalMoves = nil
		total += len(game.LegalMoves())
	}

	b.ReportMetric(float64(total), "total")
}

func BenchmarkStaticQueenMovesGeneration(b *testing.B) {
	game := NewGame()

	for i := 0; i < b.N; i++ {
		_moves := make([]Move, 0, 35)
		moves := &_moves
		ownPieces := &game.position.board.whiteSquares

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

				*moves = append(*moves, *NewMove(fromSquare, toSquare, NoPiece, NoFlag))
			}
		}
	}
}
