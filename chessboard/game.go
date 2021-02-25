package chessboard

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

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

type Game struct {
	PrecomputedData  PrecomputedData
	Position         *Position
	PositionsHistory []*Position
	Moves            []Move
}

func (game Game) String() string {
	return fmt.Sprint(game.Moves)
}

// Move applies a move in the game
func (game *Game) Move(move Move) {
	pos := game.Position.Move(move)
	game.Position = &pos
	game.PositionsHistory = append(game.PositionsHistory, game.Position)
	game.Moves = append(game.Moves, move)
}

// UndoMove undoes the last move
func (game *Game) UndoMove() {
	game.Moves = game.Moves[:len(game.Moves)-1]
	game.PositionsHistory = game.PositionsHistory[:len(game.PositionsHistory)-1]
	game.Position = game.PositionsHistory[len(game.PositionsHistory)-1]
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

	game.PrecomputedData = data
}

// NewGame initializes a new game
func NewGame() Game {
	startingPositionFEN := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	return NewGameFromFEN(startingPositionFEN)
}

// NewGameFromFEN initializes a game from a fen string
func NewGameFromFEN(fen string) Game {
	game := Game{}
	game.LoadPrecomputedData("./precomputed.json")
	game.PositionsHistory = []*Position{}
	game.Moves = []Move{}
	pos := Position{}
	game.Position = &pos

	fen = strings.TrimSpace(fen)
	pieces := strings.Split(fen, " ")
	if len(pieces) != 6 {
		panic("Invalid fen passed: it should have 6 pieces")
	}

	game.Position.board = parseFenBoard(pieces[0])

	switch pieces[1] {
	case "w":
		game.Position.turn = WhiteColor
	case "b":
		game.Position.turn = BlackColor
	default:
		panic("Invalid fen turn string")
	}

	game.Position.castleRights = parseCastleRights(pieces[2])
	enPassantSquare, ok := stringToSquare[pieces[3]]
	if !ok {
		panic("Unrecognized en passant square")
	}
	game.Position.enPassantSquare = enPassantSquare

	halfMoveClock, err := strconv.Atoi(pieces[4])
	if err != nil || halfMoveClock < 0 {
		panic("Half move clock should be a non negative number")
	}
	game.Position.halfMoveClock = halfMoveClock

	moveCount, err := strconv.Atoi(pieces[5])
	if err != nil || moveCount < 1 {
		panic("Move count should be a positive number")
	}
	game.Position.moveCount = moveCount

	game.PositionsHistory = []*Position{game.Position}
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
			board.bbBlackKnight |= currentSquare.Bitboard()
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

	return board
}
