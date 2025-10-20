package chess_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zakkbob/chess"
)

// not needed, but i don't want to rewrite the tests
func move(from, to int) chess.Move {
	return chess.Move(from<<23) | chess.Move(to<<17)
}

func TestMoveCounters(t *testing.T) {
	tests := []struct {
		Move             chess.Move
		HalfMoves        int
		QuietMoveCounter int
	}{

		//rnbqkbnr
		//pppppppp
		//
		//
		//
		//
		//PPPPPPPP
		//RNBQKBNR
		{
			Move:             chess.Move(chess.PawnType) | move(9, 17),
			HalfMoves:        1,
			QuietMoveCounter: 0,
		},
		//rnbqkbnr
		//pppppppp
		//
		//
		//
		//      P
		//PPPPPP P
		//RNBQKBNR
		{
			Move:             chess.Move(chess.KnightType) | move(62, 47),
			HalfMoves:        2,
			QuietMoveCounter: 1,
		},
		//r bqkbnr
		//pppppppp
		//n
		//
		//
		//      P
		//PPPPPP P
		//RNBQKBNR
		{
			Move:             chess.Move(chess.BishopType) | move(2, 16),
			HalfMoves:        3,
			QuietMoveCounter: 2,
		},
		//r bqkbnr
		//pppppppp
		//n
		//
		//
		//      PB
		//PPPPPP P
		//RNBQK NR
		{
			Move:             chess.Move(chess.PawnType) | move(48, 40),
			HalfMoves:        4,
			QuietMoveCounter: 0,
		},
		//r bqkbnr
		//ppppppp
		//n      p
		//
		//
		//      PB
		//PPPPPP P
		//RNBQK NR
	}

	b := chess.NewBoard()

	for i, tt := range tests {
		b.Move(tt.Move)

		t.Logf("Board after move %d\n%s", i, b.String())

		require.Equal(t, b.HalfMoves, tt.HalfMoves, "Halfmove counter wrong for move %d", i+1)
		require.Equal(t, b.QuietMoveCounter(), tt.QuietMoveCounter, "Quiet move counter wrong for move %d", i+1)
	}

	for i := len(tests) - 2; i >= 0; i-- {
		tt := tests[i]
		b.Unmove()

		t.Logf("Board after unmove %d\n%s", i, b.String())

		require.Equal(t, b.HalfMoves, tt.HalfMoves, "Halfmove counter wrong for unmove %d", i+1)
		require.Equal(t, b.QuietMoveCounter(), tt.QuietMoveCounter, "Quiet move counter wrong for unmove %d", i+1)

		// Add an extra test each step, so the history is non-linear
		if b.Turn == chess.WhiteTurn {
			b.Move(chess.Move(chess.PawnType) | move(15, 23))

		} else {
			b.Move(chess.Move(chess.PawnType) | move(55, 47))
		}

		require.Equal(t, b.HalfMoves, tt.HalfMoves+1, "Halfmove counter wrong for pawn push (after unmove %d)", i+1)
		require.Equal(t, b.QuietMoveCounter(), 0, "Quiet move counter wrong for pawn push (after unmove %d)", i+1)

		b.Unmove()
		require.Equal(t, b.HalfMoves, tt.HalfMoves, "Halfmove counter wrong for unmove pawn push (returning to move %d)", i+1)
		require.Equal(t, b.QuietMoveCounter(), tt.QuietMoveCounter, "Quiet move counter wrong for unmove pawn push (returning to move %d)", i+1)
	}

}
