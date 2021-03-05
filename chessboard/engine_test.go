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
