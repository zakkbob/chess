//go:build ignore
// +build ignore

package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
)

var r = rand.New(rand.NewSource(24157)) // NOTE: This number may not be optimal for preventing zobrist hash collisions

func generateArray(f *os.File, name string, length int) {
	fmt.Fprintf(f, "\t%s = [%d]uint64{", name, length)
	for i := range length {
		fmt.Fprintf(f, "0x%s", strconv.FormatUint(r.Uint64(), 16))
		if i != length-1 {
			fmt.Fprint(f, ", ")
		}
	}
	fmt.Fprint(f, "}\n")
}

func main() {
	f, err := os.Create("zobrist_values.go")
	if err != nil {
		panic(err)
	}

	f.WriteString("package chess\n\nfunc init() {\n")

	fmt.Fprintf(f, "\tblackToMoveZobrist = 0x%s\n", strconv.FormatUint(r.Uint64(), 16))
	generateArray(f, "enPassantZobrist", 8)
	generateArray(f, "castlingRightsZobrist", 16)
	fmt.Fprintln(f)

	generateArray(f, "whitePawnZobrist", 64)
	generateArray(f, "whiteRookZobrist", 64)
	generateArray(f, "whiteBishopZobrist", 64)
	generateArray(f, "whiteKnightZobrist", 64)
	generateArray(f, "whiteQueenZobrist", 64)
	generateArray(f, "whiteKingZobrist", 64)
	fmt.Fprintln(f)

	generateArray(f, "blackPawnZobrist", 64)
	generateArray(f, "blackRookZobrist", 64)
	generateArray(f, "blackBishopZobrist", 64)
	generateArray(f, "blackKnightZobrist", 64)
	generateArray(f, "blackQueenZobrist", 64)
	generateArray(f, "blackKingZobrist", 64)

	f.WriteString("}\n")
	f.Close()
}
