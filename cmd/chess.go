package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/zakkbob/chess"
)

func perft(b *chess.Board, depth int) int {
	if depth == 0 {
		return 1
	}

	counter := 0

	ms := b.LegalMoves()
	for _, m := range ms {
		b.Move(m)
		nodes := perft(b, depth-1)
		counter += nodes
		b.Unmove()
	}
	return counter
}

func main() {
	b := chess.BoardFromRanks(
		[8]string{
			"        ",
			"  p     ",
			"   p    ",
			"KP     r",
			" R   p k",
			"        ",
			"    P P ",
			"        ",
		},
		chess.WhiteTurn,
		chess.NoCastleRights,
	)

	//fmt.Println(b.String())

	depth, err := strconv.Atoi(os.Args[1])
	if err != nil {
		os.Exit(depth)
	}
	//fen := os.Args[2]
	if len(os.Args) == 4 {
		moves := os.Args[3]

		for a := range strings.SplitSeq(moves, " ") {
			b.DoAlgebraicMove(a)
		}
	}

	counter := 0
	ms := b.LegalMoves()
	for _, m := range ms {
		b.Move(m)
		nodes := perft(&b, depth-1)
		b.Unmove()
		fmt.Println(m.String(), nodes)
		counter += nodes
	}
	fmt.Println()
	fmt.Println(counter)

	os.Exit(0)

	for {
		fmt.Println(b.String())

		ms := b.LegalMoves()
		fmt.Println(len(ms), "legal moves")
		for _, m := range ms {
			fmt.Println(m.String())
		}
		fmt.Println(b.CanEnPassant, b.EnPassantFile)

		var from int
		var to int

		fmt.Print("From: ")
		fmt.Scanln(&from)

		fmt.Print("To: ")
		fmt.Scanln(&to)

		b.DoCoordinateMove(from, to, chess.NoPromotion)

	}
}
