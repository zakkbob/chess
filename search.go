package chess

import (
	"fmt"
	"sort"
	"time"
)

var inf = 9223372036854775807

func (e *Engine) Search(seconds int) Move {
	start := time.Now()

	ms, status := e.B.LegalMoves()

	if status != InProgress {
		panic("aghhh, the game is over")
	}

	searched := make([]MoveSearch, 0, len(ms))

	for _, m := range ms {
		searched = append(searched, MoveSearch{
			Move:  m,
			Eval:  0,
			Depth: 0,
		})
	}

	for d := 1; ; d += 1 {
		searched = e.orderMoves(d, searched)
		fmt.Println(d)
		if time.Since(start).Seconds() > float64(seconds) {
			break
		}
	}

	return searched[0].Move
}

type MoveSearch struct {
	Move  Move
	Eval  int
	Depth int
}

func (e *Engine) orderMoves(depth int, searched []MoveSearch) []MoveSearch {
	for i, m := range searched {
		e.B.Move(m.Move)
		val := -e.negamax(depth-1, -inf, inf, 1)
		searched[i].Depth = depth
		searched[i].Eval = val
		e.B.Unmove()
	}

	sort.Slice(searched, func(i, j int) bool { return searched[i].Eval > searched[j].Eval })

	return searched
}

const checkmateEval = -1000000

func (e *Engine) negamax(depth int, alpha, beta, ply int) int {
	z := e.B.Zobrist()
	if t, ok := e.TT.Get(z); ok && t.Depth >= depth {
		if t.Type == ExactEntry ||
			(t.Type == LowerBoundEntry && t.Score >= beta) ||
			(t.Type == UpperBoundEntry && t.Score < alpha) {
			return t.Score
		}
	}

	originalA := alpha
	originalB := beta

	ms, status := e.B.LegalMoves()
	value := -inf

	switch status {
	case Checkmate:
		return checkmateEval + ply
	case Stalemate, Draw:
		return 0
	}

	if len(ms) == 0 {
		return value
	}
	if depth == 0 {
		return e.Evaluate()
	}

	for _, m := range ms {
		e.B.Move(m)
		value = max(value, -e.negamax(depth-1, -beta, -alpha, ply+1))
		e.B.Unmove()
		alpha = max(alpha, value)
		if alpha >= beta {
			break
		}
	}

	t := Transposition{
		Key:   z,
		Depth: depth,
		Score: value,
		Type:  ExactEntry,
	}

	if value <= originalA {
		t.Type = UpperBoundEntry
	} else if value >= originalB {
		t.Type = LowerBoundEntry
	}

	e.TT.Save(t)

	return value
}
