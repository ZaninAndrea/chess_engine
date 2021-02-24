package main

import (
	. "github.com/ZaninAndrea/chess_engine/chessboard"
)

func generateKnight() [64]Bitboard {
	var knightMoves [64]Bitboard

	// iterate the squares in the order they are stored in the bitboard
	for r := 0; r < 8; r++ {
		for f := 0; f < 8; f++ {
			bb := Bitboard(0)

			// UP LEFT
			if r < 6 && f > 0 {
				bb |= (Square{File: File(f - 1), Rank: Rank(r + 2)}).Bitboard()
			}
			// UP RIGHT
			if r < 6 && f < 7 {
				bb |= (Square{File: File(f + 1), Rank: Rank(r + 2)}).Bitboard()
			}
			// RIGHT UP
			if r < 7 && f < 6 {
				bb |= (Square{File: File(f + 2), Rank: Rank(r + 1)}).Bitboard()
			}
			// RIGHT DOWN
			if r > 0 && f < 6 {
				bb |= (Square{File: File(f + 2), Rank: Rank(r - 1)}).Bitboard()
			}
			// DOWN RIGHT
			if r > 1 && f < 7 {
				bb |= (Square{File: File(f + 1), Rank: Rank(r - 2)}).Bitboard()
			}
			// DOWN LEFT
			if r > 1 && f > 0 {
				bb |= (Square{File: File(f - 1), Rank: Rank(r - 2)}).Bitboard()
			}
			// LEFT DOWN
			if r > 0 && f > 1 {
				bb |= (Square{File: File(f - 2), Rank: Rank(r - 1)}).Bitboard()
			}
			// LEFT UP
			if r < 7 && f > 1 {
				bb |= (Square{File: File(f - 2), Rank: Rank(r + 1)}).Bitboard()
			}

			knightMoves[f+r*8] = bb
		}
	}

	return knightMoves
}
