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
