package main

import (
	"fmt"
	"math/rand"

	"github.com/ZaninAndrea/chess_engine/chessboard"
)

func main() {
	game := chessboard.NewGame()

	fmt.Print("\033[H\033[2J")
	fmt.Println(game.Position())
	rand.Seed(5)

	for game.Result() == chessboard.NoResult {
		// time.Sleep(1200 * time.Millisecond)
		moves := game.LegalMoves()
		move := moves[rand.Intn(len(moves))]
		game.Move(move)

		fmt.Print("\033[H\033[2J")
		fmt.Println(game.Position())
	}

	fmt.Println(game.Result())
}
