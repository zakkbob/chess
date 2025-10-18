package chess_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zakkbob/chess"
)

func BenchmarkMove(b *testing.B) {
	bd := chess.NewBoard()

	for b.Loop() {
		bd.Move(0)
	}
}

func TestString(t *testing.T) {
	b := chess.NewBoard()
	expected := "rnbqkbnr\npppppppp\n        \n        \n        \n        \nPPPPPPPP\nRNBQKBNR\n"
	got := b.String()

	assert.Equal(t, expected, got)
}

func TestFromRankString(t *testing.T) {
	expected := [8]string{
		"rnbqkbnr",
		"pppppppp",
		"        ",
		"        ",
		"        ",
		"        ",
		"PPPPPPPP",
		"RNBQKBNR",
	}

	b := chess.BoardFromRanks(expected, chess.WhiteTurn)

	assert.Equal(t, expected, b.RankStrings())
}

func TestRankString(t *testing.T) {
	b := chess.NewBoard()
	expected := [8]string{
		"rnbqkbnr",
		"pppppppp",
		"        ",
		"        ",
		"        ",
		"        ",
		"PPPPPPPP",
		"RNBQKBNR",
	}
	got := b.RankStrings()

	assert.Equal(t, expected, got)
}

func TestMoveUnmove(t *testing.T) {
	move := func(from, to int) uint32 {
		return uint32(from<<23) | uint32(to<<17)
	}

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
			move: chess.PawnType | move(9, 17),
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
			move: chess.PawnType | move(49, 33),
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
			move: chess.PawnType | move(36, 43) | chess.EnPassantMask,
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
			move: chess.PawnType | move(25, 16) | chess.EnPassantMask,
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
			move: chess.KingType | move(3, 1) | chess.KingCastle,
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
			move: chess.KingType | move(3, 5) | chess.QueenCastle,
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
			move: chess.KingType | move(59, 57) | chess.KingCastle,
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
			move: chess.KingType | move(59, 61) | chess.QueenCastle,
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
			move: chess.PawnType | move(51, 59) | chess.KnightPromotion,
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
			move: chess.PawnType | move(53, 61) | chess.BishopPromotion,
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
			move: chess.PawnType | move(48, 56) | chess.RookPromotion,
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
			move: chess.PawnType | move(55, 63) | chess.QueenPromotion,
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
			move: chess.PawnType | move(11, 3) | chess.KnightPromotion,
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
			move: chess.PawnType | move(8, 0) | chess.BishopPromotion,
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
			move: chess.PawnType | move(15, 7) | chess.RookPromotion,
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
			move: chess.PawnType | move(13, 5) | chess.RookPromotion,
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
			move: chess.QueenType | move(13, 45) | chess.PawnCapture,
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
			move: chess.QueenType | move(13, 45) | chess.RookCapture,
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
			move: chess.QueenType | move(13, 45) | chess.KnightCapture,
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
			move: chess.QueenType | move(13, 45) | chess.BishopCapture,
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
			move: chess.QueenType | move(13, 45) | chess.QueenCapture,
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
			move: chess.QueenType | move(13, 45) | chess.PawnCapture,
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
			move: chess.QueenType | move(13, 45) | chess.RookCapture,
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
			move: chess.QueenType | move(13, 45) | chess.KnightCapture,
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
			move: chess.QueenType | move(13, 45) | chess.BishopCapture,
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
			move: chess.QueenType | move(13, 45) | chess.QueenCapture,
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
			})
			t.Run("Unmove", func(t *testing.T) {
				b.Unmove(tt.move)
				assert.Equal(t, tt.board, b.RankStrings())
			})
		})
	}
}
