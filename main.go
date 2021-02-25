package main

import (
	"fmt"

	. "github.com/ZaninAndrea/chess_engine/chessboard"
)

func main() {
	game := NewGame()

	fmt.Println(game.Position)
}
