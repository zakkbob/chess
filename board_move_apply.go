package chess

// WARN: must follow format <from square><to square>[<promoted to>]
func (b *Board) DoAlgebraicMove(s string) {
	if len(s) != 4 && len(s) != 5 {
		panic("Invalid algebraic move")
	}

	from := IndexFromAlgebraic(s[0:2])
	to := IndexFromAlgebraic(s[2:4])
	promotion := NoPromotion
	if len(s) == 5 {
		promotion = PromotionFromSymbol(rune(s[4]))
	}

	b.DoCoordinateMove(from, to, promotion)

	//fmt.Println(s)
	//fmt.Println(b.CastleRights.CanWhiteKing())
	//fmt.Println(b.CastleRights.CanWhiteQueen())
	//fmt.Println(b.CastleRights.CanBlackKing())
	//fmt.Println(b.CastleRights.CanBlackQueen())
	//fmt.Println()
}

// Applies given move
// Assumes it is valid and legal
func (b *Board) DoCoordinateMove(from, to int, promotion Promotion) {
	pieceType := b.pieceType(from)
	capture := b.pieceType(to).ToCapture()

	fromRank := from / 8
	fromFile := from % 8
	toRank := to / 8
	toFile := to % 8

	fileDiff := toFile - fromFile
	enPassant := (pieceType == PawnType) && (fileDiff == 1 || fileDiff == -1) && ((toRank == 2 && fromRank == 3) || (toRank == 5 && fromRank == 4))

	castle := NoCastle

	if pieceType == KingType {
		if (from == 3 && to == 1) || (from == 59 && to == 57) {
			castle = KingCastle
		} else if (from == 3 && to == 5) || (from == 59 && to == 61) {
			castle = QueenCastle
		}
	}

	m := NewMove(
		from,
		to,
		pieceType,
		promotion,
		capture,
		enPassant,
		b.CastleRights,
		castle,
	)

	b.Move(m)
}

// Applies given move
// Assumes it is valid and legal
func (b *Board) Move(m Move) {
	b.Moves = append(b.Moves, m)
	b.HalfMoves++

	var from uint32 = m.From()
	var to uint32 = m.To()

	var fromMask = uint64(1) << from
	var toMask = uint64(1) << to

	var moveMask uint64 = fromMask | toMask

	if m.IsNoisy() {
		b.noisyMoves = append(b.noisyMoves, b.HalfMoves-1)
	}

	if m.IsDoublePush() {
		b.CanEnPassant = true
		b.EnPassantFile = int(m.ToFile())
	} else {
		b.CanEnPassant = false
	}

	if b.Turn == WhiteTurn {
		// Move piece
		switch m.PieceType() {
		case PawnType:
			b.whitePawns ^= moveMask
		case RookType:
			if m.From() == 0 {
				b.CastleRights.LoseWhiteKing()
			} else if m.From() == 7 {
				b.CastleRights.LoseWhiteQueen()
			}
			b.whiteRooks ^= moveMask
		case KnightType:
			b.whiteKnights ^= moveMask
		case BishopType:
			b.whiteBishops ^= moveMask
		case QueenType:
			b.whiteQueens ^= moveMask
		case KingType:
			b.CastleRights.LoseWhiteKing()
			b.CastleRights.LoseWhiteQueen()
			b.whiteKings ^= moveMask
		}

		// Capture piece
		switch m.Capture() {
		case PawnCapture:
			b.blackPawns ^= toMask
		case RookCapture:
			if m.To() == 56 {
				b.CastleRights.LoseBlackKing()
			} else if m.To() == 63 {
				b.CastleRights.LoseBlackQueen()
			}
			b.blackRooks ^= toMask
		case KnightCapture:
			b.blackKnights ^= toMask
		case BishopCapture:
			b.blackBishops ^= toMask
		case QueenCapture:
			b.blackQueens ^= toMask
		}

		// Handle promotion
		switch m.Promotion() {
		case RookPromotion:
			b.whitePawns ^= toMask
			b.whiteRooks |= toMask
		case KnightPromotion:
			b.whitePawns ^= toMask
			b.whiteKnights |= toMask
		case BishopPromotion:
			b.whitePawns ^= toMask
			b.whiteBishops |= toMask
		case QueenPromotion:
			b.whitePawns ^= toMask
			b.whiteQueens |= toMask
		}

		// Handle castling
		switch m.Castle() {
		case KingCastle:
			b.whiteRooks ^= 0b00000101
		case QueenCastle:
			b.whiteRooks ^= 0b10010000
		}

		// Handle en passant
		if m.EnPassant() {
			b.blackPawns ^= toMask >> 8
		}
	} else {
		// Move piece
		switch m.PieceType() {
		case PawnType:
			b.blackPawns ^= moveMask
		case RookType:
			if m.From() == 56 {
				b.CastleRights.LoseBlackKing()
			} else if m.From() == 63 {
				b.CastleRights.LoseBlackQueen()
			}
			b.blackRooks ^= moveMask
		case KnightType:
			b.blackKnights ^= moveMask
		case BishopType:
			b.blackBishops ^= moveMask
		case QueenType:
			b.blackQueens ^= moveMask
		case KingType:
			b.CastleRights.LoseBlackKing()
			b.CastleRights.LoseBlackQueen()
			b.blackKings ^= moveMask
		}

		// Capture piece
		switch m.Capture() {
		case PawnCapture:
			b.whitePawns ^= toMask
		case RookCapture:
			if m.To() == 0 {
				b.CastleRights.LoseWhiteKing()
			} else if m.To() == 7 {
				b.CastleRights.LoseWhiteQueen()
			}
			b.whiteRooks ^= toMask
		case KnightCapture:
			b.whiteKnights ^= toMask
		case BishopCapture:
			b.whiteBishops ^= toMask
		case QueenCapture:
			b.whiteQueens ^= toMask
		}

		// Handle promotion
		switch m.Promotion() {
		case RookPromotion:
			b.blackPawns ^= toMask
			b.blackRooks |= toMask
		case KnightPromotion:
			b.blackPawns ^= toMask
			b.blackKnights |= toMask
		case BishopPromotion:
			b.blackPawns ^= toMask
			b.blackBishops |= toMask
		case QueenPromotion:
			b.blackPawns ^= toMask
			b.blackQueens |= toMask
		}

		// Handle castling
		switch m.Castle() {
		case KingCastle:
			b.blackRooks ^= 0b00000101 << 56
		case QueenCastle:
			b.blackRooks ^= 0b10010000 << 56
		}

		// Handle en passant
		if m.EnPassant() {
			b.whitePawns ^= toMask << 8
		}
	}

	b.Turn = !b.Turn
}

func (b *Board) Unmove() {
	b.Turn = !b.Turn

	if len(b.noisyMoves) != 0 {
		lastNoisyMove := b.noisyMoves[len(b.noisyMoves)-1]
		if lastNoisyMove == b.HalfMoves-1 {
			b.noisyMoves = b.noisyMoves[:len(b.noisyMoves)-1]
		}
	}

	m := b.Moves[b.HalfMoves-1]
	b.Moves = b.Moves[:b.HalfMoves-1]
	b.HalfMoves--

	b.CastleRights = m.CastleRights()

	if b.HalfMoves != 0 {
		lastMove := b.Moves[b.HalfMoves-1]
		b.CanEnPassant = lastMove.IsDoublePush()
		b.EnPassantFile = int(lastMove.ToFile())
	} else {
		b.CanEnPassant = false
	}

	var from uint32 = m.From()
	var to uint32 = m.To()

	var fromMask = uint64(1) << from
	var toMask = uint64(1) << to

	var moveMask uint64 = fromMask | toMask

	if b.Turn == WhiteTurn {
		// Undo promotion
		switch m.Promotion() {
		case RookPromotion:
			b.whitePawns |= toMask
			b.whiteRooks ^= toMask
		case KnightPromotion:
			b.whitePawns |= toMask
			b.whiteKnights ^= toMask
		case BishopPromotion:
			b.whitePawns |= toMask
			b.whiteBishops ^= toMask
		case QueenPromotion:
			b.whitePawns |= toMask
			b.whiteQueens ^= toMask
		}

		// Unmove piece
		switch m.PieceType() {
		case PawnType:
			b.whitePawns ^= moveMask
		case RookType:
			b.whiteRooks ^= moveMask
		case KnightType:
			b.whiteKnights ^= moveMask
		case BishopType:
			b.whiteBishops ^= moveMask
		case QueenType:
			b.whiteQueens ^= moveMask
		case KingType:
			b.whiteKings ^= moveMask
		}

		// Uncapture piece
		switch m.Capture() {
		case PawnCapture:
			b.blackPawns |= toMask
		case RookCapture:
			b.blackRooks |= toMask
		case KnightCapture:
			b.blackKnights |= toMask
		case BishopCapture:
			b.blackBishops |= toMask
		case QueenCapture:
			b.blackQueens |= toMask
		}

		// Uncastle
		switch m.Castle() {
		case KingCastle:
			b.whiteRooks ^= 0b00000101
		case QueenCastle:
			b.whiteRooks ^= 0b10010000
		}

		// Undo en passant
		if m.EnPassant() {
			b.blackPawns ^= toMask >> 8
		}
	} else {
		// Undo promotion
		switch m.Promotion() {
		case RookPromotion:
			b.blackPawns |= toMask
			b.blackRooks ^= toMask
		case KnightPromotion:
			b.blackPawns |= toMask
			b.blackKnights ^= toMask
		case BishopPromotion:
			b.blackPawns |= toMask
			b.blackBishops ^= toMask
		case QueenPromotion:
			b.blackPawns |= toMask
			b.blackQueens ^= toMask
		}

		// Unmove piece
		switch m.PieceType() {
		case PawnType:
			b.blackPawns ^= moveMask
		case RookType:
			b.blackRooks ^= moveMask
		case KnightType:
			b.blackKnights ^= moveMask
		case BishopType:
			b.blackBishops ^= moveMask
		case QueenType:
			b.blackQueens ^= moveMask
		case KingType:
			b.blackKings ^= moveMask
		}

		// Uncapture piece
		switch m.Capture() {
		case PawnCapture:
			b.whitePawns |= toMask
		case RookCapture:
			b.whiteRooks |= toMask
		case KnightCapture:
			b.whiteKnights |= toMask
		case BishopCapture:
			b.whiteBishops |= toMask
		case QueenCapture:
			b.whiteQueens |= toMask
		}

		// Uncastle
		switch m.Castle() {
		case KingCastle:
			b.blackRooks ^= 0b00000101 << 56
		case QueenCastle:
			b.blackRooks ^= 0b10010000 << 56
		}

		// Undo en passant
		if m.EnPassant() {
			b.whitePawns ^= toMask << 8
		}
	}
}
