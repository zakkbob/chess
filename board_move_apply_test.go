package chess_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestDoublePush(t *testing.T) {
	b := chess.NewBoard()

	assert.Equal(t, false, b.CanEnPassant)

	b.Move(chess.NewMove(9, 25, chess.PawnType, chess.NoPromotion, chess.NoCapture, false, chess.NoCastle))
	assert.Equal(t, true, b.CanEnPassant)
	assert.Equal(t, 1, b.EnPassantFile)

	b.Move(chess.NewMove(55, 47, chess.PawnType, chess.NoPromotion, chess.NoCapture, false, chess.NoCastle))
	assert.Equal(t, false, b.CanEnPassant)

	b.Unmove()
	assert.Equal(t, true, b.CanEnPassant)
	assert.Equal(t, 1, b.EnPassantFile)

	b.Unmove()
	assert.Equal(t, false, b.CanEnPassant)
}

func TestMoveUnmove(t *testing.T) {
	tests := []struct {
		name     string
		turn     chess.Turn
		board    [8]string
		move     chess.Move
		expected [8]string
	}{
		{
			name: "White pawn advance one",
			turn: chess.WhiteTurn,
			board: [8]string{
				"rnbqkbnr",
				"pppppppp",
				"        ",
				"        ",
				"        ",
				"        ",
				"PPPPPPPP",
				"RNBQKBNR",
			},
			move: chess.Move(chess.PawnType) | move(9, 17),
			expected: [8]string{
				"rnbqkbnr",
				"pppppppp",
				"        ",
				"        ",
				"        ",
				"      P ",
				"PPPPPP P",
				"RNBQKBNR",
			},
		},
		{
			name: "Black pawn advance two",
			turn: chess.BlackTurn,
			board: [8]string{
				"rnbqkbnr",
				"pppppppp",
				"        ",
				"        ",
				"        ",
				"        ",
				"PPPPPPPP",
				"RNBQKBNR",
			},
			move: chess.Move(chess.PawnType) | move(49, 33),
			expected: [8]string{
				"rnbqkbnr",
				"pppppp p",
				"        ",
				"      p ",
				"        ",
				"        ",
				"PPPPPPPP",
				"RNBQKBNR",
			},
		},
		{
			name: "En passant from white",
			turn: chess.WhiteTurn,
			board: [8]string{
				"rnbqkbnr",
				"pppp ppp",
				"        ",
				"   Pp   ",
				"        ",
				"        ",
				"PPP PPPP",
				"RNBQKBNR",
			},
			move: chess.Move(chess.PawnType) | move(36, 43) | chess.Move(chess.EnPassantMask),
			expected: [8]string{
				"rnbqkbnr",
				"pppp ppp",
				"    P   ",
				"        ",
				"        ",
				"        ",
				"PPP PPPP",
				"RNBQKBNR",
			},
		},
		{
			name: "En passant from black",
			turn: chess.BlackTurn,
			board: [8]string{
				"rnbqkbnr",
				"pppppp p",
				"        ",
				"        ",
				"      pP",
				"        ",
				"PPPPPPP ",
				"RNBQKBNR",
			},
			move: chess.Move(chess.PawnType) | move(25, 16) | chess.Move(chess.EnPassantMask),
			expected: [8]string{
				"rnbqkbnr",
				"pppppp p",
				"        ",
				"        ",
				"        ",
				"       p",
				"PPPPPPP ",
				"RNBQKBNR",
			},
		},
		{
			name: "White castle king side",
			turn: chess.WhiteTurn,
			board: [8]string{
				"rnbqkbnr",
				"pppppppp",
				"        ",
				"        ",
				"        ",
				"        ",
				"PPPPPPPP",
				"RNBQK  R",
			},
			move: chess.Move(chess.KingType) | move(3, 1) | chess.Move(chess.KingCastle),
			expected: [8]string{
				"rnbqkbnr",
				"pppppppp",
				"        ",
				"        ",
				"        ",
				"        ",
				"PPPPPPPP",
				"RNBQ RK ",
			},
		},
		{
			name: "White castle queen side",
			turn: chess.WhiteTurn,
			board: [8]string{
				"rnbqkbnr",
				"pppppppp",
				"        ",
				"        ",
				"        ",
				"        ",
				"PPPPPPPP",
				"R   KBNR",
			},
			move: chess.Move(chess.KingType) | move(3, 5) | chess.Move(chess.QueenCastle),
			expected: [8]string{
				"rnbqkbnr",
				"pppppppp",
				"        ",
				"        ",
				"        ",
				"        ",
				"PPPPPPPP",
				"  KR BNR",
			},
		},
		{
			name: "Black castle king side",
			turn: chess.BlackTurn,
			board: [8]string{
				"rnbqk  r",
				"pppppppp",
				"        ",
				"        ",
				"        ",
				"        ",
				"PPPPPPPP",
				"RNBQKBNR",
			},
			move: chess.Move(chess.KingType) | move(59, 57) | chess.Move(chess.KingCastle),
			expected: [8]string{
				"rnbq rk ",
				"pppppppp",
				"        ",
				"        ",
				"        ",
				"        ",
				"PPPPPPPP",
				"RNBQKBNR",
			},
		},
		{
			name: "Black castle queen side",
			turn: chess.BlackTurn,
			board: [8]string{
				"r   kbnr",
				"pppppppp",
				"        ",
				"        ",
				"        ",
				"        ",
				"PPPPPPPP",
				"RNBQKBNR",
			},
			move: chess.Move(chess.KingType) | move(59, 61) | chess.Move(chess.QueenCastle),
			expected: [8]string{
				"  kr bnr",
				"pppppppp",
				"        ",
				"        ",
				"        ",
				"        ",
				"PPPPPPPP",
				"RNBQKBNR",
			},
		},
		{
			name: "White promote knight",
			turn: chess.WhiteTurn,
			board: [8]string{
				"        ",
				"    P   ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
			},
			move: chess.Move(chess.PawnType) | move(51, 59) | chess.Move(chess.KnightPromotion),
			expected: [8]string{
				"    N   ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
			},
		},
		{
			name: "White promote bishop",
			turn: chess.WhiteTurn,
			board: [8]string{
				"        ",
				"  P     ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
			},
			move: chess.Move(chess.PawnType) | move(53, 61) | chess.Move(chess.BishopPromotion),
			expected: [8]string{
				"  B     ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
			},
		},
		{
			name: "White promote rook",
			turn: chess.WhiteTurn,
			board: [8]string{
				"        ",
				"       P",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
			},
			move: chess.Move(chess.PawnType) | move(48, 56) | chess.Move(chess.RookPromotion),
			expected: [8]string{
				"       R",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
			},
		},
		{
			name: "White promote queen",
			turn: chess.WhiteTurn,
			board: [8]string{
				"        ",
				"P       ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
			},
			move: chess.Move(chess.PawnType) | move(55, 63) | chess.Move(chess.QueenPromotion),
			expected: [8]string{
				"Q       ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
			},
		},
		{
			name: "Black promote knight",
			turn: chess.BlackTurn,
			board: [8]string{
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
				"    p   ",
				"        ",
			},
			move: chess.Move(chess.PawnType) | move(11, 3) | chess.Move(chess.KnightPromotion),
			expected: [8]string{
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
				"    n   ",
			},
		},
		{
			name: "Black promote bishop",
			turn: chess.BlackTurn,
			board: [8]string{
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
				"       p",
				"        ",
			},
			move: chess.Move(chess.PawnType) | move(8, 0) | chess.Move(chess.BishopPromotion),
			expected: [8]string{
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
				"       b",
			},
		},
		{
			name: "Black promote rook",
			turn: chess.BlackTurn,
			board: [8]string{
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
				"p       ",
				"        ",
			},
			move: chess.Move(chess.PawnType) | move(15, 7) | chess.Move(chess.RookPromotion),
			expected: [8]string{
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
				"r       ",
			},
		},
		{
			name: "Black promote queen",
			turn: chess.BlackTurn,
			board: [8]string{
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
				"  p     ",
				"        ",
			},
			move: chess.Move(chess.PawnType) | move(13, 5) | chess.Move(chess.RookPromotion),
			expected: [8]string{
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
				"  r     ",
			},
		},
		{
			name: "Capture black pawn",
			turn: chess.WhiteTurn,
			board: [8]string{
				"        ",
				"        ",
				"  p     ",
				"        ",
				"        ",
				"        ",
				"  Q     ",
				"        ",
			},
			move: chess.Move(chess.QueenType) | move(13, 45) | chess.Move(chess.PawnCapture),
			expected: [8]string{
				"        ",
				"        ",
				"  Q     ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
			},
		},
		{
			name: "Capture black rook",
			turn: chess.WhiteTurn,
			board: [8]string{
				"        ",
				"        ",
				"  r     ",
				"        ",
				"        ",
				"        ",
				"  Q     ",
				"        ",
			},
			move: chess.Move(chess.QueenType) | move(13, 45) | chess.Move(chess.RookCapture),
			expected: [8]string{
				"        ",
				"        ",
				"  Q     ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
			},
		},
		{
			name: "Capture black knight",
			turn: chess.WhiteTurn,
			board: [8]string{
				"        ",
				"        ",
				"  n     ",
				"        ",
				"        ",
				"        ",
				"  Q     ",
				"        ",
			},
			move: chess.Move(chess.QueenType) | move(13, 45) | chess.Move(chess.KnightCapture),
			expected: [8]string{
				"        ",
				"        ",
				"  Q     ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
			},
		},
		{
			name: "Capture black bishop",
			turn: chess.WhiteTurn,
			board: [8]string{
				"        ",
				"        ",
				"  b     ",
				"        ",
				"        ",
				"        ",
				"  Q     ",
				"        ",
			},
			move: chess.Move(chess.QueenType) | move(13, 45) | chess.Move(chess.BishopCapture),
			expected: [8]string{
				"        ",
				"        ",
				"  Q     ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
			},
		},
		{
			name: "Capture black queen",
			turn: chess.WhiteTurn,
			board: [8]string{
				"        ",
				"        ",
				"  q     ",
				"        ",
				"        ",
				"        ",
				"  Q     ",
				"        ",
			},
			move: chess.Move(chess.QueenType) | move(13, 45) | chess.Move(chess.QueenCapture),
			expected: [8]string{
				"        ",
				"        ",
				"  Q     ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
			},
		},
		{
			name: "Capture white pawn",
			turn: chess.BlackTurn,
			board: [8]string{
				"        ",
				"        ",
				"  P     ",
				"        ",
				"        ",
				"        ",
				"  q     ",
				"        ",
			},
			move: chess.Move(chess.QueenType) | move(13, 45) | chess.Move(chess.PawnCapture),
			expected: [8]string{
				"        ",
				"        ",
				"  q     ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
			},
		},
		{
			name: "Capture white rook",
			turn: chess.BlackTurn,
			board: [8]string{
				"        ",
				"        ",
				"  R     ",
				"        ",
				"        ",
				"        ",
				"  q     ",
				"        ",
			},
			move: chess.Move(chess.QueenType) | move(13, 45) | chess.Move(chess.RookCapture),
			expected: [8]string{
				"        ",
				"        ",
				"  q     ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
			},
		},
		{
			name: "Capture white knight",
			turn: chess.BlackTurn,
			board: [8]string{
				"        ",
				"        ",
				"  N     ",
				"        ",
				"        ",
				"        ",
				"  q     ",
				"        ",
			},
			move: chess.Move(chess.QueenType) | move(13, 45) | chess.Move(chess.KnightCapture),
			expected: [8]string{
				"        ",
				"        ",
				"  q     ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
			},
		},
		{
			name: "Capture white bishop",
			turn: chess.BlackTurn,
			board: [8]string{
				"        ",
				"        ",
				"  B     ",
				"        ",
				"        ",
				"        ",
				"  q     ",
				"        ",
			},
			move: chess.Move(chess.QueenType) | move(13, 45) | chess.Move(chess.BishopCapture),
			expected: [8]string{
				"        ",
				"        ",
				"  q     ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
			},
		},
		{
			name: "Capture white queen",
			turn: chess.BlackTurn,
			board: [8]string{
				"        ",
				"        ",
				"  Q     ",
				"        ",
				"        ",
				"        ",
				"  q     ",
				"        ",
			},
			move: chess.Move(chess.QueenType) | move(13, 45) | chess.Move(chess.QueenCapture),
			expected: [8]string{
				"        ",
				"        ",
				"  q     ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := chess.BoardFromRanks(tt.board, tt.turn)
			t.Run("Move", func(t *testing.T) {
				b.Move(tt.move)
				assert.Equal(t, tt.expected, b.RankStrings())
				assert.Equal(t, !tt.turn, b.Turn)
			})
			t.Run("Unmove", func(t *testing.T) {
				b.Unmove()
				assert.Equal(t, tt.board, b.RankStrings())
				assert.Equal(t, tt.turn, b.Turn)
			})

			t.Run("Move (By coordinates)", func(t *testing.T) {
				b.DoCoordinateMove(int(tt.move.From()), int(tt.move.To()), tt.move.Promotion()) // hacky fix till i rewrite tests
				t.Log(tt.move.From(), tt.move.To(), tt.move.FromRank(), tt.move.FromFile(), tt.move.ToRank(), tt.move.ToFile(), tt.move.Promotion())
				assert.Equal(t, tt.expected, b.RankStrings())
				assert.Equal(t, !tt.turn, b.Turn)
			})
			t.Run("Unmove (By coordinates)", func(t *testing.T) {
				b.Unmove()
				assert.Equal(t, tt.board, b.RankStrings())
				assert.Equal(t, tt.turn, b.Turn)
			})
		})
	}
}
