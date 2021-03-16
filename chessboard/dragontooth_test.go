package chessboard

import (
	"fmt"
	"testing"

	"github.com/dylhunn/dragontoothmg"
)

func countAllMovesDragontooth(board *dragontoothmg.Board, depth int) int {
	if depth == 0 {
		return 1
	}

	total := 0
	moveList := board.GenerateLegalMoves()

	if depth == 1 {
		return len(moveList)
	}

	for _, currMove := range moveList {
		// Apply it to the board
		unapplyFunc := board.Apply(currMove)
		total += countAllMovesDragontooth(board, depth-1)
		unapplyFunc()
	}

	return total
}

func BenchmarkMoveGenerationOther5PlyDragontooth(b *testing.B) {
	board := dragontoothmg.ParseFen("rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8")
	total := 0
	for i := 0; i < b.N; i++ {
		total = countAllMovesDragontooth(&board, 5)
	}
	fmt.Println(total)
}
