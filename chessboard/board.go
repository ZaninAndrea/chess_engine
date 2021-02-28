package chessboard

import "fmt"

// Board contains the position of all pieces on the chessboard
type Board struct {
	bbWhiteKing     Bitboard
	bbWhiteQueen    Bitboard
	bbWhiteRook     Bitboard
	bbWhiteBishop   Bitboard
	bbWhiteKnight   Bitboard
	bbWhitePawn     Bitboard
	bbBlackKing     Bitboard
	bbBlackQueen    Bitboard
	bbBlackRook     Bitboard
	bbBlackBishop   Bitboard
	bbBlackKnight   Bitboard
	bbBlackPawn     Bitboard
	whiteSquares    Bitboard
	blackSquares    Bitboard
	emptySquares    Bitboard
	whiteKingSquare square
	blackKingSquare square
}

// FillSupportBitboards computes and sets the support bitboard in the Board
func (b *Board) FillSupportBitboards() {
	b.whiteSquares = b.bbWhiteKing | b.bbWhiteQueen | b.bbWhiteRook |
		b.bbWhiteBishop | b.bbWhiteKnight | b.bbWhitePawn
	b.blackSquares = b.bbBlackKing | b.bbBlackQueen | b.bbBlackRook |
		b.bbBlackBishop | b.bbBlackKnight | b.bbBlackPawn
	b.emptySquares = ^(b.whiteSquares | b.blackSquares)
	b.whiteKingSquare = square(b.bbWhiteKing.LeastSignificant1Bit())
	b.blackKingSquare = square(b.bbBlackKing.LeastSignificant1Bit())
}

// Piece returns the piece in a given square of the board
func (b *Board) Piece(s square) Piece {
	bbSquare := s.Bitboard()
	switch {
	case b.bbWhiteKing&bbSquare != 0:
		return WhiteKing
	case b.bbWhiteQueen&bbSquare != 0:
		return WhiteQueen
	case b.bbWhiteRook&bbSquare != 0:
		return WhiteRook
	case b.bbWhiteBishop&bbSquare != 0:
		return WhiteBishop
	case b.bbWhiteKnight&bbSquare != 0:
		return WhiteKnight
	case b.bbWhitePawn&bbSquare != 0:
		return WhitePawn
	case b.bbBlackKing&bbSquare != 0:
		return BlackKing
	case b.bbBlackQueen&bbSquare != 0:
		return BlackQueen
	case b.bbBlackRook&bbSquare != 0:
		return BlackRook
	case b.bbBlackBishop&bbSquare != 0:
		return BlackBishop
	case b.bbBlackKnight&bbSquare != 0:
		return BlackKnight
	case b.bbBlackPawn&bbSquare != 0:
		return BlackPawn
	default:
		return NoPiece
	}
}

func (b Board) String() string {
	s := ""
	for r := 7; r >= 0; r-- {
		s += fmt.Sprintf("%d ", r+1)
		for f := 0; f < 8; f++ {
			s += b.Piece(square(f+r*8)).String() + " "
		}
		s += "\n"
	}
	s += "  A B C D E F G H"

	return s
}

// Move updates the boarding moving the piece in the starting square to the target square
// it also captures the square in the target square if needed.
// Returns the update to the zobrist hash
func (b *Board) Move(move *Move) uint64 {
	var piece Piece
	var hash uint64
	if move.promotion != NoPiece {
		piece = move.promotion
		hash = zobristHashMoves[b.Piece(move.from)-1][move.from] ^ zobristHashMoves[piece-1][move.to]
	} else {
		piece = b.Piece(move.from)
		hash = zobristHashMoves[piece-1][move.from] ^ zobristHashMoves[piece-1][move.to]
	}

	targetPiece := b.Piece(move.to)
	if targetPiece != NoPiece {
		hash ^= zobristHashMoves[targetPiece-1][move.to]
	}

	fromBB := move.from.Bitboard()
	toBB := move.to.Bitboard()
	othersBB := ^(fromBB | toBB)

	// Remove the piece in the start and target squares
	b.bbWhiteKing &= othersBB
	b.bbWhiteQueen &= othersBB
	b.bbWhiteRook &= othersBB
	b.bbWhiteBishop &= othersBB
	b.bbWhiteKnight &= othersBB
	b.bbWhitePawn &= othersBB
	b.bbBlackKing &= othersBB
	b.bbBlackQueen &= othersBB
	b.bbBlackRook &= othersBB
	b.bbBlackBishop &= othersBB
	b.bbBlackKnight &= othersBB
	b.bbBlackPawn &= othersBB

	// Place piece in target square
	switch piece {
	case WhiteKing:
		b.bbWhiteKing |= toBB
	case WhiteQueen:
		b.bbWhiteQueen |= toBB
	case WhiteRook:
		b.bbWhiteRook |= toBB
	case WhiteBishop:
		b.bbWhiteBishop |= toBB
	case WhiteKnight:
		b.bbWhiteKnight |= toBB
	case WhitePawn:
		b.bbWhitePawn |= toBB
	case BlackKing:
		b.bbBlackKing |= toBB
	case BlackQueen:
		b.bbBlackQueen |= toBB
	case BlackRook:
		b.bbBlackRook |= toBB
	case BlackBishop:
		b.bbBlackBishop |= toBB
	case BlackKnight:
		b.bbBlackKnight |= toBB
	case BlackPawn:
		b.bbBlackPawn |= toBB
	default:
		panic("From position is empty")
	}

	// Update king position
	if piece == WhiteKing {
		b.whiteKingSquare = move.to
	} else if piece == BlackKing {
		b.blackKingSquare = move.to
	}

	// Update summary bitboards
	b.emptySquares = (b.emptySquares | fromBB) & (^toBB)
	if piece.Color() == WhiteColor {
		b.whiteSquares = (b.whiteSquares | toBB) ^ fromBB
		b.blackSquares = b.blackSquares & (^toBB)
	} else {
		b.blackSquares = (b.blackSquares | toBB) ^ fromBB
		b.whiteSquares = b.whiteSquares & (^toBB)
	}

	// Move rook in castling
	if move.IsCastle() {
		if move.flags&WhiteKingCastleFlag != 0 {
			b.bbWhiteRook |= F1.Bitboard()
			b.bbWhiteRook ^= H1.Bitboard()

			b.whiteSquares |= F1.Bitboard()
			b.whiteSquares ^= H1.Bitboard()

			b.emptySquares ^= F1.Bitboard()
			b.emptySquares |= H1.Bitboard()

			hash ^= zobristHashMoves[WhiteRook-1][F1] ^ zobristHashMoves[WhiteRook-1][H1]
		} else if move.flags&WhiteQueenCastleFlag != 0 {
			b.bbWhiteRook |= D1.Bitboard()
			b.bbWhiteRook ^= A1.Bitboard()

			b.whiteSquares |= D1.Bitboard()
			b.whiteSquares ^= A1.Bitboard()

			b.emptySquares ^= D1.Bitboard()
			b.emptySquares |= A1.Bitboard()

			hash ^= zobristHashMoves[WhiteRook-1][D1] ^ zobristHashMoves[WhiteRook-1][A1]
		} else if move.flags&BlackKingCastleFlag != 0 {
			b.bbBlackRook |= F8.Bitboard()
			b.bbBlackRook ^= H8.Bitboard()

			b.blackSquares |= F8.Bitboard()
			b.blackSquares ^= H8.Bitboard()

			b.emptySquares ^= F8.Bitboard()
			b.emptySquares |= H8.Bitboard()
			hash ^= zobristHashMoves[BlackRook-1][F8] ^ zobristHashMoves[BlackRook-1][H8]
		} else if move.flags&BlackQueenCastleFlag != 0 {
			b.bbBlackRook |= D8.Bitboard()
			b.bbBlackRook ^= A8.Bitboard()

			b.blackSquares |= D8.Bitboard()
			b.blackSquares ^= A8.Bitboard()

			b.emptySquares ^= D8.Bitboard()
			b.emptySquares |= A8.Bitboard()
			hash ^= zobristHashMoves[BlackRook-1][D8] ^ zobristHashMoves[BlackRook-1][A8]
		}
	}

	// capture en passant pawn
	if move.IsEnPassant() {
		if move.flags&WhiteEnPassantFlag != 0 {
			blackPawnPosition := (move.to - 8).Bitboard()

			b.bbBlackPawn ^= blackPawnPosition
			b.blackSquares ^= blackPawnPosition
			b.emptySquares |= blackPawnPosition
		} else {
			whitePawnPosition := (move.to + 8).Bitboard()

			b.bbWhitePawn ^= whitePawnPosition
			b.whiteSquares ^= whitePawnPosition
			b.emptySquares |= whitePawnPosition
		}
	}

	return hash
}

// IsUnderAttack returns whether the current board is in check
func (board *Board) IsUnderAttack(game *Game, sq square) bool {
	var enemyKnights Bitboard
	var enemyBishopLikes Bitboard
	var enemyRookLikes Bitboard
	var enemyKing Bitboard

	if game.position.turn == WhiteColor {
		enemyKnights = board.bbBlackKnight
		enemyBishopLikes = board.bbBlackBishop | board.bbBlackQueen
		enemyRookLikes = board.bbBlackRook | board.bbBlackQueen
		enemyKing = board.bbBlackKing
	} else {
		enemyKnights = board.bbWhiteKnight
		enemyBishopLikes = board.bbWhiteBishop | board.bbWhiteQueen
		enemyRookLikes = board.bbWhiteRook | board.bbWhiteQueen
		enemyKing = board.bbWhiteKing
	}

	kingCollisions := game.precomputedData.KingMoves[sq] & enemyKing
	if kingCollisions != 0 {
		return true
	}

	// Simulate putting a knight in the square where the allied king is, if the simulated
	// knight attacks an enemy knight then our king is in check by an enemy knight
	knightCollisions := game.precomputedData.KnightMoves[sq] & enemyKnights
	if knightCollisions != 0 {
		return true
	}

	// Simulate rook and queens moving horizontally/vertically
	blockers := (^board.emptySquares) & game.precomputedData.RookMasks[sq]
	key := (uint64(blockers) * game.precomputedData.RookMagics[sq]) >> (64 - game.precomputedData.RookIndexBits[sq])
	rookCollisions := game.precomputedData.RookMoves[sq][key] & enemyRookLikes
	if rookCollisions != 0 {
		return true
	}

	// Simulate bishop and queens moving diagonally
	blockers = (^board.emptySquares) & game.precomputedData.BishopMasks[sq]
	key = (uint64(blockers) * game.precomputedData.BishopMagics[sq]) >> (64 - game.precomputedData.BishopIndexBits[sq])
	bishopCollisions := game.precomputedData.BishopMoves[sq][key] & enemyBishopLikes
	if bishopCollisions != 0 {
		return true
	}

	if game.position.turn == WhiteColor {
		if sq < H7 && sq%8 != 0 {
			upLeftSquare := (sq + 7).Bitboard()

			if (board.bbBlackPawn & upLeftSquare) != 0 {
				return true
			}
		}

		if sq < H7 && sq%8 != 7 {
			upRightSquare := (sq + 9).Bitboard()

			if (board.bbBlackPawn & upRightSquare) != 0 {
				return true
			}
		}
	} else {
		if sq > H2 && sq%8 != 0 {
			downLeftSquare := (sq - 9).Bitboard()

			if (board.bbWhitePawn & downLeftSquare) != 0 {
				return true
			}
		}
		if sq > H2 && sq%8 != 7 {
			downRightSquare := (sq - 7).Bitboard()

			if (board.bbWhitePawn & downRightSquare) != 0 {
				return true
			}
		}
	}

	return false
}
