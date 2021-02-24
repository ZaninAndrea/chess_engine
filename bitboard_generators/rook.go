package main

import (
	. "github.com/ZaninAndrea/chess_engine/chessboard"
)

func generateRookMasks() [64]Bitboard {
	var maskBBs [64]Bitboard

	fileA := uint64(0x0101010101010101)
	rank1 := uint64(0b11111111)

	for r := 0; r < 8; r++ {
		for f := 0; f < 8; f++ {
			bb := uint64(0)

			// add the current file and rank
			bb |= fileA << f
			bb |= rank1 << uint64(8*r)

			// remove the current square
			squareBB := (Square{File: File(f), Rank: Rank(r)}).Bitboard()
			bb &= ^uint64(squareBB)

			outer := uint64(0)
			if r > 0 {
				outer |= rank1
			}
			if r < 7 {
				outer |= rank1 << 56
			}
			if f > 0 {
				outer |= fileA
			}
			if f < 7 {
				outer |= fileA << 7
			}

			bb &= ^outer

			maskBBs[f+r*8] = Bitboard(bb)
		}
	}

	return maskBBs
}
