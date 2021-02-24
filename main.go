package main

import (
	. "github.com/ZaninAndrea/chess_engine/chessboard"
)

func main() {
	game := Game{}

	game.LoadPrecomputedData("./precomputed.json")
}
