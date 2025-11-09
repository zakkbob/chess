package chess

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

//func BenchmarkSearch(b *testing.B) {
//bd := NewBoard()
//ms, _ := bd.LegalMoves()
//
//searched := make([]MoveSearch, 0, len(ms))
//for _, m := range ms {
//searched = append(searched, MoveSearch{
//Move:  m,
//Eval:  0,
//Depth: 0,
//})
//}
//
//for b.Loop() {
//b.StopTimer()
//tt := NewTranspositionTable(20)
//b.StartTimer()
//searched = orderMoves(bd, 4, searched, tt)
//}
//}

func TestSearch(t *testing.T) {
	b, err := BoardFromFEN("r2qkb1r/pp2nppp/3p4/2pNN1B1/2BnP3/3P4/PPP2PPP/R2bK2R w KQkq - 1 0")
	assert.NoError(t, err)

	e := Engine{
		B:  b,
		TT: *NewTranspositionTable(10),
		EP: DefaultParams,
	}

	m := e.Search(1)
	t.Log(m.String())
}
