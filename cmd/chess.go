package main

import (
	"fmt"

	"github.com/zakkbob/chess"
)

func main() {
	b := chess.NewBoard()

	for {
		fmt.Println(b.String())

		var from int
		var to int

		fmt.Print("From: ")
		fmt.Scanln(&from)

		fmt.Print("To: ")
		fmt.Scanln(&to)

		b.DoCoordinateMove(from, to, chess.NoPromotion)
	}
}
