package main

import (
	"fmt"

	"github.com/zakkbob/chess"
)

func main() {
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
