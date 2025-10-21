package chess

var inf = 9223372036854775807

func Search(b Board, depth int) Move {
	ms := b.LegalMoves()

	if len(ms) == 0 {
		panic("aghhh, the game is over")
	}

	bestMove := ms[0]
	bestValue := -inf

	for _, m := range ms {
		b.Move(m)
		val := -negamax(b, depth-1, -inf, inf)
		if val > bestValue {
			bestValue = val
			bestMove = m
		}
		b.Unmove()
	}

	return bestMove
}

func negamax(b Board, depth int, alpha, beta int) int {
	ms := b.LegalMoves()
	value := -inf
	if len(ms) == 0 {
		return value
	}
	if depth == 0 {
		return Evaluate(b)
	}
	for _, m := range ms {
		b.Move(m)
		value = max(value, -negamax(b, depth-1, -beta, -alpha))
		b.Unmove()
		alpha = max(alpha, value)
		if alpha >= beta {
			break
		}
	}
	return value
}
