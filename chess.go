package chess

import (
	"strings"
)

type Turn = bool

const (
	WhiteTurn Turn = false
	BlackTurn Turn = true
)

type Board struct {
	whitePawns   uint64
	whiteRooks   uint64
	whiteKnights uint64
	whiteBishops uint64
	whiteQueens  uint64
	whiteKings   uint64

	blackPawns   uint64
	blackRooks   uint64
	blackKnights uint64
	blackBishops uint64
	blackQueens  uint64
	blackKings   uint64

	Turn       Turn
	HalfMoves  int
	Moves      []Move
	noisyMoves []int
}

// Returns a board in the proper starting configuration
func NewBoard() Board {
	return Board{
		whitePawns:   0b0000000000000000000000000000000000000000000000001111111100000000,
		whiteRooks:   0b0000000000000000000000000000000000000000000000000000000010000001,
		whiteKnights: 0b0000000000000000000000000000000000000000000000000000000001000010,
		whiteBishops: 0b0000000000000000000000000000000000000000000000000000000000100100,
		whiteQueens:  0b0000000000000000000000000000000000000000000000000000000000010000,
		whiteKings:   0b0000000000000000000000000000000000000000000000000000000000001000,
		blackPawns:   0b0000000011111111000000000000000000000000000000000000000000000000,
		blackRooks:   0b1000000100000000000000000000000000000000000000000000000000000000,
		blackKnights: 0b0100001000000000000000000000000000000000000000000000000000000000,
		blackBishops: 0b0010010000000000000000000000000000000000000000000000000000000000,
		blackQueens:  0b0001000000000000000000000000000000000000000000000000000000000000,
		blackKings:   0b0000100000000000000000000000000000000000000000000000000000000000,

		Turn:       WhiteTurn,
		HalfMoves:  0,
		Moves:      []Move{},
		noisyMoves: []int{},
	}
}

func BoardFromRanks(rs [8]string, turn Turn) Board {
	b := Board{
		Turn:       turn,
		HalfMoves:  0,
		Moves:      []Move{},
		noisyMoves: []int{},
	}

	for i, r := range rs {
		for j, p := range r {
			posMask := uint64(1) << ((7-i)*8 + (7 - j)) // annoying backwards board thing

			switch p {
			case 'P':
				b.whitePawns |= posMask
			case 'R':
				b.whiteRooks |= posMask
			case 'N':
				b.whiteKnights |= posMask
			case 'B':
				b.whiteBishops |= posMask
			case 'Q':
				b.whiteQueens |= posMask
			case 'K':
				b.whiteKings |= posMask
			case 'p':
				b.blackPawns |= posMask
			case 'r':
				b.blackRooks |= posMask
			case 'n':
				b.blackKnights |= posMask
			case 'b':
				b.blackBishops |= posMask
			case 'q':
				b.blackQueens |= posMask
			case 'k':
				b.blackKings |= posMask
			}
		}
	}

	return b
}

func (b *Board) QuietMoveCounter() int {
	if len(b.noisyMoves) == 0 {
		return 0
	}
	return b.HalfMoves - b.noisyMoves[len(b.noisyMoves)-1] - 1
}

func (b *Board) pieceType(i int) PieceType {
	posMask := uint64(1) << i

	var (
		pawns   = b.whitePawns | b.blackPawns
		rooks   = b.whiteRooks | b.blackRooks
		knights = b.whiteKnights | b.blackKnights
		bishops = b.whiteBishops | b.blackBishops
		queens  = b.whiteQueens | b.blackQueens
		kings   = b.whiteKings | b.blackKings
	)

	if pawns&posMask != 0 {
		return PawnType
	} else if rooks&posMask != 0 {
		return RookType
	} else if knights&posMask != 0 {
		return KnightType
	} else if bishops&posMask != 0 {
		return BishopType
	} else if queens&posMask != 0 {
		return QueenType
	} else if kings&posMask != 0 {
		return KingType
	} else {
		return NoType
	}
}

// Applies given move
// Assumes it is valid and legal
func (b *Board) CoordinateMove(from, to int, promotion Promotion) {
	pieceType := b.pieceType(from)
	capture := PieceTypeToCapture(b.pieceType(to))

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

	noisyMove := (m.PieceType() == PawnType) || (m.Capture() != NoCapture)
	if noisyMove {
		b.noisyMoves = append(b.noisyMoves, b.HalfMoves-1)
	}

	if b.Turn == WhiteTurn {
		// Move piece
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

		// Capture piece
		switch m.Capture() {
		case PawnCapture:
			b.blackPawns ^= toMask
		case RookCapture:
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

		// Capture piece
		switch m.Capture() {
		case PawnCapture:
			b.whitePawns ^= toMask
		case RookCapture:
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

// Returns an array of ranks represented by strings
// Assumes the bitboards are in a valid state
func (b *Board) RankStrings() [8]string {
	var r [8]string
	var rs [8]byte

	var m uint64 = uint64(1) << 63
	for i := range 64 {
		switch {
		case b.whitePawns&m != 0:
			rs[i%8] = 'P'
		case b.whiteRooks&m != 0:
			rs[i%8] = 'R'
		case b.whiteKnights&m != 0:
			rs[i%8] = 'N'
		case b.whiteBishops&m != 0:
			rs[i%8] = 'B'
		case b.whiteQueens&m != 0:
			rs[i%8] = 'Q'
		case b.whiteKings&m != 0:
			rs[i%8] = 'K'
		case b.blackPawns&m != 0:
			rs[i%8] = 'p'
		case b.blackRooks&m != 0:
			rs[i%8] = 'r'
		case b.blackKnights&m != 0:
			rs[i%8] = 'n'
		case b.blackBishops&m != 0:
			rs[i%8] = 'b'
		case b.blackQueens&m != 0:
			rs[i%8] = 'q'
		case b.blackKings&m != 0:
			rs[i%8] = 'k'
		default:
			rs[i%8] = ' '
		}

		if i%8 == 7 {
			r[i/8] = string(rs[:])
		}

		m >>= 1
	}

	return r
}

// Returns a string representation of the board
// Assumes the bitboards are in a valid state
func (b *Board) String() string {
	var s strings.Builder

	var m uint64 = uint64(1) << 63
	for i := range 64 {
		switch {
		case b.whitePawns&m != 0:
			s.WriteByte('P')
		case b.whiteRooks&m != 0:
			s.WriteByte('R')
		case b.whiteKnights&m != 0:
			s.WriteByte('N')
		case b.whiteBishops&m != 0:
			s.WriteByte('B')
		case b.whiteQueens&m != 0:
			s.WriteByte('Q')
		case b.whiteKings&m != 0:
			s.WriteByte('K')
		case b.blackPawns&m != 0:
			s.WriteByte('p')
		case b.blackRooks&m != 0:
			s.WriteByte('r')
		case b.blackKnights&m != 0:
			s.WriteByte('n')
		case b.blackBishops&m != 0:
			s.WriteByte('b')
		case b.blackQueens&m != 0:
			s.WriteByte('q')
		case b.blackKings&m != 0:
			s.WriteByte('k')
		default:
			s.WriteByte(' ')
		}

		if i%8 == 7 {
			s.WriteByte('\n')
		}

		m >>= 1
	}

	return s.String()
}
