package chessboard

import (
	"testing"
)

func countAllType1Collisions(game *Game, depth int, zobristTable *[1 << 27]uint64, positionsMap *map[uint64]Position) int {
	collisions := 0
	key := game.position.Hash() >> (64 - 27)

	if zobristTable[key] != 0 && zobristTable[key] == game.position.Hash() &&
		((*positionsMap)[key].board != (*game.position).board ||
			(*positionsMap)[key].castleRights != (*game.position).castleRights ||
			(*positionsMap)[key].turn != (*game.position).turn ||
			(*positionsMap)[key].enPassantSquare != (*game.position).enPassantSquare) {
		collisions = 1
	}

	zobristTable[key] = game.position.Hash()
	(*positionsMap)[key] = *game.position

	if depth == 0 {
		return collisions
	}

	moves := game.LegalMoves()

	for _, move := range moves {
		game.Move(move)
		collisions += countAllType1Collisions(game, depth-1, zobristTable, positionsMap)
		game.UndoMove()
	}

	return collisions
}

func BenchmarkZobristType1Collisions(b *testing.B) {
	// game := NewGame()
	game := NewGameFromFEN("rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8")

	var zobristTable [1 << 27]uint64
	positionsMap := map[uint64]Position{}

	total := countAllType1Collisions(&game, 5, &zobristTable, &positionsMap)

	b.ReportMetric(float64(total), "type1collisions")
}

func countAllType2Collisions(game *Game, depth int, zobristTable *[1 << 27]uint64) int {
	collisions := 0
	key := game.position.Hash() >> (64 - 27)

	if zobristTable[key] != 0 && zobristTable[key] != game.position.Hash() {
		collisions = 1
	}

	zobristTable[key] = game.position.Hash()
	if depth == 0 {
		return collisions
	}

	moves := game.LegalMoves()

	for _, move := range moves {
		game.Move(move)
		collisions += countAllType2Collisions(game, depth-1, zobristTable)
		game.UndoMove()
	}

	return collisions
}

func BenchmarkZobristType2Collisions(b *testing.B) {
	// game := NewGame()
	game := NewGameFromFEN("rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8")
	var zobristTable [1 << 27]uint64

	total := countAllType2Collisions(&game, 5, &zobristTable)

	b.ReportMetric(float64(total), "type2collisions")
}
