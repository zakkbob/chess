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

func perftCommand(args []string) {
	depth, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("Cannot parse depth: ", args[0])
		os.Exit(1)
	}

	fen := args[1]
	b, err := chess.BoardFromFEN(fen)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if len(args) == 3 {
		moves := args[2]

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
}

func playCommand(args []string) {
	b := chess.NewBoard()
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

func main() {
	switch os.Args[1] {
	case "perft":
		perftCommand(os.Args[2:])
	case "play":
		playCommand(os.Args[2:])
	default:
		fmt.Println("expected 'perft' or 'play' subcommands")
		os.Exit(1)
	}

}
