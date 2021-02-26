package chessboard

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"
)

type Result int

const (
	NoResult Result = iota
	Draw
	Checkmate
)

func (res Result) String() string {
	switch res {
	case NoResult:
		return "NoResult"
	case Draw:
		return "Draw"
	case Checkmate:
		return "Checkmate"
	default:
		panic("Unknown result")
	}
}

// PrecomputedData contains all the precalculated bitboards used in move generation
type PrecomputedData struct {
	KingMoves       [64]Bitboard
	KnightMoves     [64]Bitboard
	RookMagics      [64]uint64
	RookIndexBits   [64]int
	RookMasks       [64]Bitboard
	RookMoves       [64][4096]Bitboard
	BishopMagics    [64]uint64
	BishopIndexBits [64]int
	BishopMasks     [64]Bitboard
	BishopMoves     [64][1024]Bitboard
}

// Game contains all information about the game
type Game struct {
	precomputedData  PrecomputedData
	position         *Position
	positionsHistory []*Position
	moves            []*Move
}

// Result returns the result of the current game
func (game *Game) Result() Result {
	// Check draws by insufficient material
	if game.position.board.bbWhitePawn == 0 && game.position.board.bbBlackPawn == 0 &&
		game.position.board.bbWhiteRook == 0 && game.position.board.bbBlackRook == 0 &&
		game.position.board.bbWhiteQueen == 0 && game.position.board.bbBlackQueen == 0 {
		knightsAndBishops := game.position.board.bbWhiteBishop | game.position.board.bbBlackBishop |
			game.position.board.bbWhiteKnight | game.position.board.bbBlackKnight

		// King vs King, King+Bishop vs King, King+Knight vs King
		if knightsAndBishops.KernighanPopCount() <= 1 {
			return Draw
		}

		// King+Bishop vs King+Bishop with Bishops on the same colour
		knights := game.position.board.bbWhiteKnight | game.position.board.bbBlackKnight
		if knights.KernighanPopCount() == 0 &&
			game.position.board.bbWhiteBishop.KernighanPopCount() == 1 &&
			game.position.board.bbBlackBishop.KernighanPopCount() == 1 {
			whiteBishopSquare := square(game.position.board.bbWhiteBishop.LeastSignificantBit())
			blackBishopSquare := square(game.position.board.bbBlackBishop.LeastSignificantBit())

			if whiteBishopSquare.Color() == blackBishopSquare.Color() {
				return Draw
			}
		}
	}

	// Draw by 75 moves rule
	if game.position.halfMoveClock >= 75 {
		return Draw
	}

	legalMoves := game.LegalMoves()
	if len(legalMoves) == 0 {
		var kingSquare square

		if game.position.turn == WhiteColor {
			kingSquare = game.position.board.whiteKingSquare
		} else {
			kingSquare = game.position.board.blackKingSquare
		}

		if game.position.board.IsUnderAttack(game, kingSquare) {
			return Checkmate
		}

		// Stalemate
		return Draw
	}

	return NoResult
}

// Move applies a move in the game
func (game *Game) Move(move *Move) {
	pos := game.position.Move(move)
	game.position = &pos
	game.positionsHistory = append(game.positionsHistory, game.position)
	game.moves = append(game.moves, move)
}

// UndoMove undoes the last move
func (game *Game) UndoMove() {
	game.moves = game.moves[:len(game.moves)-1]
	game.positionsHistory = game.positionsHistory[:len(game.positionsHistory)-1]
	game.position = game.positionsHistory[len(game.positionsHistory)-1]
}

// Position returns the current position in the game
func (game *Game) Position() Position {
	return *game.position
}

// LoadPrecomputedData loads all the precomputed data for fast move generation
func (game *Game) LoadPrecomputedData(path string) {
	jsonBytes, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var data PrecomputedData
	err = json.Unmarshal(jsonBytes, &data)
	if err != nil {
		panic(err)
	}

	game.precomputedData = data
}

// NewGame initializes a new game
func NewGame() Game {
	startingPositionFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	return NewGameFromFEN(startingPositionFEN)
}

// NewGameFromFEN initializes a game from a fen string
func NewGameFromFEN(fen string) Game {
	game := Game{}
	game.LoadPrecomputedData("/Users/andreazanin/Code/go/chess_engine/precomputed.json")
	game.positionsHistory = make([]*Position, 0, 40)
	game.moves = make([]*Move, 0, 40)
	pos := Position{}
	game.position = &pos

	fen = strings.TrimSpace(fen)
	pieces := strings.Split(fen, " ")
	if len(pieces) != 6 {
		panic("Invalid fen passed: it should have 6 pieces")
	}

	game.position.board = parseFenBoard(pieces[0])

	switch pieces[1] {
	case "w":
		game.position.turn = WhiteColor
	case "b":
		game.position.turn = BlackColor
	default:
		panic("Invalid fen turn string")
	}

	game.position.castleRights = parseCastleRights(pieces[2])
	enPassantSquare, ok := stringToSquare[pieces[3]]
	if !ok {
		panic("Unrecognized en passant square")
	}
	game.position.enPassantSquare = enPassantSquare

	halfMoveClock, err := strconv.Atoi(pieces[4])
	if err != nil || halfMoveClock < 0 {
		panic("Half move clock should be a non negative number")
	}
	game.position.halfMoveClock = halfMoveClock

	moveCount, err := strconv.Atoi(pieces[5])
	if err != nil || moveCount < 1 {
		panic("Move count should be a positive number")
	}
	game.position.moveCount = moveCount

	game.positionsHistory = []*Position{game.position}
	// TODO: update in check status

	return game
}

func parseCastleRights(rawRights string) CastleRights {
	rights := CastleRights{}
	if strings.Contains(rawRights, "K") {
		rights.WhiteKingSide = true
	}
	if strings.Contains(rawRights, "Q") {
		rights.WhiteQueenSide = true
	}
	if strings.Contains(rawRights, "k") {
		rights.BlackKingSide = true
	}
	if strings.Contains(rawRights, "q") {
		rights.BlackQueenSide = true
	}

	return rights
}

// example string: rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR
func parseFenBoard(rawBoard string) Board {
	board := Board{}

	currentSquare := A8
	index := 0

	// parse fen character by character
	for index < len(rawBoard) {
		char := rawBoard[index]

		switch char {
		case 'K':
			board.bbWhiteKing |= currentSquare.Bitboard()
			currentSquare++
		case 'Q':
			board.bbWhiteQueen |= currentSquare.Bitboard()
			currentSquare++
		case 'R':
			board.bbWhiteRook |= currentSquare.Bitboard()
			currentSquare++
		case 'B':
			board.bbWhiteBishop |= currentSquare.Bitboard()
			currentSquare++
		case 'N':
			board.bbWhiteKnight |= currentSquare.Bitboard()
			currentSquare++
		case 'P':
			board.bbWhitePawn |= currentSquare.Bitboard()
			currentSquare++
		case 'k':
			board.bbBlackKing |= currentSquare.Bitboard()
			currentSquare++
		case 'q':
			board.bbBlackQueen |= currentSquare.Bitboard()
			currentSquare++
		case 'r':
			board.bbBlackRook |= currentSquare.Bitboard()
			currentSquare++
		case 'b':
			board.bbBlackBishop |= currentSquare.Bitboard()
			currentSquare++
		case 'n':
			board.bbBlackKnight |= currentSquare.Bitboard()
			currentSquare++
		case 'p':
			board.bbBlackPawn |= currentSquare.Bitboard()
			currentSquare++
		case '/':
			currentSquare -= 16
		default:
			jump, err := strconv.Atoi(string(char))
			if err != nil {
				panic("Unknown character in FEN board")
			}

			currentSquare += square(jump)
		}

		index++
	}

	board.FillSupportBitboards()

	return board
}
