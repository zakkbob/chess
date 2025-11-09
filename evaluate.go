package chess

import "math/bits"

func sumWhiteValues(bb uint64, values [64]int) int {
	val := 0
	for bb != 0 {
		i := bits.TrailingZeros64(bb)
		val += values[63-i]
		bb &= bb - 1
	}
	return val
}

func sumBlackValues(bb uint64, values [64]int) int {
	val := 0
	for bb != 0 {
		i := bits.TrailingZeros64(bb)
		val += values[i]
		bb &= bb - 1
	}
	return val
}

func (e *Engine) Evaluate() int {
	var multiplier int
	if e.B.Turn == WhiteTurn {
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

		wP = bits.OnesCount64(e.B.whitePawns)
		bP = bits.OnesCount64(e.B.blackPawns)
		wR = bits.OnesCount64(e.B.whiteRooks)
		bR = bits.OnesCount64(e.B.blackRooks)
		wN = bits.OnesCount64(e.B.whiteKnights)
		bN = bits.OnesCount64(e.B.blackKnights)
		wB = bits.OnesCount64(e.B.whiteBishops)
		bB = bits.OnesCount64(e.B.blackBishops)
		wQ = bits.OnesCount64(e.B.whiteQueens)
		bQ = bits.OnesCount64(e.B.blackQueens)
		wK = bits.OnesCount64(e.B.whiteKings)
		bK = bits.OnesCount64(e.B.blackKings)
	)

	materialScore := pawnWt*(wP-bP) +
		knightWt*(wN-bN) +
		bishopWt*(wB-bB) +
		rookWt*(wR-bR) +
		queenWt*(wQ-bQ) +
		kingWt*(wK-bK)

	positionalScore := sumWhiteValues(e.B.whitePawns, e.EP.PawnVals) - sumBlackValues(e.B.blackPawns, e.EP.PawnVals) +
		sumWhiteValues(e.B.whiteKnights, e.EP.KnightVals) - sumBlackValues(e.B.blackKnights, e.EP.KnightVals) +
		sumWhiteValues(e.B.whiteRooks, e.EP.RookVals) - sumBlackValues(e.B.blackKnights, e.EP.RookVals) +
		sumWhiteValues(e.B.whiteKings, e.EP.KingVals) - sumBlackValues(e.B.blackKings, e.EP.KingVals) +
		sumWhiteValues(e.B.whiteBishops, e.EP.BishopVals) - sumBlackValues(e.B.blackBishops, e.EP.BishopVals) +
		sumWhiteValues(e.B.whiteQueens, e.EP.QueenVals) - sumBlackValues(e.B.blackQueens, e.EP.QueenVals)

	score := materialScore + positionalScore

	return score * multiplier
}
