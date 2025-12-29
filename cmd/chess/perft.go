package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/zakkbob/chess"
)

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
	ms, _ := b.LegalMoves()
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

func perft(b *chess.Board, depth int) int {
	if depth == 0 {
		return 1
	}

	counter := 0

	ms, _ := b.LegalMoves()
	for _, m := range ms {
		b.Move(m)
		nodes := perft(b, depth-1)
		counter += nodes
		b.Unmove()
	}
	return counter
}
