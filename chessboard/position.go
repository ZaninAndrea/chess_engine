package chessboard

// CastleRights stored information about which side each player
// is still allowed to castle to
type CastleRights struct {
	WhiteKingSide  bool
	WhiteQueenSide bool
	BlackKingSide  bool
	BlackQueenSide bool
}

// Position contains all the information about a give position in a game
// including turn, enpassant information, castle rights, valid moves in this position...
type Position struct {
	board           *Board
	turn            Color
	castleRights    CastleRights
	enPassantSquare Square
	halfMoveClock   int
	moveCount       int
	inCheck         bool
	legalMoves      []*Move
}