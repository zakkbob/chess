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
			Board: chess.NewBoard(),
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
			Board: chess.BoardFromRanks(
				[8]string{
					"r   k  r",
					"p ppqpb ",
					"bn  pnp ",
					"   PN   ",
					" p  P   ",
					"  N  Q p",
					"PPPBBPPP",
					"R   K  R",
				},
				chess.WhiteTurn,
				chess.AllCastleRights,
			),
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
			Board: chess.BoardFromRanks(
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
			),
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
			Board: chess.BoardFromRanks(
				[8]string{
					"r   k  r",
					"Pppp ppp",
					" b   nbN",
					"nP      ",
					"BBP P   ",
					"q    N  ",
					"Pp P  PP",
					"R  Q RK ",
				},
				chess.WhiteTurn,
				chess.NewCastleRights(false, false, true, true),
			),
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
			Board: chess.BoardFromRanks(
				[8]string{
					"rnbq k r",
					"pp Pbppp",
					"  p     ",
					"        ",
					"  B     ",
					"        ",
					"PPP NnPP",
					"RNBQK  R",
				},
				chess.WhiteTurn,
				chess.NewCastleRights(true, true, false, false),
			),
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
			Board: chess.BoardFromRanks(
				[8]string{
					"r    rk ",
					" pp qppp",
					"p np n  ",
					"  b p B ",
					"  B P b ",
					"P NP N  ",
					" PP QPPP",
					"R    RK ",
				},
				chess.WhiteTurn,
				chess.NoCastleRights,
			),
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
