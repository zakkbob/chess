package chess

import "strconv"

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
type Move uint32

func (m Move) String() string {
	s := strconv.Itoa(int(m.From())) + " to " + strconv.Itoa(int(m.To()))

	if m.Capture() != NoCapture {
		s += ", captures " + m.Capture().String()
	}

	if m.Promotion() != NoPromotion {
		s += ", promotes " + m.Promotion().String()
	}

	if m.EnPassant() {
		s += ", en passants"
	}

	return s
}

func (m Move) IsNoisy() bool {
	return (m.PieceType() == PawnType) || (m.Capture() != NoCapture)
}

func (m Move) IsDoublePush() bool {
	return (m.PieceType() == PawnType) && ((m.FromRank() == 1 && m.ToRank() == 3) || (m.FromRank() == 6 && m.ToRank() == 4))
}

func (m Move) PieceType() PieceType {
	return PieceType(uint32(m) & PieceTypeMask)
}

func (m Move) From() uint32 {
	return (uint32(m) & FromMask) >> 23
}

func (m Move) FromRank() uint32 {
	return m.From() / 8
}

func (m Move) FromFile() uint32 {
	return m.From() % 8
}

func (m Move) To() uint32 {
	return (uint32(m) & ToMask) >> 17
}

func (m Move) ToRank() uint32 {
	return m.To() / 8
}

func (m Move) ToFile() uint32 {
	return m.To() % 8
}

func (m Move) Promotion() Promotion {
	return Promotion(uint32(m) & PromotionMask)
}

func (m Move) Capture() Capture {
	return Capture(uint32(m) & CaptureMask)
}

func (m Move) EnPassant() bool {
	return uint32(m)&EnPassantMask != 0
}

func (m Move) Castle() Castle {
	return Castle(uint32(m) & CastleMask)
}

func NewMove(from, to int, pieceType PieceType, promotion Promotion, capture Capture, enPassant bool, castle Castle) Move {
	m := Move(uint32(from<<23) | uint32(to<<17) | uint32(pieceType) | uint32(promotion) | uint32(capture) | uint32(castle))
	if enPassant {
		m |= Move(EnPassantMask)
	}
	return m
}

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

type PieceType uint32

const (
	// Piece Type       0b11100000000000000000000000000000
	PawnType   PieceType = 0b00000000000000000000000000000000
	RookType   PieceType = 0b00100000000000000000000000000000
	KnightType PieceType = 0b01000000000000000000000000000000
	BishopType PieceType = 0b01100000000000000000000000000000
	QueenType  PieceType = 0b10000000000000000000000000000000
	KingType   PieceType = 0b10100000000000000000000000000000
	NoType     PieceType = 0b11000000000000000000000000000000
)

func (p PieceType) ToCapture() Capture {
	switch p {
	case PawnType:
		return PawnCapture
	case RookType:
		return RookCapture
	case KnightType:
		return KnightCapture
	case BishopType:
		return BishopCapture
	case QueenType:
		return QueenCapture
	case KingType:
		panic("unable to convert KingType to capture - cannot capture king")
	case NoType:
		return NoCapture
	default:
		panic("unable to convert PieceType to capture - invalid piece type")
	}
}

type Promotion uint32

const (
	// Promotion             0b00000000000000011100000000000000
	NoPromotion     Promotion = 0b00000000000000000000000000000000
	RookPromotion   Promotion = 0b00000000000000010000000000000000
	KnightPromotion Promotion = 0b00000000000000010100000000000000
	BishopPromotion Promotion = 0b00000000000000011000000000000000
	QueenPromotion  Promotion = 0b00000000000000011100000000000000
)

func (p Promotion) String() string {
	switch p {
	case NoPromotion:
		return "nothing"
	case RookPromotion:
		return "rook"
	case KnightPromotion:
		return "knight"
	case BishopPromotion:
		return "bishop"
	case QueenPromotion:
		return "queen"
	default:
		panic("invalid promotion")
	}
}

type Capture uint32

const (
	// Capture             0b00000000000000000011100000000000
	NoCapture     Capture = 0b00000000000000000000000000000000
	PawnCapture   Capture = 0b00000000000000000000100000000000
	RookCapture   Capture = 0b00000000000000000001000000000000
	KnightCapture Capture = 0b00000000000000000001100000000000
	BishopCapture Capture = 0b00000000000000000010000000000000
	QueenCapture  Capture = 0b00000000000000000010100000000000
)

func (c Capture) String() string {
	switch c {
	case NoCapture:
		return "nothing"
	case PawnCapture:
		return "pawn"
	case RookCapture:
		return "rook"
	case KnightCapture:
		return "knight"
	case BishopCapture:
		return "bishop"
	case QueenCapture:
		return "queen"
	default:
		panic("invalid capture")
	}
}

type Castle uint32

const (
	// Castle            0b00000000000000000000000000110000
	NoCastle    Castle = 0b00000000000000000000000000000000
	KingCastle  Castle = 0b00000000000000000000000000100000
	QueenCastle Castle = 0b00000000000000000000000000110000
)
