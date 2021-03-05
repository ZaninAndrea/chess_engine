package chessboard

import "math/rand"

// RandomEngine simply picks a random move as best move
type RandomEngine struct {
	game *Game
}

// NewRandomEngine initializes a RandomEngine
func NewRandomEngine(game *Game) Engine {
	return &RandomEngine{game: game}
}

// BestMove returns the best move as computed by the AI
func (eng *RandomEngine) BestMove(remainingTime int) *Move {
	legalMoves := eng.game.LegalMoves()

	return legalMoves[rand.Intn(len(legalMoves))]
}
