package chessboard

import "fmt"

// Move contains the informations about a move
type Move struct {
	from      Square
	to        Square
	promotion Piece
}

func (m Move) String() string {
	if m.promotion != NoPiece {
		return fmt.Sprintf("%s-%s=%s", m.from, m.to, m.promotion)
	}

	return fmt.Sprintf("%s-%s", m.from, m.to)
}
