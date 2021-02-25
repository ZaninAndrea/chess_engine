package main

import (
	. "github.com/ZaninAndrea/chess_engine/chessboard"
)

func generateKing() [64]Bitboard {
	var kingMoves [64]Bitboard

	// iterate the squares in the order they are stored in the bitboard
	for r := 0; r < 8; r++ {
		for f := 0; f < 8; f++ {
			bb := Bitboard(0)

			// Move down
			if r > 0 {
				bb |= SquareFromFileRank(f, r-1).Bitboard()
			}
			// Move up
			if r < 7 {
				bb |= SquareFromFileRank(f, r+1).Bitboard()
			}
			// Move left
			if f > 0 {
				bb |= SquareFromFileRank(f-1, r).Bitboard()
			}
			// Move right
			if f < 7 {
				bb |= SquareFromFileRank(f+1, r).Bitboard()
			}
			// Move down-left
			if r > 0 && f > 0 {
				bb |= SquareFromFileRank(f-1, r-1).Bitboard()
			}
			// Move down-right
			if r > 0 && f < 7 {
				bb |= SquareFromFileRank(f+1, r-1).Bitboard()
			}
			// Move up-left
			if r < 7 && f > 0 {
				bb |= SquareFromFileRank(f-1, r+1).Bitboard()
			}
			// Move up-right
			if r < 7 && f < 7 {
				bb |= SquareFromFileRank(f+1, r+1).Bitboard()
			}

			kingMoves[f+r*8] = bb
		}
	}

	return kingMoves
}
