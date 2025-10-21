package chess

import "math/bits"

func sumValues(bb uint64, values [64]int) int {
	val := 0
	for bb != 0 {
		i := bits.TrailingZeros64(bb)
		val += values[i]
		bb &= bb - 1
	}
	return val
}

var knightVals = [64]int{
	-50, -40, -30, -30, -30, -30, -40, -50,
	-40, -20, 0, 5, 5, 0, -20, -40,
	-30, 0, 10, 15, 15, 10, 0, -30,
	-30, 5, 15, 20, 20, 15, 5, -30,
	-30, 5, 15, 20, 20, 15, 5, -30,
	-30, 0, 10, 15, 15, 10, 0, -30,
	-40, -20, 0, 5, 5, 0, -20, -40,
	-50, -40, -30, -30, -30, -30, -40, -50,
}

func Evaluate(b Board) int {
	var multiplier int
	if b.Turn == WhiteTurn {
		multiplier = 1
	} else {
		multiplier = -1
	}

	var (
		pawnWt   = 100
		knightWt = 300
		bishopWt = 300
		rookWt   = 500
		queenWt  = 900
		kingWt   = 20000

		wP = bits.OnesCount64(b.whitePawns)
		bP = bits.OnesCount64(b.blackPawns)
		wR = bits.OnesCount64(b.whiteRooks)
		bR = bits.OnesCount64(b.blackRooks)
		wN = bits.OnesCount64(b.whiteKnights)
		bN = bits.OnesCount64(b.blackKnights)
		wB = bits.OnesCount64(b.whiteBishops)
		bB = bits.OnesCount64(b.blackBishops)
		wQ = bits.OnesCount64(b.whiteQueens)
		bQ = bits.OnesCount64(b.blackQueens)
		wK = bits.OnesCount64(b.whiteKings)
		bK = bits.OnesCount64(b.blackKings)
	)

	materialScore := pawnWt*(wP-bP) +
		knightWt*(wN-bN) +
		bishopWt*(wB-bB) +
		rookWt*(wR-bR) +
		queenWt*(wQ-bQ) +
		kingWt*(wK-bK)

	positionalScore := sumValues(b.whiteKnights, knightVals) - sumValues(b.blackKnights, knightVals)

	score := materialScore + positionalScore

	return score * multiplier
}
