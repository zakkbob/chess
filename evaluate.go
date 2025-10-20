package chess

import "math/bits"

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

	score := materialScore

	return score * multiplier
}
