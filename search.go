package chess

import (
	"fmt"
	"sort"
	"time"
)

var inf = 9223372036854775807

var nodes = 0

func Search(b Board, seconds int) Move {
	start := time.Now()

	ms, status := b.LegalMoves()

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
		searched = orderMoves(b, d, searched)
		if time.Since(start).Seconds() > float64(seconds) {
			break
		}
	}

	fmt.Println(nodes)
	return searched[0].Move
}

type MoveSearch struct {
	Move  Move
	Eval  int
	Depth int
}

func orderMoves(b Board, depth int, searched []MoveSearch) []MoveSearch {
	for i, m := range searched {
		b.Move(m.Move)
		val := -negamax(b, depth-1, -inf, inf, 1)
		searched[i].Depth = depth
		searched[i].Eval = val
		b.Unmove()
	}

	sort.Slice(searched, func(i, j int) bool { return searched[i].Eval > searched[j].Eval })

	return searched
}

const checkmateEval = -1000000

func negamax(b Board, depth int, alpha, beta, ply int) int {
	nodes++
	ms, status := b.LegalMoves()
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
		return Evaluate(b)
	}
	for _, m := range ms {
		b.Move(m)
		value = max(value, -negamax(b, depth-1, -beta, -alpha, ply+1))
		b.Unmove()
		alpha = max(alpha, value)
		if alpha >= beta {
			break
		}
	}
	return value
}
