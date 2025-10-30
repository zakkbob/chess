package chess

import (
	"testing"
)

func BenchmarkSearch(b *testing.B) {
	bd := NewBoard()
	ms, _ := bd.LegalMoves()

	searched := make([]MoveSearch, 0, len(ms))
	for _, m := range ms {
		searched = append(searched, MoveSearch{
			Move:  m,
			Eval:  0,
			Depth: 0,
		})
	}

	for b.Loop() {
		b.StopTimer()
		tt := NewTranspositionTable(20)
		b.StartTimer()
		searched = orderMoves(bd, 4, searched, tt)
	}
}
