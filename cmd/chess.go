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
		fmt.Println("Value:", chess.Evaluate(b))

		ms := b.LegalMoves()
		/*
			for _, m := range ms {
				fmt.Println(m.String())
			}
		*/

		var move string
		var from, to int
		var err error

		legalMove := false

		for !legalMove {
			fmt.Print("Move: ")
			fmt.Scanln(&move)

			if len(move) != 4 {

				fmt.Println("Ermmm, that doesn't look like a move to me")
				continue
			}

			from, err = chess.IndexFromAlgebraic(move[0:2])
			if err != nil {
				fmt.Println("Ermmm, that doesn't look like a move to me")
				continue
			}
			to, err = chess.IndexFromAlgebraic(move[2:4])
			if err != nil {
				fmt.Println("Ermmm, that doesn't look like a move to me")
				continue
			}

			for _, m := range ms {
				if int(m.To()) == to && int(m.From()) == from {
					legalMove = true
					break
				}
			}
			if legalMove {
				break
			} else {
				fmt.Println("Aha! Caught you cheating!!")
			}
		}

		b.DoCoordinateMove(from, to, chess.NoPromotion)

		fmt.Println(b.String())

		m := chess.Search(b, 4)
		b.Move(m)
		fmt.Println("Computer did", m.String())
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
