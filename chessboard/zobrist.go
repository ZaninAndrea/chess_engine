package chessboard

import (
	"math/rand"
)

var (
	zobristHashWhiteKingCastle  uint64 = 0
	zobristHashWhiteQueenCastle uint64 = 0
	zobristHashBlackKingCastle  uint64 = 0
	zobristHashBlackQueenCastle uint64 = 0
	zobristHashBlackTurn        uint64 = 0
	zobristHashEnPassant        [8]uint64
	zobristHashMoves            [12][64]uint64
)

func initializeZobristHashes() {
	rand.Seed(1)
	zobristHashWhiteKingCastle = rand.Uint64() << 20
	zobristHashWhiteQueenCastle = rand.Uint64() << 20
	zobristHashBlackKingCastle = rand.Uint64() << 20
	zobristHashBlackQueenCastle = rand.Uint64() << 20
	zobristHashBlackTurn = rand.Uint64() << 20

	for i := 0; i < 12; i++ {
		for j := 0; j < 64; j++ {
			zobristHashMoves[i][j] = rand.Uint64() << 20
		}
	}

	for i := 0; i < 8; i++ {
		zobristHashEnPassant[i] = rand.Uint64() << 20
	}

}
