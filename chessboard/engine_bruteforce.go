package chessboard

import (
	"fmt"
	"time"
)

// BruteForceEngine explores all the tree to find the best move
type BruteForceEngine struct {
	trackedGame               *Game
	game                      Game
	MaterialDifferenceEval    bool
	PositionDifferenceEval    bool
	CenterControlEval         bool
	DoubledIsolatedPawnsEval  bool
	PassedPawnsEval           bool
	QuiescentSearchEnabled    bool
	AlphaBetaPruningEnabled   bool
	TranspositionTableEnabled bool
	MoveSortingEnabled        bool
}

// NewBruteForceEngine initializes a BruteForceEngine
func NewBruteForceEngine(game *Game) Engine {
	return &BruteForceEngine{trackedGame: game,
		MaterialDifferenceEval:    true,
		PositionDifferenceEval:    true,
		QuiescentSearchEnabled:    true,
		AlphaBetaPruningEnabled:   true,
		CenterControlEval:         true,
		TranspositionTableEnabled: true,
		MoveSortingEnabled:        true,
		DoubledIsolatedPawnsEval:  true,
		PassedPawnsEval:           true,
	}
}

var nodes = 0

// BestMove returns the best move as computed by the AI
func (eng *BruteForceEngine) BestMove(remainingTime int) *Move {
	endTime := time.Now().Add(time.Duration(remainingTime) * (time.Second / 40))
	eng.game = *eng.trackedGame

	nodes = 0
	zobristCacheHits = 0
	zobristCacheMisses = 0

	move := eng.game.LegalMoves()[0]
	evaluations := &(ZobristTable{})
	quiescentEvaluations := &(ZobristTable{})
	depth := 1
	for {
		aborted, bestMove, newEvals, newQuiescentEvals := eng.NegaMax(depth, endTime, evaluations, quiescentEvaluations)
		evaluations = newEvals
		quiescentEvaluations = newQuiescentEvals

		if aborted {
			break
		}

		move = bestMove
		depth++
	}

	fmt.Printf("Depth reached: %d, Nodes explored: %d, Cache hits: %d, Cache misses: %d\n", depth-1, nodes, zobristCacheHits, zobristCacheMisses)

	return move
}

func (eng *BruteForceEngine) sortMoves(moves []*Move, evaluationTable *ZobristTable) []*Move {
	scores := make([]int, len(moves))

	// Fill the scores slice
	N := len(moves)
	for i := 0; i < N; i++ {
		eng.game.Move(moves[i])

		got, hash := evaluationTable.Get(eng.game.position.hash)
		if got {
			scores[i] = hash.Evaluation()
		}

		eng.game.UndoMove()
	}

	// A quicksort implementation that sorts moves by their score
	// var recQuickSort func(int, int)
	// recQuickSort = func(low int, high int) {
	// 	if low < high {
	// 		pivot := low
	// 		pivotVal := scores[low]

	// 		for i := low + 1; i <= high; i++ {
	// 			if scores[i] < pivotVal {
	// 				pivot++

	// 				// swap moves at i and pivot
	// 				swap := moves[i]
	// 				moves[i] = moves[pivot]
	// 				moves[pivot] = swap

	// 				// swap scores at i and pivot
	// 				swapScore := scores[i]
	// 				scores[i] = scores[pivot]
	// 				scores[pivot] = swapScore
	// 			}
	// 		}

	// 		// Swap scores and moves at i and pivot
	// 		scores[low] = scores[pivot]
	// 		scores[pivot] = pivotVal
	// 		swap := moves[low]
	// 		moves[low] = moves[pivot]
	// 		moves[pivot] = swap

	// 		recQuickSort(low, pivot-1)
	// 		recQuickSort(pivot+1, high)
	// 	}
	// }
	// recQuickSort(0, len(moves)-1)

	// An alternative implementation with insertion sort
	for i := 0; i < len(moves); i++ {
		score := scores[i]
		move := moves[i]

		j := i - 1
		for j >= 0 && scores[j] > score {
			moves[j+1] = moves[j]
			scores[j+1] = scores[j]

			j--
		}

		scores[j+1] = score
		moves[j+1] = move
	}

	return moves
}

// NegaMax does a negamax search of the tree up to the passed depth
func (eng *BruteForceEngine) NegaMax(depth int, endTime time.Time, evaluations *ZobristTable, quiescentEvaluations *ZobristTable) (bool, *Move, *ZobristTable, *ZobristTable) {
	var legalMoves []*Move

	if eng.MoveSortingEnabled {
		legalMoves = eng.sortMoves(eng.game.LegalMoves(), evaluations)
	} else {
		legalMoves = eng.game.LegalMoves()
	}

	bestMove := legalMoves[0]
	bestScore := -1_000_000
	mainLine := []*Move{}

	for i := 0; i < len(legalMoves); i++ {
		// Abort search if running out of time
		if time.Now().After(endTime) {
			return true, nil, nil, nil
		}

		eng.game.Move(legalMoves[i])

		score, variation := eng.recNegaMax(depth-1, -1_000_000, -bestScore, evaluations, quiescentEvaluations)
		score = -score
		if score > bestScore {
			bestScore = score
			bestMove = legalMoves[i]
			mainLine = append(variation, legalMoves[i])
		}

		eng.game.UndoMove()
	}

	fmt.Printf("Score: %d, ", bestScore)
	fmt.Print("Main line: ")
	for i := len(mainLine) - 1; i >= 0; i-- {
		fmt.Printf("%s ", *mainLine[i])
	}
	fmt.Println()

	return false, bestMove, evaluations, quiescentEvaluations
}

func (eng *BruteForceEngine) recNegaMax(depth int, alpha int, beta int, evaluationCache *ZobristTable, quiescentCache *ZobristTable) (int, []*Move) {
	nodes++

	switch eng.game.Result() {
	case Draw:
		return 0, []*Move{}
	case Checkmate:
		return -1_000_000, []*Move{}
	}

	if depth == 0 {
		if eng.QuiescentSearchEnabled {
			nodes--
			return eng.quiescentSearch(7, alpha, beta, evaluationCache, quiescentCache)
		}

		return eng.StaticEvaluation(), []*Move{}
	}

	bestScore := -1_000_000
	mainLine := []*Move{}
	if found, value := evaluationCache.Get(eng.game.position.hash); eng.TranspositionTableEnabled && found && value.Depth() >= depth {
		if !value.LowerBound() {
			return value.Evaluation(), mainLine
		}

		bestScore := value.Evaluation()
		if bestScore > alpha {
			alpha = bestScore

			if alpha >= beta {
				return alpha, mainLine
			}
		}
	}

	var legalMoves []*Move
	if eng.MoveSortingEnabled {
		legalMoves = eng.sortMoves(eng.game.LegalMoves(), evaluationCache)
	} else {
		legalMoves = eng.game.LegalMoves()
	}

	for i := 0; i < len(legalMoves); i++ {
		eng.game.Move(legalMoves[i])

		score, variation := eng.recNegaMax(depth-1, -beta, -alpha, evaluationCache, quiescentCache)
		score = -score
		if score > bestScore {
			bestScore = score
			mainLine = append(variation, legalMoves[i])

			if bestScore > alpha {
				alpha = bestScore

				if alpha >= beta && eng.AlphaBetaPruningEnabled {
					eng.game.UndoMove()

					hash := eng.game.position.hash.HashValue().SetData(int16(alpha), int8(depth), true)
					evaluationCache.Set(eng.game.position.hash.Key(), hash)
					return alpha, mainLine
				}
			}
		}

		eng.game.UndoMove()
	}

	hash := eng.game.position.hash.HashValue().SetData(int16(bestScore), int8(depth), false)
	evaluationCache.Set(eng.game.position.hash.Key(), hash)
	return bestScore, mainLine
}

func (eng *BruteForceEngine) quiescentSearch(depth int, alpha int, beta int, evaluationCache *ZobristTable, quiescentCache *ZobristTable) (int, []*Move) {
	nodes++

	switch eng.game.Result() {
	case Draw:
		return 0, []*Move{}
	case Checkmate:
		return -1_000_000, []*Move{}
	}

	if depth == 0 {
		return eng.StaticEvaluation(), []*Move{}
	}

	var legalMoves []*Move

	if eng.MoveSortingEnabled {
		legalMoves = eng.sortMoves(eng.game.LegalMoves(), evaluationCache)
	} else {
		legalMoves = eng.game.LegalMoves()
	}

	var bestScore int
	mainLine := []*Move{}
	// Lookup in both caches
	if found, value := quiescentCache.Get(eng.game.position.hash); eng.TranspositionTableEnabled && found {
		if !value.LowerBound() {
			return value.Evaluation(), mainLine
		}

		bestScore := value.Evaluation()
		if bestScore > alpha {
			alpha = bestScore

			if alpha >= beta {
				return alpha, mainLine
			}
		}
	} else if found, value := evaluationCache.Get(eng.game.position.hash); eng.TranspositionTableEnabled && found {
		if !value.LowerBound() {
			return value.Evaluation(), mainLine
		}

		bestScore := value.Evaluation()
		if bestScore > alpha {
			alpha = bestScore

			if alpha >= beta {
				return alpha, mainLine
			}
		}
	} else {
		// Replace with null move evaluation
		bestScore = eng.StaticEvaluation()
	}

	captureMoveFound := false

	for i := 0; i < len(legalMoves); i++ {
		if legalMoves[i].IsCapture() {
			captureMoveFound = true
			eng.game.Move(legalMoves[i])

			score, variation := eng.quiescentSearch(depth-1, -beta, -alpha, evaluationCache, quiescentCache)
			score = -score
			if score > bestScore {
				bestScore = score
				mainLine = append(variation, legalMoves[i])

				if bestScore > alpha {
					alpha = bestScore

					if alpha >= beta && eng.AlphaBetaPruningEnabled {
						eng.game.UndoMove()

						hash := eng.game.position.hash.HashValue().SetData(int16(alpha), int8(depth), true)
						quiescentCache.Set(eng.game.position.hash.Key(), hash)
						return alpha, mainLine
					}
				}
			}

			eng.game.UndoMove()
		}
	}

	if !captureMoveFound {
		return eng.StaticEvaluation(), []*Move{}
	}

	hash := eng.game.position.hash.HashValue().SetData(int16(bestScore), int8(depth), true)
	quiescentCache.Set(eng.game.position.hash.Key(), hash)
	return bestScore, mainLine
}

// StaticEvaluation returns an evaluation of the current position from a
// strategic standpoint (e.g. material imbalances, pawn structures, ...) without
// considering any tactical advantages (e.g. ability to capture a piece)
func (eng *BruteForceEngine) StaticEvaluation() int {
	score := 0
	if eng.MaterialDifferenceEval {
		score += materialDifference(eng.game.position)
	}
	if eng.PositionDifferenceEval {
		score += positionDifference(eng.game.position)
	}
	if eng.CenterControlEval {
		score += centerControl(eng.game.position)
	}
	if eng.DoubledIsolatedPawnsEval {
		score += doubledOrIsolatedPawnsPenalties(eng.game.position, &eng.game.precomputedData)
	}
	if eng.PassedPawnsEval {
		score += passedPawnsBonuses(eng.game.position, &eng.game.precomputedData)
	}

	return score * int(eng.game.position.turn)
}
