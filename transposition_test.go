package chess_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zakkbob/chess"
)

// bit of a janky test
func TestTranspositionTable(t *testing.T) {
	t1 := chess.Transposition{
		Key:      0,
		BestMove: 0,
		Depth:    1,
		Score:    1,
		Type:     chess.ExactEntry,
	}

	t2 := chess.Transposition{
		Key:      2,
		BestMove: 0,
		Depth:    2,
		Score:    2,
		Type:     chess.ExactEntry,
	}

	tt := chess.NewTranspositionTable(1) // 2 entries

	tt.Save(t1)
	got, ok := tt.Get(0)
	assert.Equal(t, true, ok)
	assert.Equal(t, t1, got)

	_, ok = tt.Get(2)
	assert.Equal(t, false, ok)

	tt.Save(t2)
	_, ok = tt.Get(0)
	assert.Equal(t, false, ok)

	got, ok = tt.Get(2)
	assert.Equal(t, true, ok)
	assert.Equal(t, t2, got)

	got, ok = tt.Get(1)
	assert.Equal(t, false, ok)
}
