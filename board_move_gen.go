package chess

import (
	"math/bits"
)

func (b *Board) LegalMoves() []Move {
	ms := make([]Move, 0, 128)

	if b.Turn == WhiteTurn {
		var (
			pawns   = b.whitePawns
			rooks   = b.whiteRooks
			knights = b.whiteKnights
			bishops = b.whiteBishops
			queens  = b.whiteQueens
			kings   = b.whiteKings

			enemyPawns   = b.blackPawns
			enemyRooks   = b.blackRooks
			enemyKnights = b.blackKnights
			enemyBishops = b.blackBishops
			enemyQueens  = b.blackQueens
			enemyKings   = b.blackKings

			own      = pawns | rooks | knights | bishops | queens | kings
			enemies  = enemyPawns | enemyRooks | enemyKnights | enemyBishops | enemyQueens | enemyKings
			occupied = own | enemies
			empty    = ^occupied
		)

		var (
			rank0 uint64 = 0x00000000000000ff
			rank1 uint64 = 0x000000000000ff00
			//rank2 uint64 = 0x0000000000ff0000
			rank3 uint64 = 0x00000000ff000000
			//rank4 uint64 = 0x000000ff00000000
			//rank5 uint64 = 0x0000ff0000000000
			rank6 uint64 = 0x00ff000000000000
			rank7 uint64 = 0xff00000000000000
			file0 uint64 = 0x0101010101010101
			file1 uint64 = 0x0202020202020202
			//file2 uint64 = 0x0404040404040404
			//file3 uint64 = 0x0808080808080808
			//file4 uint64 = 0x1010101010101010
			//file5 uint64 = 0x2020202020202020
			file6 uint64 = 0x4040404040404040
			file7 uint64 = 0x8080808080808080
		)

		addMoves := func(cells uint64, from int, pieceType PieceType, promotion Promotion, capture Capture, enPassant bool, castle Castle) {
			for cells != 0 {
				i := bits.TrailingZeros64(cells)
				ms = append(ms, NewMove(from, i, pieceType, promotion, capture, enPassant, castle))
				cells &= cells - 1
			}
		}

		addPawnMove := func(cells uint64, from int, capture Capture, enPassant bool, castle Castle) {
			if from/8 == 6 { // promotions
				addMoves(cells, from, PawnType, RookPromotion, capture, enPassant, castle)
				addMoves(cells, from, PawnType, KnightPromotion, capture, enPassant, castle)
				addMoves(cells, from, PawnType, BishopPromotion, capture, enPassant, castle)
				addMoves(cells, from, PawnType, QueenPromotion, capture, enPassant, castle)

			} else { // non promotions
				addMoves(cells, from, PawnType, NoPromotion, capture, enPassant, castle)
			}
		}

		forEachEnemyBoard := func(getCaptures func(enemies uint64, c Capture)) {
			getCaptures(enemyPawns, PawnCapture)
			getCaptures(enemyRooks, RookCapture)
			getCaptures(enemyKnights, KnightCapture)
			getCaptures(enemyBishops, BishopCapture)
			getCaptures(enemyQueens, QueenCapture)
		}

		addMovesAndCaptures := func(cells uint64, from int, pieceType PieceType, promotion Promotion, enPassant bool, castle Castle) {
			addMoves(cells&empty, from, pieceType, promotion, NoCapture, enPassant, castle)

			forEachEnemyBoard(func(enemies uint64, c Capture) {
				addMoves(cells&enemies, from, pieceType, promotion, c, enPassant, castle)
			})

		}

		leftRay := func(i int, shift int, stopPropagating uint64) uint64 {
			piece := uint64(1) << i

			piece = ((^(stopPropagating | enemies) & piece) << shift) & ^own
			moves := piece
			piece = ((^(stopPropagating | enemies) & piece) << shift) & ^own
			moves |= piece
			piece = ((^(stopPropagating | enemies) & piece) << shift) & ^own
			moves |= piece
			piece = ((^(stopPropagating | enemies) & piece) << shift) & ^own
			moves |= piece
			piece = ((^(stopPropagating | enemies) & piece) << shift) & ^own
			moves |= piece
			piece = ((^(stopPropagating | enemies) & piece) << shift) & ^own
			moves |= piece
			piece = ((^(stopPropagating | enemies) & piece) << shift) & ^own
			moves |= piece

			return moves
		}

		rightRay := func(i int, shift int, stopPropagating uint64) uint64 {
			piece := uint64(1) << i

			piece = ((^(stopPropagating | enemies) & piece) >> shift) & ^own
			moves := piece
			piece = ((^(stopPropagating | enemies) & piece) >> shift) & ^own
			moves |= piece
			piece = ((^(stopPropagating | enemies) & piece) >> shift) & ^own
			moves |= piece
			piece = ((^(stopPropagating | enemies) & piece) >> shift) & ^own
			moves |= piece
			piece = ((^(stopPropagating | enemies) & piece) >> shift) & ^own
			moves |= piece
			piece = ((^(stopPropagating | enemies) & piece) >> shift) & ^own
			moves |= piece
			piece = ((^(stopPropagating | enemies) & piece) >> shift) & ^own
			moves |= piece

			return moves
		}

		northWestRay := func(i int) uint64 {
			return leftRay(i, 9, rank7|file7)
		}

		northRay := func(i int) uint64 {
			return leftRay(i, 8, rank7)
		}

		northEastRay := func(i int) uint64 {
			return leftRay(i, 7, rank7|file0)
		}

		westRay := func(i int) uint64 {
			return leftRay(i, 1, file7)
		}

		eastRay := func(i int) uint64 {
			return rightRay(i, 1, file0)
		}

		southWestRay := func(i int) uint64 {
			return rightRay(i, 7, rank0|file7)
		}

		southRay := func(i int) uint64 {
			return rightRay(i, 8, rank0)
		}

		southEastRay := func(i int) uint64 {
			return rightRay(i, 9, rank0|file0)
		}

		rookRays := func(i int) uint64 {
			return northRay(i) | eastRay(i) | southRay(i) | westRay(i)
		}

		bishopRays := func(i int) uint64 {
			return northEastRay(i) | northWestRay(i) | southEastRay(i) | southWestRay(i)
		}

		// Pawns
		p := pawns
		for p != 0 {
			i := bits.TrailingZeros64(p)
			rank := i / 8
			file := i % 8

			board := uint64(1) << i

			// Pushes
			pushes := (board << 8)          // single push
			pushes |= (board << 16) & rank3 // double push
			pushes &= empty                 // remove occupied squares
			addPawnMove(pushes, i, NoCapture, false, NoCastle)

			// Captures
			forEachEnemyBoard(func(enemies uint64, c Capture) {
				captures := (board & ^file7) << 9 // capture left
				captures |= (board & ^file0) << 7 // capture right
				captures &= enemies               // only allow enemy captures
				addPawnMove(captures, i, c, false, NoCastle)
			})

			// En passant
			if b.CanEnPassant && rank == 4 && ((b.EnPassantFile == file-1) || (b.EnPassantFile == file+1)) {
				ms = append(ms, NewMove(i, Index(5, b.EnPassantFile), PawnType, NoPromotion, NoCapture, true, NoCastle))
			}

			p &= p - 1
		}

		// Rook
		r := rooks
		for r != 0 {
			i := bits.TrailingZeros64(r)
			moves := rookRays(i)
			addMovesAndCaptures(moves, i, RookType, NoPromotion, false, NoCastle)
			r &= r - 1
		}

		// Bishop
		b := bishops
		for b != 0 {
			i := bits.TrailingZeros64(b)
			moves := bishopRays(i)
			addMovesAndCaptures(moves, i, BishopType, NoPromotion, false, NoCastle)
			b &= b - 1
		}

		// Queen
		q := queens
		for q != 0 {
			i := bits.TrailingZeros64(q)
			moves := rookRays(i) | bishopRays(i)
			addMovesAndCaptures(moves, i, QueenType, NoPromotion, false, NoCastle)
			q &= q - 1
		}

		// Knight
		n := knights
		for n != 0 {
			i := bits.TrailingZeros64(n)

			board := uint64(1) << i

			moves := (^(rank7 | file1 | file0) & board) << 6
			moves |= (^(rank7 | file6 | file7) & board) << 10
			moves |= (^(rank7 | rank6 | file0) & board) << 15
			moves |= (^(rank7 | rank6 | file7) & board) << 17
			moves |= (^(rank0 | file6 | file7) & board) >> 6
			moves |= (^(rank0 | file1 | file0) & board) >> 10
			moves |= (^(rank0 | rank1 | file7) & board) >> 15
			moves |= (^(rank0 | rank1 | file1) & board) >> 17

			addMovesAndCaptures(moves, i, KnightType, NoPromotion, false, NoCastle)

			n &= n - 1
		}

		// King
		k := kings
		for k != 0 {
			i := bits.TrailingZeros64(k)

			board := uint64(1) << i

			moves := (board & ^rank7) << 8
			moves |= (board & ^file7) << 1
			moves |= (board & ^rank0) >> 8
			moves |= (board & ^file0) >> 1
			moves |= (board & ^rank7 & ^file7) << 9
			moves |= (board & ^file7 & ^file0) << 7
			moves |= (board & ^rank0 & ^file0) >> 9
			moves |= (board & ^file0 & ^file7) >> 7

			addMovesAndCaptures(moves, i, KingType, NoPromotion, false, NoCastle)

			k &= k - 1
		}
	}

	return ms
}
