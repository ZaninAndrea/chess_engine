package main

import (
	. "github.com/ZaninAndrea/chess_engine/chessboard"
)

func generateDoubledPawnMasks() ([64]Bitboard, [64]Bitboard) {
	ownColumnMask := [64]Bitboard{}
	sideColumnsMask := [64]Bitboard{}

	fileA := Bitboard(0x101010101010101)

	for file := 0; file < 8; file++ {
		for rank := 0; rank < 8; rank++ {
			sq := SquareFromFileRank(file, rank)

			ownColumn := fileA << (file + rank*8 + 8)
			sideColumns := Bitboard(0)
			if file > 0 {
				sideColumns |= fileA << (file - 1)
			}
			if file < 7 {
				sideColumns |= fileA << (file + 1)
			}

			ownColumnMask[sq] = ownColumn
			sideColumnsMask[sq] = sideColumns
		}
	}

	return ownColumnMask, sideColumnsMask
}

func passedPawnMasks() (whiteMasks [64]Bitboard, blackMasks [64]Bitboard) {
	fileA := Bitboard(0x101010101010101)

	for file := 0; file < 8; file++ {
		for rank := 0; rank < 8; rank++ {
			sq := SquareFromFileRank(file, rank)

			mask := fileA << (file + rank*8 + 8)
			if file > 0 {
				mask |= fileA << (file - 1 + rank*8 + 8)
			}
			if file < 7 {
				mask |= fileA << (file + 1 + rank*8 + 8)
			}

			whiteMasks[sq] = mask
		}
	}

	for file := 0; file < 8; file++ {
		for rank := 0; rank < 8; rank++ {
			sq := SquareFromFileRank(file, rank)

			mask := (fileA << file) >> ((8 - rank) * 8)
			if file > 0 {
				mask |= (fileA << (file - 1)) >> ((8 - rank) * 8)
			}
			if file < 7 {
				mask |= (fileA << (file + 1)) >> ((8 - rank) * 8)
			}

			blackMasks[sq] = mask
		}
	}

	return
}
