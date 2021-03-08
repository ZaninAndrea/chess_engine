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
	zobristHashWhiteKingCastle = ZobristHash(rand.Uint64())
	zobristHashWhiteQueenCastle = ZobristHash(rand.Uint64())
	zobristHashBlackKingCastle = ZobristHash(rand.Uint64())
	zobristHashBlackQueenCastle = ZobristHash(rand.Uint64())
	zobristHashBlackTurn = ZobristHash(rand.Uint64())

	for i := 0; i < 12; i++ {
		for j := 0; j < 64; j++ {
			zobristHashMoves[i][j] = ZobristHash(rand.Uint64())
		}
	}

	for i := 0; i < 8; i++ {
		zobristHashEnPassant[i] = ZobristHash(rand.Uint64())
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

// The following parameter can be safely decreased down to 21, because it still
// leaves enough space in the hash table content for all the metadata
const zobristCacheSize int = 22

// Key returns the 22 bit key for the hash tables
func (h ZobristHash) Key() int32 {
	return int32(h >> (64 - zobristCacheSize))
}

// HashValue returns the value to store in the hash table as verification
func (h ZobristHash) HashValue() ZobristHash {
	return h << zobristCacheSize
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

// ZobristTable is an HashTable using the zobrist hash algorithm
type ZobristTable [1 << zobristCacheSize]ZobristHash

var zobristCacheHits = 0
var zobristCacheMisses = 0

// Get gets an element from the table, returns a boolean value representing whether the
// value was found and the value itself if found
func (tb *ZobristTable) Get(hash ZobristHash) (bool, ZobristHash) {
	key := hash.Key()
	if tb[key] == 0 {
		zobristCacheMisses++
		return false, 0
	}

	if (tb[key] >> zobristCacheSize) != (hash.HashValue() >> zobristCacheSize) {
		zobristCacheMisses++
		return false, 0
	}

	zobristCacheHits++
	return true, tb[key]
}

// Set saves an hash in the table
func (tb *ZobristTable) Set(key int32, hash ZobristHash) {
	tb[key] = hash
}
