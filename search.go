package chess

func Search(b Board, depth int) Move {
	ms := b.LegalMoves()

	if len(ms) == 0 {
		panic("aghhh, the game is over")
	}

	bestMove := ms[0]
	bestValue := -9223372036854775808

	for _, m := range ms {
		b.Move(m)
		val := -negamax(b, depth-1)
		if val > bestValue {
			bestValue = val
			bestMove = m
		}
		b.Unmove()
	}

	return bestMove
}

func negamax(b Board, depth int) int {
	ms := b.LegalMoves()
	value := -9223372036854775808
	if len(ms) == 0 {
		return value
	}
	if depth == 0 {
		return Evaluate(b)
	}
	for _, m := range ms {
		b.Move(m)
		value = max(value, -negamax(b, depth-1))
		b.Unmove()
	}
	return value
}
