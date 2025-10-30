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

	ms, _ := b.LegalMoves()
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

func isMoveLegal(from, to int, p chess.Promotion, ms []chess.Move) bool {
	for _, m := range ms {
		if int(m.To()) == to && int(m.From()) == from && m.Promotion() == p {
			return true
		}
	}
	return false
}

func displayLegalMoves(ms []chess.Move) {
	fmt.Print("Legal moves: ")
	for _, m := range ms {
		fmt.Print(m.String() + " ")
	}
	fmt.Println("(" + fmt.Sprint(len(ms)) + ")")
}

func doHumanMove(b *chess.Board, ms []chess.Move) {
	var (
		move     string
		from, to int
		p        chess.Promotion
		err      error
	)

	legalMove := false

	for !legalMove {
		fmt.Print("Move: ")
		fmt.Scanln(&move)

		from, to, p, err = chess.ParseAlgebraicMove(move)
		if err != nil {
			fmt.Println("Errrm, that doesn't look like a valid move to me")
			continue
		}

		legalMove = isMoveLegal(from, to, p, ms)

		if !legalMove {
			fmt.Println("Aha! Caught you cheating!!")
		}
	}

	b.DoCoordinateMove(from, to, p)
}

func playCommand(args []string) {
	b := chess.NewBoard()

	var (
		whiteIsEngine bool
		blackIsEngine bool
	)

	fmt.Print("White is engine? ")
	fmt.Scanln(&whiteIsEngine)
	fmt.Print("Black is engine? ")
	fmt.Scanln(&blackIsEngine)

	for {
		fmt.Println(b.String())
		fmt.Println("Value:", chess.Evaluate(b))

		ms, _ := b.LegalMoves()
		displayLegalMoves(ms)

		if (b.Turn == chess.WhiteTurn && whiteIsEngine) || (b.Turn == chess.BlackTurn && blackIsEngine) {
			m := chess.Search(b, 1)
			b.Move(m)
			fmt.Println("Engine did", m.String())
		} else {
			doHumanMove(&b, ms)
		}

		if len(ms) == 0 {
			fmt.Println("Game is over, I wonder who won")
		}

		fmt.Println()
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
