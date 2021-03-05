package chessboard

import "fmt"

// BruteForceEngine explores all the tree to find the best move
type BruteForceEngine struct {
	trackedGame               *Game
	game                      Game
	MaterialDifferenceEval    bool
	PositionDifferenceEval    bool
	CenterControlEval         bool
	QuiescentSearchEnabled    bool
	AlphaBetaPruningEnabled   bool
	TranspositionTableEnabled bool
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
	}
}

var nodes = 0

// BestMove returns the best move as computed by the AI
func (eng *BruteForceEngine) BestMove(remainingTime int) *Move {
	// end := time.Now().Add(time.Duration(remainingTime / 40))
	eng.game = *eng.trackedGame

	nodes = 0

	move := eng.NegaMax(3)
	fmt.Printf("Nodes explored: %d\n", nodes)
	// for depth := 2; time.Now() < remainingTime/100; depth++ {
	// 	fmt.Printf("Reached depth %d in %f seconds\n", depth, start.Sub(time.Now()).Seconds())
	// 	move = eng.NegaMax(depth)
	// }

	return move
}

// NegaMax does a negamax search of the tree up to the passed depth
func (eng *BruteForceEngine) NegaMax(depth int) *Move {
	legalMoves := eng.game.LegalMoves()

	var bestMove *Move
	bestScore := -1_000_000
	mainLine := []*Move{}
	evaluationCache := ZobristTable{}

	for i := 0; i < len(legalMoves); i++ {
		eng.game.Move(legalMoves[i])

		score, variation := eng.recNegaMax(depth-1, -1_000_000, -bestScore, &evaluationCache)
		score = -score
		if score > bestScore {
			bestScore = score
			bestMove = legalMoves[i]
			mainLine = append(variation, legalMoves[i])
		}

		eng.game.UndoMove()
	}

	fmt.Printf("Score: %d\n", bestScore)
	fmt.Print("Main line: ")
	for i := len(mainLine) - 1; i >= 0; i-- {
		fmt.Printf("%s ", *mainLine[i])
	}

	return bestMove
}

func (eng *BruteForceEngine) recNegaMax(depth int, alpha int, beta int, evaluationCache *ZobristTable) (int, []*Move) {
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
			return eng.quiescentSearch(4, alpha, beta, evaluationCache)
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

	legalMoves := eng.game.LegalMoves()

	for i := 0; i < len(legalMoves); i++ {
		eng.game.Move(legalMoves[i])

		score, variation := eng.recNegaMax(depth-1, -beta, -alpha, evaluationCache)
		score = -score
		if score > bestScore {
			bestScore = score
			mainLine = append(variation, legalMoves[i])

			if bestScore > alpha {
				alpha = bestScore

				if alpha >= beta && eng.AlphaBetaPruningEnabled {
					eng.game.UndoMove()

					hash := eng.game.position.hash.SetData(int16(alpha), int8(depth), true)
					evaluationCache.Set(hash)
					return alpha, []*Move{}
				}
			}
		}

		eng.game.UndoMove()
	}

	hash := eng.game.position.hash.SetData(int16(bestScore), int8(depth), false)
	evaluationCache.Set(hash)
	return bestScore, mainLine
}

func (eng *BruteForceEngine) quiescentSearch(depth int, alpha int, beta int, evaluationCache *ZobristTable) (int, []*Move) {
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

	legalMoves := eng.game.LegalMoves()

	var bestScore int
	mainLine := []*Move{}
	if found, value := evaluationCache.Get(eng.game.position.hash); eng.TranspositionTableEnabled && found {
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

			score, variation := eng.quiescentSearch(depth-1, -beta, -alpha, evaluationCache)
			score = -score
			if score > bestScore {
				bestScore = score
				mainLine = append(variation, legalMoves[i])

				if bestScore > alpha {
					alpha = bestScore

					if alpha >= beta && eng.AlphaBetaPruningEnabled {
						eng.game.UndoMove()
						return alpha, []*Move{}
					}
				}
			}

			eng.game.UndoMove()
		}
	}

	if !captureMoveFound {
		return eng.StaticEvaluation(), []*Move{}
	}

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

	return score * int(eng.game.position.turn)
}
