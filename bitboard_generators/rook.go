package main

import (
	"fmt"

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
			squareBB := SquareFromFileRank(f, r).Bitboard()
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

// Generates rook moves for the given magics and checks that there are no collisions
func generateRookMoves(rookMagics [64]uint64, rookIndexBits [64]int) [64][4096]Bitboard {
	var rookMoves [64][4096]Bitboard

	for r := 0; r < 8; r++ {
		for f := 0; f < 8; f++ {
			fillRookMovesSquare(f, r, rookMagics, rookIndexBits, &rookMoves)
		}
	}

	return rookMoves
}

func fillRookMovesSquare(file int, rank int, rookMagics [64]uint64, rookIndexBits [64]int, rookMoves *[64][4096]Bitboard) {
	square := file + rank*8

	// fill choices with all the squares that can be blockers
	choices := []Bitboard{}

	// left squares
	for f := 1; f < file; f++ {
		choices = append(choices, SquareFromFileRank(f, rank).Bitboard())
	}
	// right squares
	for f := 6; f > file; f-- {
		choices = append(choices, SquareFromFileRank(f, rank).Bitboard())
	}
	// up squares
	for r := 6; r > rank; r-- {
		choices = append(choices, SquareFromFileRank(file, r).Bitboard())
	}
	// down squares
	for r := 1; r < rank; r++ {
		choices = append(choices, SquareFromFileRank(file, r).Bitboard())
	}

	// All the combinations of blocked squares can be iterated by
	// counting up to 2^len(choices) and parsing the bits as blocked squares
	combinations := 1 << len(choices)
	for blockedPieces := int(0); blockedPieces < combinations; blockedPieces++ {
		blockers := Bitboard(0)

		// fill the map of blockers
		for i := 0; i < len(choices); i++ {
			if blockedPieces&(1<<i) != 0 {
				blockers |= choices[i]
			}
		}

		// compute allowed moves
		moves := Bitboard(0)

		// moving left
		for f := file - 1; f >= 0; f-- {
			squareBB := SquareFromFileRank(f, rank).Bitboard()
			moves |= squareBB

			// if a blocker is encountered exit the loop
			if blockers&squareBB != 0 {
				break
			}
		}
		// moving right
		for f := file + 1; f <= 7; f++ {
			squareBB := SquareFromFileRank(f, rank).Bitboard()
			moves |= squareBB

			// if a blocker is encountered exit the loop
			if blockers&squareBB != 0 {
				break
			}
		}
		// moving up
		for r := rank + 1; r <= 7; r++ {
			squareBB := SquareFromFileRank(file, r).Bitboard()
			moves |= squareBB

			// if a blocker is encountered exit the loop
			if blockers&squareBB != 0 {
				break
			}
		}
		// moving down
		for r := rank - 1; r >= 0; r-- {
			squareBB := SquareFromFileRank(file, r).Bitboard()
			moves |= squareBB

			// if a blocker is encountered exit the loop
			if blockers&squareBB != 0 {
				break
			}
		}

		key := (uint64(blockers) * rookMagics[square]) >> (64 - rookIndexBits[square])
		if rookMoves[square][key] != 0 && rookMoves[square][key] != moves {
			panic(fmt.Sprintf("Invalid magic number for square %d", square))
		}

		rookMoves[square][key] = moves
	}
}
