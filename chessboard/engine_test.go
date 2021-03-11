package chessboard

import (
	"fmt"
	"testing"
)

func BenchmarkRandomSelfPlay(b *testing.B) {
	for i := 0; i < b.N; i++ {
		game := NewGame()
		engine := NewRandomEngine(&game)

		fmt.Println(game.Position())

		for game.Result() == NoResult {
			game.Move(engine.BestMove(60))

			fmt.Println()
			fmt.Println(game.Position())
		}

		fmt.Println(game.Result())
	}
}

func BenchmarkBruteForceSelfPlay(b *testing.B) {
	for i := 0; i < b.N; i++ {
		game := NewGame()
		engine := NewBruteForceEngine(&game)

		fmt.Println(game.Position())

		for game.Result() == NoResult {
			game.Move(engine.BestMove(60))

			fmt.Println()
			fmt.Println(game.Position())
		}

		fmt.Println(game.Result())
	}
}

func evaluateAllMoves(eng *BruteForceEngine, depth int) {
	if depth == 0 {
		eng.StaticEvaluation()
		return
	}

	moves := eng.game.LegalMoves()
	for _, move := range moves {
		eng.game.Move(move)
		evaluateAllMoves(eng, depth-1)
		eng.game.UndoMove()
	}
}

func BenchmarkMoveEvaluationStart5Ply(b *testing.B) {
	game := NewGame()
	eng := NewBruteForceEngine(&game)

	for i := 0; i < b.N; i++ {
		evaluateAllMoves(eng, 5)
	}
}

func checkResultAndEvaluateAllMoves(eng *BruteForceEngine, depth int) {
	switch eng.game.Result() {
	case Draw:
		return
	case Checkmate:
		return
	}

	if depth == 0 {
		eng.StaticEvaluation()
		return
	}

	moves := eng.game.LegalMoves()
	for _, move := range moves {
		eng.game.Move(move)
		checkResultAndEvaluateAllMoves(eng, depth-1)
		eng.game.UndoMove()
	}
}

func BenchmarkResultAndEvaluationStart5Ply(b *testing.B) {
	game := NewGame()
	eng := NewBruteForceEngine(&game)

	for i := 0; i < b.N; i++ {
		checkResultAndEvaluateAllMoves(eng, 5)
	}
}

func BenchmarkAnalyzeStart5PlyWithNoPruning(b *testing.B) {
	game := NewGame()
	eng := NewBruteForceEngine(&game)
	eng.AlphaBetaPruningEnabled = false
	eng.TranspositionTableEnabled = false
	eng.AspirationSearchEnabled = false
	eng.QuiescentSearchEnabled = false
	eng.MaxDepth = 5

	for i := 0; i < b.N; i++ {
		eng.BestMove(600000)
	}
}

func BenchmarkAnalyzeStart5PlyWithAllPruning(b *testing.B) {
	game := NewGame()
	eng := NewBruteForceEngine(&game)
	eng.MaxDepth = 5
	eng.QuiescentSearchEnabled = false

	for i := 0; i < b.N; i++ {
		eng.BestMove(600000)
	}
}
