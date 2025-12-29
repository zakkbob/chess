package main

import (
	"fmt"

	"github.com/zakkbob/chess"
)

func playCommand(args []string) {
	var (
		whiteIsEngine bool
		blackIsEngine bool
	)

	e := chess.Engine{
		B:  chess.NewBoard(),
		TT: *chess.NewTranspositionTable(10),
		EP: chess.DefaultParams,
	}

	fmt.Print("White is engine? ")
	fmt.Scanln(&whiteIsEngine)
	fmt.Print("Black is engine? ")
	fmt.Scanln(&blackIsEngine)

	for {
		fmt.Println(e.B.String())
		fmt.Println("Value:", e.Evaluate())

		ms, _ := e.B.LegalMoves()
		displayLegalMoves(ms)

		if (e.B.Turn == chess.WhiteTurn && whiteIsEngine) || (e.B.Turn == chess.BlackTurn && blackIsEngine) {
			m := e.Search(1)
			e.B.Move(m)
			fmt.Println("Engine did", m.String())
		} else {
			doHumanMove(&e.B, ms)
		}

		if len(ms) == 0 {
			fmt.Println("Game is over, I wonder who won")
		}

		fmt.Println()
	}
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
