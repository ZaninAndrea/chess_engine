package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ZaninAndrea/chess_engine/chessboard"
)

func main() {
	game := chessboard.NewGame()

	fmt.Print("\033[H\033[2J")
	fmt.Println(game.Position())
	rand.Seed(1)

	for game.Result() == chessboard.NoResult {
		time.Sleep(800 * time.Millisecond)
		moves := game.LegalMoves()
		move := moves[rand.Intn(len(moves))]
		game.Move(move)

		fmt.Print("\033[H\033[2J")
		fmt.Println(game.Position())
	}

	fmt.Println(game.Result())
}
