package chess

import (
	"strings"
)

// --- Move Representation ---
// Bits Overview (inclusive)
// 0-2   - Piece type
// 3-8   - From
// 9-14  - To
// 15-17 - Promotion
// 18-20 - Capture
// 21    - En passant
// 22-25 - Castling rights (before move)
// 26-27 - Castle side (if the move is a castle)
//
// Piece type
// 000 - Pawn
// 001 - Rook
// 010 - Knight
// 011 - Bishop
// 100 - Queen
// 101 - King
//
// Promotion
// 000 - None
// 100 - Rook
// 101 - Knight
// 110 - Bishop
// 111 - Queen
//
// Capture
// 000 - None
// 001 - Pawn
// 010 - Rook
// 011 - Knight
// 100 - Bishop
// 101 - Queen
//
// Castling rights
// 1xxx - White kingside
// x1xx - White queenside
// xx1x - Black kingside
// xxx1 - Black queenside
//
// Castle side
// 00 - No castle
// 10 - King side
// 11 - Queen side
type Move = uint32

const (
	PieceTypeMask    uint32 = 0b11100000000000000000000000000000
	FromMask         uint32 = 0b00011111100000000000000000000000
	ToMask           uint32 = 0b00000000011111100000000000000000
	PromotionMask    uint32 = 0b00000000000000011100000000000000
	CaptureMask      uint32 = 0b00000000000000000011100000000000
	EnPassantMask    uint32 = 0b00000000000000000000010000000000
	CastleRightsMask uint32 = 0b00000000000000000000001111000000
	CastleMask       uint32 = 0b00000000000000000000000000110000
)

const (
	// Piece Type       0b11100000000000000000000000000000
	PawnType   uint32 = 0b00000000000000000000000000000000
	RookType   uint32 = 0b00100000000000000000000000000000
	KnightType uint32 = 0b01000000000000000000000000000000
	BishopType uint32 = 0b01100000000000000000000000000000
	QueenType  uint32 = 0b10000000000000000000000000000000
	KingType   uint32 = 0b10100000000000000000000000000000
)

const (
	// Promotion             0b00000000000000011100000000000000
	NoPromotion     uint32 = 0b00000000000000000000000000000000
	RookPromotion   uint32 = 0b00000000000000010000000000000000
	KnightPromotion uint32 = 0b00000000000000010100000000000000
	BishopPromotion uint32 = 0b00000000000000011000000000000000
	QueenPromotion  uint32 = 0b00000000000000011100000000000000
)

const (
	// Capture             0b00000000000000000011100000000000
	NoCapture     uint32 = 0b00000000000000000000000000000000
	PawnCapture   uint32 = 0b00000000000000000000100000000000
	RookCapture   uint32 = 0b00000000000000000001000000000000
	KnightCapture uint32 = 0b00000000000000000001100000000000
	BishopCapture uint32 = 0b00000000000000000010000000000000
	QueenCapture  uint32 = 0b00000000000000000010100000000000
)
const (
	// Castle            0b00000000000000000000000000110000
	NoCastle    uint32 = 0b00000000000000000000000000000000
	KingCastle  uint32 = 0b00000000000000000000000000100000
	QueenCastle uint32 = 0b00000000000000000000000000110000
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

	turn Turn
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
	}
}

func BoardFromRanks(rs [8]string, turn Turn) Board {
	b := Board{
		turn: turn,
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

func (b *Board) Move(m Move) {
	var from uint32 = (m & 0b00011111100000000000000000000000) >> 23
	var to uint32 = (m & 0b00000000011111100000000000000000) >> 17

	var fromMask = uint64(1) << from
	var toMask = uint64(1) << to

	var moveMask uint64 = fromMask | toMask

	if b.turn == WhiteTurn {
		// Move piece
		switch m & PieceTypeMask {
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
		switch m & CaptureMask {
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
		switch m & PromotionMask {
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
		switch m & CastleMask {
		case KingCastle:
			b.whiteRooks ^= 0b00000101
		case QueenCastle:
			b.whiteRooks ^= 0b10010000
		}

		// Handle en passant
		if m&EnPassantMask != 0 {
			b.blackPawns ^= toMask >> 8
		}
	} else {
		// Move piece
		switch m & PieceTypeMask {
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
		switch m & CaptureMask {
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
		switch m & PromotionMask {
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
		switch m & CastleMask {
		case KingCastle:
			b.blackRooks ^= 0b00000101 << 56
		case QueenCastle:
			b.blackRooks ^= 0b10010000 << 56
		}

		// Handle en passant
		if m&EnPassantMask != 0 {
			b.whitePawns ^= toMask << 8
		}
	}

}

func (b *Board) Unmove(m Move) {
	var from uint32 = (m & 0b00011111100000000000000000000000) >> 23
	var to uint32 = (m & 0b00000000011111100000000000000000) >> 17

	var fromMask = uint64(1) << from
	var toMask = uint64(1) << to

	var moveMask uint64 = fromMask | toMask

	if b.turn == WhiteTurn {
		// Undo promotion
		switch m & PromotionMask {
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
		switch m & PieceTypeMask {
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
		switch m & CaptureMask {
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
		switch m & CastleMask {
		case KingCastle:
			b.whiteRooks ^= 0b00000101
		case QueenCastle:
			b.whiteRooks ^= 0b10010000
		}

		// Undo en passant
		if m&EnPassantMask != 0 {
			b.blackPawns ^= toMask >> 8
		}
	} else {
		// Undo promotion
		switch m & PromotionMask {
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
		switch m & PieceTypeMask {
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
		switch m & CaptureMask {
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
		switch m & CastleMask {
		case KingCastle:
			b.blackRooks ^= 0b00000101 << 56
		case QueenCastle:
			b.blackRooks ^= 0b10010000 << 56
		}

		// Undo en passant
		if m&EnPassantMask != 0 {
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
