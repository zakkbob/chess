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

var pawnVals = [64]int{
	00, 00, 00, 00, 00, 00, 00, 00,
	10, 20, 20, 30, 30, 20, 20, 10,
	0, 10, 10, 20, 20, 10, 10, 00,
	-10, 00, 00, 10, 10, 00, 00, -10,
	-20, -10, -10, 00, 00, -10, -10, -20,
	-30, -20, -20, -20, -20, -20, -20, -30,
	-40, -30, -30, -30, -30, -30, -30, -40,
	00, 00, 00, 00, 00, 00, 00, 00,
}

var rookVals = [64]int{
	-20, -20, -10, 00, 00, -10, -20, -20,
	-20, -10, 00, 10, 10, 00, -10, -20,
	-10, 00, 10, 20, 20, 10, 00, -10,
	00, 10, 20, 30, 30, 20, 10, 00,
	00, 10, 20, 30, 30, 20, 10, 00,
	-10, 00, 10, 20, 20, 10, 00, -10,
	-20, -10, 00, 10, 10, 00, -10, -20,
	-20, -20, -10, 00, 00, -10, -20, -20,
}

var queenVals = [64]int{
	-20, -20, -10, 00, 00, -10, -20, -20,
	-20, -10, 00, 10, 10, 00, -10, -20,
	-10, 00, 10, 20, 20, 10, 00, -10,
	00, 10, 20, 30, 30, 20, 10, 00,
	00, 10, 20, 30, 30, 20, 10, 00,
	-10, 00, 10, 20, 20, 10, 00, -10,
	-20, -10, 00, 10, 10, 00, -10, -20,
	-20, -20, -10, 00, 00, -10, -20, -20,
}

var kingVals = [64]int{
	-50, -50, -50, -50, -50, -50, -50, -50,
	-40, -40, -50, -50, -50, -50, -40, -40,
	-30, -30, -40, -50, -50, -40, -30, -30,
	-20, -20, -30, -40, -40, -30, -20, -20,
	-10, -10, -20, -30, -30, -20, -10, -10,
	0, 0, -10, -20, -20, -10, 0, 0,
	10, 10, 0, -10, -10, 0, 10, 10,
	20, 20, 10, 0, 0, 10, 20, 20,
}

var knightVals = [64]int{
	-50, -40, -30, -30, -30, -30, -40, -50,
	-40, -20, 00, 5, 5, 00, -20, -40,
	-30, 00, 10, 15, 15, 10, 00, -30,
	-30, 5, 15, 20, 20, 15, 5, -30,
	-30, 5, 15, 20, 20, 15, 5, -30,
	-30, 00, 10, 15, 15, 10, 00, -30,
	-40, -20, 00, 5, 5, 00, -20, -40,
	-50, -40, -30, -30, -30, -30, -40, -50,
}

var bishopVals = [64]int{
	-50, -40, -30, -30, -30, -30, -40, -50,
	-40, -20, 00, 5, 5, 00, -20, -40,
	-30, 00, 10, 15, 15, 10, 00, -30,
	-30, 5, 15, 20, 20, 15, 5, -30,
	-30, 5, 15, 20, 20, 15, 5, -30,
	-30, 00, 10, 15, 15, 10, 00, -30,
	-40, -20, 00, 5, 5, 00, -20, -40,
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

	positionalScore := sumWhiteValues(b.whitePawns, pawnVals) - sumBlackValues(b.blackPawns, pawnVals) +
		sumWhiteValues(b.whiteKnights, knightVals) - sumBlackValues(b.blackKnights, knightVals) +
		sumWhiteValues(b.whiteRooks, rookVals) - sumBlackValues(b.blackKnights, rookVals) +
		sumWhiteValues(b.whiteKings, kingVals) - sumBlackValues(b.blackKings, kingVals) +
		sumWhiteValues(b.whiteBishops, bishopVals) - sumBlackValues(b.blackBishops, bishopVals) +
		sumWhiteValues(b.whiteQueens, queenVals) - sumBlackValues(b.blackQueens, queenVals)

	score := materialScore + positionalScore

	return score * multiplier
}
