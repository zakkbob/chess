package chess_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zakkbob/chess"
)

func perft(b *chess.Board, depth int) int {
	if depth == 1 {
		return len(b.LegalMoves())
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

func TestPerft(t *testing.T) {
	tests := []struct {
		Name   string
		Board  chess.Board
		Depths []int
	}{
		{
			Name: "Initial Position",
			Depths: []int{
				20,
				400,
				8902,
				197281,
				4865609,
				//119060324,
			},
			Board: chess.BoardFromFEN("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"),
		},
		{
			Name: "Kiwipete",
			Depths: []int{
				48,
				2039,
				97862,
				4085603,
				//193690690,
			},
			Board: chess.BoardFromFEN("r3k2r/p1ppqpb1/bn2pnp1/3PN3/1p2P3/2N2Q1p/PPPBBPPP/R3K2R w KQkq -"),
		},
		{
			Name: "Position 3",
			Depths: []int{
				14,
				191,
				2812,
				43238,
				674624,
				//11030083,
			},
			Board: chess.BoardFromFEN("8/2p5/3p4/KP5r/1R3p1k/8/4P1P1/8 w - - 0 1"),
		},
		{
			Name: "Position 4",
			Depths: []int{
				6,
				264,
				9467,
				422333,
				//15833292,
			},
			Board: chess.BoardFromFEN("r3k2r/Pppp1ppp/1b3nbN/nP6/BBP1P3/q4N2/Pp1P2PP/R2Q1RK1 w kq - 0 1"),
		},
		{
			Name: "Position 5",
			Depths: []int{
				44,
				1486,
				62379,
				2103487,
				//89941194,
			},
			Board: chess.BoardFromFEN("rnbq1k1r/pp1Pbppp/2p5/8/2B5/8/PPP1NnPP/RNBQK2R w KQ - 1 8"),
		},
		{
			Name: "Position 6",
			Depths: []int{
				46,
				2079,
				89890,
				3894594,
				//164075551,
			},
			Board: chess.BoardFromFEN("r4rk1/1pp1qppp/p1np1n2/2b1p1B1/2B1P1b1/P1NP1N2/1PP1QPPP/R4RK1 w - - 0 10"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			for i := range len(tt.Depths) {
				t.Run(fmt.Sprintf("Depth %d", i+1), func(t *testing.T) {
					got := perft(&tt.Board, i+1)

					assert.Equal(t, tt.Depths[i], got, "Incorrect perft result at depth %d", i+1)
				})
			}
		})
	}
}
