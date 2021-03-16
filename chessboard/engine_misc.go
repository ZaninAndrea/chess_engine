package chessboard

import (
	"fmt"

	"github.com/dylhunn/dragontoothmg"
)

// PositionAnalysisString returns a string containing all the factors of the positional analysis
func (eng *BruteForceEngine) PositionAnalysisString() string {
	return fmt.Sprintf("Material difference: %d, Position difference: %d, Center control: %d, Doubled and Isolated Pawns: %d, Passed Pawns: %d",
		materialDifference(eng.trackedGame.position),
		positionDifference(eng.trackedGame.position),
		centerControl(eng.trackedGame.position),
		doubledOrIsolatedPawnsPenalties(eng.trackedGame.position, &eng.trackedGame.precomputedData),
		passedPawnsBonuses(eng.trackedGame.position, &eng.trackedGame.precomputedData),
	)
}

func (eng *BruteForceEngine) sortMoves(moves []dragontoothmg.Move, evaluationTable *ZobristTable) []dragontoothmg.Move {
	return moves
	// scores := make([]int, len(moves))

	// // Fill the scores slice
	// N := len(moves)
	// for i := 0; i < N; i++ {
	// 	undo := eng.game.Apply(moves[i])

	// 	got, hash := evaluationTable.Get(eng.game.position.hash)
	// 	if got {
	// 		scores[i] = hash.Evaluation()
	// 	}

	// 	undo()
	// }

	// // Sort the moves
	// for i := 0; i < len(moves); i++ {
	// 	score := scores[i]
	// 	move := moves[i]

	// 	j := i - 1
	// 	for j >= 0 && scores[j] > score {
	// 		moves[j+1] = moves[j]
	// 		scores[j+1] = scores[j]

	// 		j--
	// 	}

	// 	scores[j+1] = score
	// 	moves[j+1] = move
	// }

	// return moves
}
