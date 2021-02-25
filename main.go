package main

import (
	"fmt"
	"math/rand"
	"time"

	. "github.com/ZaninAndrea/chess_engine/chessboard"
)

func main() {
	game := NewGame()

	fmt.Print("\033[H\033[2J")
	fmt.Println(game.Position())

	for i := 0; i < 100; i++ {
		time.Sleep(1200 * time.Millisecond)
		moves := game.LegalMoves()

		if len(moves) == 0 {
			fmt.Println("No legal moves")
			break
		}
		move := *(moves[rand.Intn(len(moves))])

		game.Move(move)
		fmt.Print("\033[H\033[2J")
		fmt.Println(game.Position())

	}
}
