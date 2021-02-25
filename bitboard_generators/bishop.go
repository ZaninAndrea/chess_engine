package main

import (
	"fmt"

	. "github.com/ZaninAndrea/chess_engine/chessboard"
)

func generateBishopMasks() [64]Bitboard {
	var maskBBs [64]Bitboard

	mainDiagonal := uint64(0x8040201008040201)
	mainAntiDiagonal := uint64(0x102040810204080)

	fileA := uint64(0x0101010101010101)
	rank1 := uint64(0b11111111)

	outer := uint64(0)
	outer |= rank1
	outer |= rank1 << 56
	outer |= fileA
	outer |= fileA << 7

	inner := ^outer

	for r := 0; r < 8; r++ {
		for f := 0; f < 8; f++ {
			bb := uint64(0)

			// add the diagonals crossing this square
			var diagonal uint64
			if f <= r {
				diagonal = mainDiagonal << ((r - f) * 8)
			} else {
				diagonal = mainDiagonal >> ((f - r) * 8)
			}
			bb |= diagonal

			var antiDiagonal uint64
			if f+r <= 7 {
				antiDiagonal = mainAntiDiagonal >> ((7 - f - r) * 8)
			} else {
				antiDiagonal = mainAntiDiagonal << ((f + r - 7) * 8)
			}
			bb |= antiDiagonal

			// remove the current square and the edges of the board
			squareBB := SquareFromFileRank(f, r).Bitboard()
			bb &= ^uint64(squareBB)
			bb &= inner

			maskBBs[f+r*8] = Bitboard(bb)
		}
	}

	return maskBBs
}

// Generates rook moves for the given magics and checks that there are no collisions
func generateBishopMoves(bishopMagics [64]uint64, bishopIndexBits [64]int) [64][1024]Bitboard {
	var bishopMoves [64][1024]Bitboard

	for r := 0; r < 8; r++ {
		for f := 0; f < 8; f++ {
			fillBishopMovesSquare(f, r, bishopMagics, bishopIndexBits, &bishopMoves)
		}
	}

	return bishopMoves
}

func fillBishopMovesSquare(file int, rank int, bishopMagics [64]uint64, bishopIndexBits [64]int, bishopMoves *[64][1024]Bitboard) {
	square := file + rank*8

	// fill choices with all the squares that can be blockers
	choices := []Bitboard{}

	// left-up squares
	for i := 1; (file-i >= 1) && (rank+i <= 6); i++ {
		choices = append(choices, SquareFromFileRank(file-i, rank+i).Bitboard())
	}
	// left-down squares
	for i := 1; (file-i >= 1) && (rank-i >= 1); i++ {
		choices = append(choices, SquareFromFileRank(file-i, rank-i).Bitboard())
	}
	// right-up squares
	for i := 1; (file+i <= 6) && (rank+i <= 6); i++ {
		choices = append(choices, SquareFromFileRank(file+i, rank+i).Bitboard())
	}
	// right-down squares
	for i := 1; (file+i <= 6) && (rank-i >= 1); i++ {
		choices = append(choices, SquareFromFileRank(file+i, rank-i).Bitboard())
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

		// moving left-up
		for i := 1; (file-i >= 0) && (rank+i <= 7); i++ {
			squareBB := SquareFromFileRank(file-i, rank+i).Bitboard()
			moves |= squareBB

			// if a blocker is encountered exit the loop
			if blockers&squareBB != 0 {
				break
			}
		}
		// moving left-down
		for i := 1; (file-i >= 0) && (rank-i >= 0); i++ {
			squareBB := SquareFromFileRank(file-i, rank-i).Bitboard()
			moves |= squareBB

			// if a blocker is encountered exit the loop
			if blockers&squareBB != 0 {
				break
			}
		}
		// moving right-down
		for i := 1; (file+i <= 7) && (rank-i >= 0); i++ {
			squareBB := SquareFromFileRank(file+i, rank-i).Bitboard()
			moves |= squareBB

			// if a blocker is encountered exit the loop
			if blockers&squareBB != 0 {
				break
			}
		}
		// moving right-up
		for i := 1; (file+i <= 7) && (rank+i <= 7); i++ {
			squareBB := SquareFromFileRank(file+i, rank+i).Bitboard()
			moves |= squareBB

			// if a blocker is encountered exit the loop
			if blockers&squareBB != 0 {
				break
			}
		}

		key := (uint64(blockers) * bishopMagics[square]) >> (64 - bishopIndexBits[square])
		if bishopMoves[square][key] != 0 && bishopMoves[square][key] != moves {
			panic(fmt.Sprintf("Invalid magic number for square %d", square))
		}

		bishopMoves[square][key] = moves
	}
}
