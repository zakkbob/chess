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

func BenchmarkZobrist(b *testing.B) {
	bd := chess.NewBoard()
	for b.Loop() {
		bd.Zobrist()
	}
}

func TestString(t *testing.T) {
	b := chess.NewBoard()
	expected := "8|♜ ♞ ♝ ♛ ♚ ♝ ♞ ♜ \n7|♟ ♟ ♟ ♟ ♟ ♟ ♟ ♟ \n6|. . . . . . . . \n5|. . . . . . . . \n4|. . . . . . . . \n3|. . . . . . . . \n2|♙ ♙ ♙ ♙ ♙ ♙ ♙ ♙ \n1|♖ ♘ ♗ ♕ ♔ ♗ ♘ ♖ \n  ---------------\n  a b c d e f g h"
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

	b := chess.BoardFromRanks(expected, chess.WhiteTurn, chess.AllCastleRights)

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
