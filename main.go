package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ZaninAndrea/chess_engine/chessboard"
)

func main() {
	game := chessboard.NewGame()
	// game := chessboard.NewGameFromFEN("r1bqk2r/pp1nbpp1/2p1pn1p/3p4/2PP3B/2NBPN2/PP3PPP/R2QK2R b KQkq - 0 1")
	engine := chessboard.NewBruteForceEngine(&game)

	// fmt.Print("\033[H\033[2J")
	fmt.Println(game.Position())
	rand.Seed(1)

	for game.Result() == chessboard.NoResult {
		time.Sleep(800 * time.Millisecond)
		fmt.Println()
		game.Move(engine.BestMove(60))

		// fmt.Print("\033[H\033[2J")
		fmt.Println(game.Position())
	}

	fmt.Println(game.Result())
}
