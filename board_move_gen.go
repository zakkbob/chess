package chess

import (
	"fmt"
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
		)

		var (
			rank3 uint64 = 0x00000000ff000000
			file0 uint64 = 0x0101010101010101
			file7 uint64 = 0x8080808080808080
		)

		addMoves := func(cells uint64, from int, pieceType PieceType, promotion Promotion, capture Capture, enPassant bool, castle Castle) {
			for cells != 0 {
				i := bits.TrailingZeros64(cells)
				ms = append(ms, NewMove(from, i, pieceType, promotion, capture, enPassant, castle))
				cells &= cells - 1
			}
		}

		addPawnMoves := func(cells uint64, from int, capture Capture, enPassant bool, castle Castle) {
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
			pushes &= ^occupied             // remove occupied squares
			addPawnMoves(pushes, i, NoCapture, false, NoCastle)

			// Captures
			forEachEnemyBoard(func(enemies uint64, c Capture) {
				captures := (board & ^file7) << 9 // capture left
				captures |= (board & ^file0) << 7 // capture right
				captures &= enemies               // only allow enemy captures
				addPawnMoves(captures, i, c, false, NoCastle)
			})

			// En passant
			fmt.Println("e", file, b.EnPassantFile)
			if b.CanEnPassant && rank == 4 && ((b.EnPassantFile == file-1) || (b.EnPassantFile == file+1)) {
				fmt.Println("a")
				ms = append(ms, NewMove(i, Index(5, b.EnPassantFile), PawnType, NoPromotion, NoCapture, true, NoCastle))
			}

			p &= p - 1
		}
	}

	return ms
}
