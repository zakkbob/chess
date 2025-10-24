package chess

import (
	"math/bits"
)

type GameStatus int

const (
	InProgress = iota
	Draw
	Stalemate
	Checkmate
)

// also checks for draw/stalemate/checkmate
func (b *Board) LegalMoves() ([]Move, GameStatus) {
	ms := make([]Move, 0, 128)

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
	)

	if b.Turn == BlackTurn {
		pawns = b.blackPawns
		rooks = b.blackRooks
		knights = b.blackKnights
		bishops = b.blackBishops
		queens = b.blackQueens
		kings = b.blackKings

		enemyPawns = b.whitePawns
		enemyRooks = b.whiteRooks
		enemyKnights = b.whiteKnights
		enemyBishops = b.whiteBishops
		enemyQueens = b.whiteQueens
		enemyKings = b.whiteKings
	}

	var (
		own                      = pawns | rooks | knights | bishops | queens | kings
		enemies                  = enemyPawns | enemyRooks | enemyKnights | enemyBishops | enemyQueens | enemyKings
		diagonalSlidingEnemies   = enemyBishops | enemyQueens
		orthogonalSlidingEnemies = enemyRooks | enemyQueens
		occupied                 = own | enemies
		empty                    = ^occupied
	)

	var (
		rank0 uint64 = 0x00000000000000ff
		rank1 uint64 = 0x000000000000ff00
		//rank2 uint64 = 0x0000000000ff0000
		rank3 uint64 = 0x00000000ff000000
		rank4 uint64 = 0x000000ff00000000
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

	var (
		pins           [64]uint64
		permittedMoves uint64 = ^uint64(0)
	)

	addMoves := func(cells uint64, from int, pieceType PieceType, promotion Promotion, capture Capture, enPassant bool, castle Castle) {
		cells &= ^pins[from] & permittedMoves
		for cells != 0 {
			i := bits.TrailingZeros64(cells)
			ms = append(ms, NewMove(from, i, pieceType, promotion, capture, enPassant, b.CastleRights, castle))
			cells &= cells - 1
		}
	}

	addPawnMove := func(cells uint64, from int, capture Capture, enPassant bool, castle Castle) {
		promotionRank := 6
		if b.Turn == BlackTurn {
			promotionRank = 1
		}
		if from/8 == promotionRank { // promotions
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

	// casts a ray in a direction which increases index
	// stopPropogating is the bitboard which contains 1's the ray should stop at, but the ray will include that position
	// blockPropogation is the bitboard which contains 1's the ray should stop at, but the ray won't include that position
	leftRay := func(i int, shift int, stopPropagating uint64, blockPropogation uint64) uint64 {
		piece := uint64(1) << i

		piece = ((^(stopPropagating) & piece) << shift) & ^blockPropogation
		moves := piece
		piece = ((^(stopPropagating) & piece) << shift) & ^blockPropogation
		moves |= piece
		piece = ((^(stopPropagating) & piece) << shift) & ^blockPropogation
		moves |= piece
		piece = ((^(stopPropagating) & piece) << shift) & ^blockPropogation
		moves |= piece
		piece = ((^(stopPropagating) & piece) << shift) & ^blockPropogation
		moves |= piece
		piece = ((^(stopPropagating) & piece) << shift) & ^blockPropogation
		moves |= piece
		piece = ((^(stopPropagating) & piece) << shift) & ^blockPropogation
		moves |= piece

		return moves
	}

	// casts a ray in a direction which decreases index
	// stopPropogating is the bitboard which contains 1's the ray should stop at, but the ray will include that position
	// blockPropogation is the bitboard which contains 1's the ray should stop at, but the ray won't include that position
	rightRay := func(i int, shift int, stopPropagation uint64, blockPropogation uint64) uint64 {
		piece := uint64(1) << i

		piece = ((^(stopPropagation) & piece) >> shift) & ^blockPropogation
		moves := piece
		piece = ((^(stopPropagation) & piece) >> shift) & ^blockPropogation
		moves |= piece
		piece = ((^(stopPropagation) & piece) >> shift) & ^blockPropogation
		moves |= piece
		piece = ((^(stopPropagation) & piece) >> shift) & ^blockPropogation
		moves |= piece
		piece = ((^(stopPropagation) & piece) >> shift) & ^blockPropogation
		moves |= piece
		piece = ((^(stopPropagation) & piece) >> shift) & ^blockPropogation
		moves |= piece
		piece = ((^(stopPropagation) & piece) >> shift) & ^blockPropogation
		moves |= piece

		return moves
	}

	northWestRay := func(i int, stopPropogation uint64, blockPropogation uint64) uint64 {
		return leftRay(i, 9, stopPropogation|rank7|file7, blockPropogation)
	}

	northRay := func(i int, stopPropogation uint64, blockPropogation uint64) uint64 {
		return leftRay(i, 8, stopPropogation|rank7, blockPropogation)
	}

	northEastRay := func(i int, stopPropogation uint64, blockPropogation uint64) uint64 {
		return leftRay(i, 7, stopPropogation|rank7|file0, blockPropogation)
	}

	westRay := func(i int, stopPropogation uint64, blockPropogation uint64) uint64 {
		return leftRay(i, 1, stopPropogation|file7, blockPropogation)
	}

	eastRay := func(i int, stopPropogation uint64, blockPropogation uint64) uint64 {
		return rightRay(i, 1, stopPropogation|file0, blockPropogation)
	}

	southWestRay := func(i int, stopPropogation uint64, blockPropogation uint64) uint64 {
		return rightRay(i, 7, stopPropogation|rank0|file7, blockPropogation)
	}

	southRay := func(i int, stopPropogation uint64, blockPropogation uint64) uint64 {
		return rightRay(i, 8, stopPropogation|rank0, blockPropogation)
	}

	southEastRay := func(i int, stopPropogation uint64, blockPropogation uint64) uint64 {
		return rightRay(i, 9, stopPropogation|rank0|file0, blockPropogation)
	}

	orthogonalRays := func(i int, stopPropogation uint64, blockPropogation uint64) uint64 {
		return northRay(i, stopPropogation, blockPropogation) | eastRay(i, stopPropogation, blockPropogation) | southRay(i, stopPropogation, blockPropogation) | westRay(i, stopPropogation, blockPropogation)
	}

	diagonalRays := func(i int, stopPropogation uint64, blockPropogation uint64) uint64 {
		return northEastRay(i, stopPropogation, blockPropogation) | northWestRay(i, stopPropogation, blockPropogation) | southEastRay(i, stopPropogation, blockPropogation) | southWestRay(i, stopPropogation, blockPropogation)
	}

	// Pin detection
	bb := kings
	for bb != 0 {
		i := bits.TrailingZeros64(bb)

		checkForPin := func(rayFunc func(int, uint64, uint64) uint64, slidingEnemies uint64) {
			firstRay := rayFunc(i, own&^kings, enemies)
			ownIndex := bits.TrailingZeros64(firstRay & own)
			secondRay := rayFunc(ownIndex, enemies, own)
			pinner := secondRay & slidingEnemies
			if pinner != 0 {
				pins[ownIndex] = ^(firstRay | secondRay)
			}
		}

		checkForPin(northRay, orthogonalSlidingEnemies)
		checkForPin(eastRay, orthogonalSlidingEnemies)
		checkForPin(southRay, orthogonalSlidingEnemies)
		checkForPin(westRay, orthogonalSlidingEnemies)
		checkForPin(northWestRay, diagonalSlidingEnemies)
		checkForPin(northEastRay, diagonalSlidingEnemies)
		checkForPin(southWestRay, diagonalSlidingEnemies)
		checkForPin(southEastRay, diagonalSlidingEnemies)

		bb &= bb - 1
	}

	var enPassantPawnIndex int
	var enPassantRank int
	if b.Turn == WhiteTurn {
		enPassantRank = 4
	} else {
		enPassantRank = 3
	}
	if b.CanEnPassant {
		enPassantPawnIndex = Index(enPassantRank, b.EnPassantFile)
	} else {
		enPassantPawnIndex = 64
	}
	enPassantPawn := uint64(1) << enPassantPawnIndex

	enPassantPawnIsOnlyPieceAttackingKing := false //Wow thats a long name

	// pieces under enemy attack
	var (
		enemyAttackedSquares uint64
	)
	{
		// Pawns
		bb = enemyPawns & ^enPassantPawn
		for bb != 0 {
			i := bits.TrailingZeros64(bb)

			board := uint64(1) << i

			// Captures
			var captures uint64
			if b.Turn == BlackTurn {
				captures |= (board & ^file7) << 9 // capture left
				captures |= (board & ^file0) << 7 // capture right
			} else {
				captures |= (board & ^file7) >> 7 // capture left
				captures |= (board & ^file0) >> 9 // capture right
			}
			enemyAttackedSquares |= captures

			if kings&captures != 0 {
				permittedMoves &= board
			}

			bb &= bb - 1
		}

		// Rook
		bb = enemyRooks
		for bb != 0 {
			i := bits.TrailingZeros64(bb)
			board := uint64(1) << i
			moves := orthogonalRays(i, (own|enemies)&^(board|kings), 0) //WARN: doesn't handle properly, current the ray passes through the king
			enemyAttackedSquares |= moves
			bb &= bb - 1
		}

		// Bishop
		bb = enemyBishops
		for bb != 0 {
			i := bits.TrailingZeros64(bb)
			board := uint64(1) << i
			moves := diagonalRays(i, (own|enemies)&^(board|kings), 0)
			enemyAttackedSquares |= moves
			bb &= bb - 1
		}

		// Queen
		bb = enemyQueens
		for bb != 0 {
			i := bits.TrailingZeros64(bb)
			board := uint64(1) << i
			moves := orthogonalRays(i, (own|enemies)&^(board|kings), 0) | diagonalRays(i, (own|enemies)&^(board|kings), 0)
			enemyAttackedSquares |= moves
			bb &= bb - 1
		}

		// Knight
		bb = enemyKnights
		for bb != 0 {
			i := bits.TrailingZeros64(bb)

			board := uint64(1) << i

			moves := (^(rank7 | file1 | file0) & board) << 6
			moves |= (^(rank7 | file6 | file7) & board) << 10
			moves |= (^(rank7 | rank6 | file0) & board) << 15
			moves |= (^(rank7 | rank6 | file7) & board) << 17
			moves |= (^(rank0 | file6 | file7) & board) >> 6
			moves |= (^(rank0 | file1 | file0) & board) >> 10
			moves |= (^(rank0 | rank1 | file7) & board) >> 15
			moves |= (^(rank0 | rank1 | file0) & board) >> 17

			enemyAttackedSquares |= moves

			if kings&moves != 0 {
				permittedMoves &= board
			}

			bb &= bb - 1
		}

		// King
		bb = enemyKings
		i := bits.TrailingZeros64(bb)

		board := uint64(1) << i

		moves := (^(rank7) & board) << 8
		moves |= (^(file7) & board) << 1
		moves |= (^(rank0) & board) >> 8
		moves |= (^(file0) & board) >> 1
		moves |= (^(rank7 | file7) & board) << 9
		moves |= (^(rank7 | file0) & board) << 7
		moves |= (^(rank0 | file0) & board) >> 9
		moves |= (^(rank0 | file7) & board) >> 7

		enemyAttackedSquares |= moves

		// handle en passant capturable pawn seperately to check if the enPassantPawn is the only one attacking the king
		bb = enPassantPawn
		if bb != 0 {
			// Captures
			var captures uint64
			if b.Turn == BlackTurn {
				captures |= (enPassantPawn & ^file7) << 9 // capture left
				captures |= (enPassantPawn & ^file0) << 7 // capture right
			} else {
				captures |= (enPassantPawn & ^file7) >> 7 // capture left
				captures |= (enPassantPawn & ^file0) >> 9 // capture right
			}

			if (kings&captures != 0) && (enemyAttackedSquares&kings == 0) {
				enPassantPawnIsOnlyPieceAttackingKing = true //Wow thats a long name
			}

			enemyAttackedSquares |= captures

			if kings&captures != 0 {
				permittedMoves &= enPassantPawn
			}
		}
	}

	// Blockable check detection
	bb = kings
	for bb != 0 {
		i := bits.TrailingZeros64(bb)

		checkForSlidingCheck := func(rayFunc func(int, uint64, uint64) uint64, slidingEnemies uint64) {
			ray := rayFunc(i, slidingEnemies, occupied&^(kings|slidingEnemies))
			attacker := ray & slidingEnemies
			if attacker != 0 {
				permittedMoves &= ray
			}

		}

		checkForSlidingCheck(northRay, orthogonalSlidingEnemies)
		checkForSlidingCheck(eastRay, orthogonalSlidingEnemies)
		checkForSlidingCheck(southRay, orthogonalSlidingEnemies)
		checkForSlidingCheck(westRay, orthogonalSlidingEnemies)
		checkForSlidingCheck(northWestRay, diagonalSlidingEnemies)
		checkForSlidingCheck(northEastRay, diagonalSlidingEnemies)
		checkForSlidingCheck(southWestRay, diagonalSlidingEnemies)
		checkForSlidingCheck(southEastRay, diagonalSlidingEnemies)

		bb &= bb - 1
	}

	// Check if en passant will put king in check
	enPassantPutsKingInCheck := false
	if b.CanEnPassant {
		caputuringPawn1 := (uint64(1) << Index(enPassantRank, b.EnPassantFile-1)) & pawns
		capturingPawn2 := (uint64(1) << Index(enPassantRank, b.EnPassantFile+1)) & pawns

		// if two pawns can do enPassant. then an east/west ray will always be blocked
		var capturingPawn uint64
		if caputuringPawn1 != 0 && capturingPawn2 != 0 {
			capturingPawn = caputuringPawn1
		} else {
			capturingPawn = caputuringPawn1 | capturingPawn2
		}

		capturedPiece := uint64(1) << uint64(enPassantPawnIndex)
		kingIndex := bits.TrailingZeros64(kings)
		fileDiff := enPassantPawnIndex%8 - kingIndex%8
		rankDiff := enPassantPawnIndex/8 - kingIndex/8

		var ray uint64 = 0
		switch {
		case fileDiff == 0 && rankDiff > 0 && b.Turn != WhiteTurn: // north, on white turn, an enpassant will block nort ray
			ray = northRay(kingIndex, occupied&^(kings|capturedPiece), 0) & orthogonalSlidingEnemies
		case fileDiff == 0 && rankDiff < 0 && b.Turn != BlackTurn: // south, on black turn, an en passant will block a south rau
			ray = southRay(kingIndex, occupied&^(kings|capturedPiece), 0) & orthogonalSlidingEnemies
		case rankDiff == 0 && fileDiff < 0: // east
			ray = eastRay(kingIndex, occupied&^(kings|capturingPawn|capturedPiece), 0) & orthogonalSlidingEnemies
		case rankDiff == 0 && fileDiff > 0: // west
			ray = westRay(kingIndex, occupied&^(kings|capturingPawn|capturedPiece), 0) & orthogonalSlidingEnemies
		case rankDiff == fileDiff && rankDiff > 0: // north-west
			ray = northWestRay(kingIndex, occupied&^(kings|capturedPiece), 0) & diagonalSlidingEnemies
		case rankDiff == -fileDiff && rankDiff > 0: // north-east
			ray = northEastRay(kingIndex, occupied&^(kings|capturedPiece), 0) & diagonalSlidingEnemies
		case rankDiff == fileDiff && rankDiff < 0: // south-east
			ray = southEastRay(kingIndex, occupied&^(kings|capturedPiece), 0) & diagonalSlidingEnemies
		case rankDiff == -fileDiff && rankDiff < 0: // south-west
			ray = southWestRay(kingIndex, occupied&^(kings|capturedPiece), 0) & diagonalSlidingEnemies
		default:

		}
		if ray != 0 {
			enPassantPutsKingInCheck = true
		}
	}

	// Pawns
	bb = pawns
	for bb != 0 {
		i := bits.TrailingZeros64(bb)
		rank := i / 8
		file := i % 8

		board := uint64(1) << i

		// Pushes
		var pushes uint64
		if b.Turn == WhiteTurn {
			pushes = (board << 8)                           // single push
			pushes |= (((board << 8) & empty) << 8) & rank3 // double push
		} else {
			pushes = (board >> 8)                           // single push
			pushes |= (((board >> 8) & empty) >> 8) & rank4 // double push

		}
		pushes &= empty // remove occupied squares
		addPawnMove(pushes, i, NoCapture, false, NoCastle)

		// Captures
		forEachEnemyBoard(func(enemies uint64, c Capture) {
			var captures uint64
			if b.Turn == WhiteTurn {
				captures |= (board & ^file7) << 9 // capture left
				captures |= (board & ^file0) << 7 // capture right
			} else {
				captures |= (board & ^file7) >> 7 // capture left
				captures |= (board & ^file0) >> 9 // capture right
			}
			captures &= enemies // only allow enemy captures
			addPawnMove(captures, i, c, false, NoCastle)
		})

		if !enPassantPutsKingInCheck {
			if b.Turn == WhiteTurn {
				// En passant
				to := Index(5, b.EnPassantFile)
				permittedMove := permittedMoves&(uint64(1)<<to) != 0 || enPassantPawnIsOnlyPieceAttackingKing
				pinned := pins[i]&(uint64(1)<<to) != 0
				if b.CanEnPassant && rank == 4 && ((b.EnPassantFile == file-1) || (b.EnPassantFile == file+1)) && !pinned && permittedMove {
					ms = append(ms, NewMove(i, to, PawnType, NoPromotion, NoCapture, true, b.CastleRights, NoCastle))
				}
			} else {
				// En passant
				to := Index(2, b.EnPassantFile)
				permittedMove := permittedMoves&(uint64(1)<<to) != 0 || enPassantPawnIsOnlyPieceAttackingKing
				pinned := pins[i]&(uint64(1)<<to) != 0
				if b.CanEnPassant && rank == 3 && ((b.EnPassantFile == file-1) || (b.EnPassantFile == file+1)) && !pinned && permittedMove {
					ms = append(ms, NewMove(i, to, PawnType, NoPromotion, NoCapture, true, b.CastleRights, NoCastle))
				}

			}
		}

		bb &= bb - 1
	}

	// Rook
	bb = rooks
	for bb != 0 {
		i := bits.TrailingZeros64(bb)
		moves := orthogonalRays(i, enemies, own)
		addMovesAndCaptures(moves, i, RookType, NoPromotion, false, NoCastle)
		bb &= bb - 1
	}

	// Bishop
	bb = bishops
	for bb != 0 {
		i := bits.TrailingZeros64(bb)
		moves := diagonalRays(i, enemies, own)
		addMovesAndCaptures(moves, i, BishopType, NoPromotion, false, NoCastle)
		bb &= bb - 1
	}

	// Queen
	bb = queens
	for bb != 0 {
		i := bits.TrailingZeros64(bb)
		moves := orthogonalRays(i, enemies, own) | diagonalRays(i, enemies, own)
		addMovesAndCaptures(moves, i, QueenType, NoPromotion, false, NoCastle)
		bb &= bb - 1
	}

	// Knight
	bb = knights
	for bb != 0 {
		i := bits.TrailingZeros64(bb)

		board := uint64(1) << i

		moves := (^(rank7 | file1 | file0) & board) << 6
		moves |= (^(rank7 | file6 | file7) & board) << 10
		moves |= (^(rank7 | rank6 | file0) & board) << 15
		moves |= (^(rank7 | rank6 | file7) & board) << 17
		moves |= (^(rank0 | file6 | file7) & board) >> 6
		moves |= (^(rank0 | file1 | file0) & board) >> 10
		moves |= (^(rank0 | rank1 | file7) & board) >> 15
		moves |= (^(rank0 | rank1 | file0) & board) >> 17

		addMovesAndCaptures(moves, i, KnightType, NoPromotion, false, NoCastle)

		bb &= bb - 1
	}

	addKingMoves := func(cells uint64, from int, capture Capture) {
		for cells != 0 {
			i := bits.TrailingZeros64(cells)
			ms = append(ms, NewMove(from, i, KingType, NoPromotion, capture, false, b.CastleRights, NoCastle))
			cells &= cells - 1
		}
	}

	addKingMovesAndCaptures := func(cells uint64, from int) {
		addKingMoves(cells&empty, from, NoCapture)

		forEachEnemyBoard(func(enemies uint64, c Capture) {
			addKingMoves(cells&enemies, from, c)
		})
	}

	// King
	bb = kings
	i := bits.TrailingZeros64(bb)

	board := uint64(1) << i

	moves := (^(rank7) & board) << 8
	moves |= (^(file7) & board) << 1
	moves |= (^(rank0) & board) >> 8
	moves |= (^(file0) & board) >> 1
	moves |= (^(rank7 | file7) & board) << 9
	moves |= (^(rank7 | file0) & board) << 7
	moves |= (^(rank0 | file0) & board) >> 9
	moves |= (^(rank0 | file7) & board) >> 7
	moves &= ^enemyAttackedSquares

	addKingMovesAndCaptures(moves, i)

	// Castling
	if enemyAttackedSquares&kings == 0 {
		if b.Turn == WhiteTurn {
			if b.CastleRights.CanWhiteKing() && ((occupied|enemyAttackedSquares)&0b00000110 == 0) && (i == 3) {
				ms = append(ms, NewMove(i, 1, KingType, NoPromotion, NoCapture, false, b.CastleRights, KingCastle))
			}
			if b.CastleRights.CanWhiteQueen() && (occupied&0b01110000 == 0) && (enemyAttackedSquares&0b00110000 == 0) && (i == 3) {
				ms = append(ms, NewMove(i, 5, KingType, NoPromotion, NoCapture, false, b.CastleRights, QueenCastle))
			}
		} else {
			if b.CastleRights.CanBlackKing() && ((occupied|enemyAttackedSquares)&(0b00000110<<56) == 0) && (i == 59) {
				ms = append(ms, NewMove(i, 57, KingType, NoPromotion, NoCapture, false, b.CastleRights, KingCastle))
			}
			if b.CastleRights.CanBlackQueen() && (occupied&(0b01110000<<56) == 0) && (enemyAttackedSquares&(0b00110000<<56) == 0) && (i == 59) {
				ms = append(ms, NewMove(i, 61, KingType, NoPromotion, NoCapture, false, b.CastleRights, QueenCastle))
			}
		}

	}

	// Game status
	if len(ms) == 0 {
		if enemyAttackedSquares&kings != 0 {
			return ms, Checkmate
		} else {
			return ms, Stalemate
		}
	}

	if b.QuietMoveCounter() == 50 {
		return ms, Draw
	}

	return ms, InProgress
}
