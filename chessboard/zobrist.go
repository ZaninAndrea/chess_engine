package chessboard

import (
	"math/rand"
)

type ZobristHash uint64

var (
	zobristHashWhiteKingCastle  ZobristHash = 0
	zobristHashWhiteQueenCastle ZobristHash = 0
	zobristHashBlackKingCastle  ZobristHash = 0
	zobristHashBlackQueenCastle ZobristHash = 0
	zobristHashBlackTurn        ZobristHash = 0
	zobristHashEnPassant        [8]ZobristHash
	zobristHashMoves            [12][64]ZobristHash
)

func initializeZobristHashes() {
	rand.Seed(1)
	zobristHashWhiteKingCastle = ZobristHash(rand.Uint64() << 21)
	zobristHashWhiteQueenCastle = ZobristHash(rand.Uint64() << 21)
	zobristHashBlackKingCastle = ZobristHash(rand.Uint64() << 21)
	zobristHashBlackQueenCastle = ZobristHash(rand.Uint64() << 21)
	zobristHashBlackTurn = ZobristHash(rand.Uint64() << 21)

	for i := 0; i < 12; i++ {
		for j := 0; j < 64; j++ {
			zobristHashMoves[i][j] = ZobristHash(rand.Uint64() << 21)
		}
	}

	for i := 0; i < 8; i++ {
		zobristHashEnPassant[i] = ZobristHash(rand.Uint64() << 21)
	}
}

// SetData sets the evaluation data in the zobrist hash for a given position
func (h ZobristHash) SetData(evaluation int16, depth int8, lowerBound bool) ZobristHash {
	h |= ZobristHash(evaluation << 7)
	h |= ZobristHash(depth << 1)

	if lowerBound {
		h |= 1
	}

	return h
}

// Key returns the 27 bit key for the hash tables
func (h ZobristHash) Key() int32 {
	return int32(h >> (64 - 27))
}

// PositionHash returns the hash for the current position in long form (43 bits)
func (h ZobristHash) PositionHash() int64 {
	return int64(h >> 21)
}

const evaluationMask ZobristHash = 0b111111111111110000000

// Evaluation returns the evaluation encoded in the zobrist hash
func (h ZobristHash) Evaluation() int {
	return int((h & evaluationMask) >> 6)
}

const depthMask ZobristHash = 0b1111110

// Depth returns the depth of the stored evaluation
func (h ZobristHash) Depth() int {
	return int((h & depthMask) >> 1)
}

// LowerBound returns whether the current hash evaluation is a lower bound (otherwise its the correct value)
func (h ZobristHash) LowerBound() bool {
	return h&1 != 0
}

const zobristCacheSize int = 1 << 27

// ZobristTable is an HashTable using the zobrist hash algorithm
type ZobristTable [zobristCacheSize]ZobristHash

// Get gets an element from the table, returns a boolean value representing whether the
// value was found and the value itself if found
func (tb *ZobristTable) Get(hash ZobristHash) (bool, ZobristHash) {
	if tb[hash.Key()] == 0 {
		return false, 0
	}

	return true, tb[hash.Key()]
}

// Set saves an hash in the table
func (tb *ZobristTable) Set(hash ZobristHash) {
	tb[hash.Key()] = hash
}
