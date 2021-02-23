package main

func main() {
	c2 := Square{File(2), Rank(1)}
	c3 := Square{File(2), Rank(2)}
	e3 := Square{File(4), Rank(2)}

	b := Board{
		bbWhiteKing: SquareBitboard(c2),
		bbBlackPawn: SquareBitboard(c3),
		bbBlackKing: SquareBitboard(e3),
	}
	move := Move{from: c2, to: c3}
	b.Move(move)
}
