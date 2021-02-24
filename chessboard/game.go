package chessboard

import (
	"encoding/json"
	"io/ioutil"
)

type PrecomputedData struct {
	KingMoves       [64]Bitboard
	KnightMoves     [64]Bitboard
	RookMagics      [64]uint64
	RookIndexBits   [64]int
	RookMasks       [64]Bitboard
	RookMoves       [64][4096]Bitboard
	BishopMagics    [64]uint64
	BishopIndexBits [64]int
	BishopMasks     [64]Bitboard
	BishopMoves     [64][1024]Bitboard
}

type Game struct {
	PrecomputedData PrecomputedData
	Position        Position
	History         []Position
	Moves           []Move
}

// LoadPrecomputedData loads all the precomputed data for fast move generation
func (game *Game) LoadPrecomputedData(path string) {
	jsonBytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var data PrecomputedData
	err = json.Unmarshal(jsonBytes, &data)
	if err != nil {
		panic(err)
	}

	game.PrecomputedData = data
}
