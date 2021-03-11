package chessboard

import (
	"fmt"
	"time"
)

// Infinity contains a very high int16 number
const Infinity = 15_000_000

// CheckmateScore contains the score given to a checkmate loss
const CheckmateScore = -1_000_000

// DrawScore contains the score given to a draw
const DrawScore = 0

// BruteForceEngine explores all the tree to find the best move
type BruteForceEngine struct {
	trackedGame               *Game
	game                      Game
	MaxDepth                  int
	MaterialDifferenceEval    bool
	PositionDifferenceEval    bool
	CenterControlEval         bool
	DoubledIsolatedPawnsEval  bool
	PassedPawnsEval           bool
	QuiescentSearchEnabled    bool
	AlphaBetaPruningEnabled   bool
	TranspositionTableEnabled bool
	MoveSortingEnabled        bool
	AspirationSearchEnabled   bool
	AspirationWindowWidth     int
}

// NewBruteForceEngine initializes a BruteForceEngine
func NewBruteForceEngine(game *Game) *BruteForceEngine {
	return &BruteForceEngine{trackedGame: game,
		MaterialDifferenceEval:    true,
		PositionDifferenceEval:    true,
		QuiescentSearchEnabled:    true,
		AlphaBetaPruningEnabled:   true,
		CenterControlEval:         true,
		TranspositionTableEnabled: false,
		MoveSortingEnabled:        true,
		DoubledIsolatedPawnsEval:  true,
		PassedPawnsEval:           true,
		MaxDepth:                  -1,
		AspirationSearchEnabled:   true,
		AspirationWindowWidth:     180,
	}
}

var _nodeHits = 0

// BestMove returns the best move as computed by the AI
func (eng *BruteForceEngine) BestMove(remainingTime int) *Move {
	var endTime time.Time
	if eng.MaxDepth == -1 {
		endTime = time.Now().Add(time.Duration(remainingTime) * (time.Second / 40))
	} else {
		endTime = time.Now().Add(720 * time.Hour)
	}

	eng.game = *eng.trackedGame

	_nodeHits = 0
	_zobristCacheHits = 0
	_zobristCacheMisses = 0

	move := eng.game.LegalMoves()[0]
	evaluations := &(ZobristTable{})
	quiescentEvaluations := &(ZobristTable{})
	depth := 1
	previousScore := eng.StaticEvaluation()
	for eng.MaxDepth == -1 || depth <= eng.MaxDepth {
		var aborted bool
		var bestMove *Move
		var score int

		if eng.AspirationSearchEnabled {
			aborted, bestMove, score = eng.NegaMax(
				depth,
				previousScore-eng.AspirationWindowWidth,
				previousScore+eng.AspirationWindowWidth,
				endTime,
				evaluations,
				quiescentEvaluations,
			)

			if aborted {
				break
			}
		}

		// If the aspiration search is disabled than we perform a full search directly
		// If the evaluation from the aspiration search is outside the bound of the aspiration window we must research
		if !eng.AspirationSearchEnabled || score >= previousScore+eng.AspirationWindowWidth || score <= previousScore-eng.AspirationWindowWidth {
			aborted, bestMove, score = eng.NegaMax(
				depth,
				-Infinity,
				Infinity,
				endTime,
				evaluations,
				quiescentEvaluations,
			)
		}

		if aborted {
			break
		}

		move = bestMove
		previousScore = score
		depth++
	}

	fmt.Printf("Depth reached: %d, Nodes explored: %d, Cache hits: %d, Cache misses: %d\n", depth-1, _nodeHits, _zobristCacheHits, _zobristCacheMisses)

	return move
}

// NegaMax does a negamax search of the tree up to the passed depth
func (eng *BruteForceEngine) NegaMax(depth int, alpha int, beta int, endTime time.Time, evaluations *ZobristTable, quiescentEvaluations *ZobristTable) (bool, *Move, int) {
	var legalMoves []*Move

	// We can use the evaluation score from the previous iteration to sort the moves,
	// exploring first the moves that have the highest chance of being the best one
	// allows the alpha-beta algorithm to prune more branches
	if eng.MoveSortingEnabled {
		legalMoves = eng.sortMoves(eng.game.LegalMoves(), evaluations)
	} else {
		legalMoves = eng.game.LegalMoves()
	}

	bestMove := legalMoves[0]
	bestScore := -Infinity
	bestPositionalScore := -Infinity

	// mainLine is a diagnostic value that tracks what the bot thinks
	// would be the perfect game from now on
	mainLine := []*Move{}

	// Try each move, recursively compute the score of the resulting position and
	// choose the best move for us (that is the worst for our opponent)
	for i := 0; i < len(legalMoves); i++ {
		// Abort search if running out of time
		if time.Now().After(endTime) {
			return true, nil, 0
		}

		eng.game.Move(legalMoves[i])

		// Get the evaluation of the position from our opponents point of view and flip it (best for us is worst for our opponent)
		score, variation := eng.recNegaMax(depth-1, -beta, -alpha, evaluations, quiescentEvaluations)
		score = -score

		if score > bestScore {
			bestScore = score
			bestPositionalScore = -eng.StaticEvaluation()
			bestMove = legalMoves[i]
			mainLine = append(variation, legalMoves[i])

			if bestScore > alpha {
				alpha = bestScore

				// The score will be outside the bounds of the aspiration window, so we can
				// stop the search already
				if alpha > beta && eng.AlphaBetaPruningEnabled {
					eng.game.UndoMove()
					return false, bestMove, alpha
				}
			}
		} else if score == bestScore {
			// When the score is the same choose the move with the better static evaluation
			positionalScore := -eng.StaticEvaluation()
			if positionalScore > bestPositionalScore {
				bestScore = score
				bestPositionalScore = positionalScore
				bestMove = legalMoves[i]
				mainLine = append(variation, legalMoves[i])
			}
		}

		eng.game.UndoMove()
	}

	// Log diagnostics about the best move found with this depth of search
	fmt.Printf("Depth: %d, Score: %d, Main line: ", depth, bestScore)
	for i := len(mainLine) - 1; i >= 0; i-- {
		fmt.Printf("%s ", *mainLine[i])
	}
	fmt.Println()

	return false, bestMove, bestScore
}

func (eng *BruteForceEngine) recNegaMax(depth int, alpha int, beta int, evaluationCache *ZobristTable, quiescentCache *ZobristTable) (int, []*Move) {
	_nodeHits++

	switch eng.game.Result() {
	case Draw:
		return DrawScore, []*Move{}
	case Checkmate:
		return CheckmateScore, []*Move{}
	}

	if depth == 0 {
		// When reaching depth 0 we can procede the search deeper but considering only capture
		// moves, this way mitigate the horizon effect and correctly assess trades
		if eng.QuiescentSearchEnabled {
			_nodeHits-- // avoid double counting this node
			return eng.quiescentSearch(7, alpha, beta, evaluationCache, quiescentCache)
		}

		return eng.StaticEvaluation(), []*Move{}
	}

	bestScore := -Infinity
	mainLine := []*Move{}

	// If this position's evaluation is cached we don't need to recompute it, we could have also stored a lower bound
	// because the search had been stopped by alpha-beta pruning.
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
	// Sorting the moves using the past evaluations allows us to evaluate first the moves that have a high
	// probability of being the best, this makes alpha-beta pruning more effective
	if eng.MoveSortingEnabled {
		legalMoves = eng.sortMoves(eng.game.LegalMoves(), evaluationCache)
	} else {
		legalMoves = eng.game.LegalMoves()
	}

	for i := 0; i < len(legalMoves); i++ {
		eng.game.Move(legalMoves[i])

		score, variation := eng.recNegaMax(depth-1, -beta, -alpha, evaluationCache, quiescentCache)
		score = -score
		eng.game.UndoMove()

		if score > bestScore {
			bestScore = score
			mainLine = append(variation, legalMoves[i])

			if bestScore > alpha {
				alpha = bestScore

				// If we find a position that is too good we can stop the search, because the player making the
				// previous move will opt for a move giving us a weaker position.
				if alpha > beta && eng.AlphaBetaPruningEnabled {
					// Store evaluation in the cache
					hash := eng.game.position.hash.HashValue().SetData(int16(alpha), int8(depth), true)
					evaluationCache.Set(eng.game.position.hash.Key(), hash)

					return alpha, mainLine
				}
			}
		}
	}

	// Store evaluation in the cache
	hash := eng.game.position.hash.HashValue().SetData(int16(bestScore), int8(depth), false)
	evaluationCache.Set(eng.game.position.hash.Key(), hash)

	return bestScore, mainLine
}

func (eng *BruteForceEngine) quiescentSearch(depth int, alpha int, beta int, evaluationCache *ZobristTable, quiescentCache *ZobristTable) (int, []*Move) {
	_nodeHits++

	switch eng.game.Result() {
	case Draw:
		return DrawScore, []*Move{}
	case Checkmate:
		return CheckmateScore, []*Move{}
	}

	// At depth 0 we statically evaluate the position with the implemented heuristics
	if depth == 0 {
		return eng.StaticEvaluation(), []*Move{}
	}

	var legalMoves []*Move
	if eng.MoveSortingEnabled {
		legalMoves = eng.sortMoves(eng.game.LegalMoves(), evaluationCache)
	} else {
		legalMoves = eng.game.LegalMoves()
	}

	var bestScore int = -Infinity
	mainLine := []*Move{}
	// Inside the quiescent search we can read evaluations from both the quiescent evaluations cache
	// and the full evaluation cache; the latter don't need to be depth checked, because they surely
	// searched deeper than the quiescent search would do.
	if found, value := quiescentCache.Get(eng.game.position.hash); eng.TranspositionTableEnabled && found && value.Depth() >= depth {
		if !value.LowerBound() {
			return value.Evaluation(), mainLine
		}

		bestScore = value.Evaluation()
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

		bestScore = value.Evaluation()
		if bestScore > alpha {
			alpha = bestScore

			if alpha >= beta {
				return alpha, mainLine
			}
		}
	} else {
		// Replace with null move evaluation
		bestScore = eng.StaticEvaluation()

		// Null move heuristic
		// eng.game.Move(&NullMove)
		// score, _ := eng.quiescentSearch(depth-1, -beta, -alpha, evaluationCache, quiescentCache)
		// bestScore = score
		// eng.game.UndoMove()

		// if bestScore > alpha {
		// 	alpha = bestScore

		// 	if alpha >= beta && eng.AlphaBetaPruningEnabled {
		// 		hash := eng.game.position.hash.HashValue().SetData(int16(alpha), int8(depth), true)
		// 		quiescentCache.Set(eng.game.position.hash.Key(), hash)
		// 		return alpha, mainLine
		// 	}
		// }
	}

	// If there are no disruptive moves, then we have found a quiescent position,
	// we can stop the search and evaluate statically this position
	disruptiveMoveFound := false

	for i := 0; i < len(legalMoves); i++ {
		if eng.game.position.inCheck || legalMoves[i].IsCapture() {
			disruptiveMoveFound = true

			eng.game.Move(legalMoves[i])
			score, variation := eng.quiescentSearch(depth-1, -beta, -alpha, evaluationCache, quiescentCache)
			score = -score
			eng.game.UndoMove()

			if score > bestScore {
				bestScore = score
				mainLine = append(variation, legalMoves[i])

				if bestScore > alpha {
					alpha = bestScore

					if alpha > beta && eng.AlphaBetaPruningEnabled {
						// Store the evaluation in the cache
						hash := eng.game.position.hash.HashValue().SetData(int16(alpha), int8(depth), true)
						quiescentCache.Set(eng.game.position.hash.Key(), hash)
						return alpha, mainLine
					}
				}
			}
		}
	}

	// When finding a quiescent position return the static evaluation
	if !disruptiveMoveFound {
		return eng.StaticEvaluation(), []*Move{}
	}

	// Store the evaluation in the cache
	hash := eng.game.position.hash.HashValue().SetData(int16(bestScore), int8(depth), true)
	quiescentCache.Set(eng.game.position.hash.Key(), hash)

	return bestScore, mainLine
}
