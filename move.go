package chess

import (
	"errors"
	"strings"
)

var (
	ErrInvalidAlgebraicNotation = errors.New("Invalid algebraic notation")
)

func IndexFromAlgebraic(a string) (int, error) {
	if len(a) != 2 {
		return 0, ErrInvalidAlgebraicNotation
	}
	var i int
	switch a[0] {
	case 'a', 'A':
		i = 7
	case 'b', 'B':
		i = 6
	case 'c', 'C':
		i = 5
	case 'd', 'D':
		i = 4
	case 'e', 'E':
		i = 3
	case 'f', 'F':
		i = 2
	case 'g', 'G':
		i = 1
	case 'h', 'H':
		i = 0
	default:
		return 0, ErrInvalidAlgebraicNotation
	}
	switch a[1] {
	case '1':
		i += 0
	case '2':
		i += 8
	case '3':
		i += 16
	case '4':
		i += 24
	case '5':
		i += 32
	case '6':
		i += 40
	case '7':
		i += 48
	case '8':
		i += 56
	default:
		return 0, ErrInvalidAlgebraicNotation
	}

	return i, nil
}

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
	var b strings.Builder

	files := []rune{'h', 'g', 'f', 'e', 'd', 'c', 'b', 'a'}
	ranks := []rune{'1', '2', '3', '4', '5', '6', '7', '8'}

	b.WriteRune(files[m.FromFile()])
	b.WriteRune(ranks[m.FromRank()])
	b.WriteRune(files[m.ToFile()])
	b.WriteRune(ranks[m.ToRank()])

	if m.Promotion() != NoPromotion {
		b.WriteRune(m.Promotion().Symbol())
	}

	return b.String()
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

func (m Move) CastleRights() CastleRights {
	return CastleRights(uint32(m) & CastleRightsMask)
}

func (m Move) Castle() Castle {
	return Castle(uint32(m) & CastleMask)
}

func NewMove(from, to int, pieceType PieceType, promotion Promotion, capture Capture, enPassant bool, castleRights CastleRights, castle Castle) Move {
	m := Move(uint32(from<<23) | uint32(to<<17) | uint32(pieceType) | uint32(promotion) | uint32(capture) | uint32(castleRights) | uint32(castle))
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

func PieceTypeFromRune(r rune) PieceType {
	switch r {
	case 'P', 'p':
		return PawnType
	case 'R', 'r':
		return RookType
	case 'N', 'n':
		return KnightType
	case 'B', 'b':
		return BishopType
	case 'Q', 'q':
		return QueenType
	case 'K', 'k':
		return KingType
	default:
		return NoType
	}
}

func (p PieceType) Symbol(t Turn) rune {
	if t == WhiteTurn {
		switch p {
		case PawnType:
			return 'P'
		case RookType:
			return 'R'
		case KnightType:
			return 'N'
		case BishopType:
			return 'B'
		case QueenType:
			return 'Q'
		case KingType:
			return 'K'
		case NoType:
			return ' '
		default:
			panic("unable to get PieceType symbol - invalid PieceType")
		}
	} else {
		switch p {
		case PawnType:
			return 'p'
		case RookType:
			return 'r'
		case KnightType:
			return 'n'
		case BishopType:
			return 'b'
		case QueenType:
			return 'q'
		case KingType:
			return 'k'
		case NoType:
			return ' '
		default:
			panic("unable to get PieceType symbol - invalid PieceType")
		}
	}
}

func (p PieceType) String() string {
	switch p {
	case PawnType:
		return "pawn"
	case RookType:
		return "rook"
	case KnightType:
		return "knight"
	case BishopType:
		return "bishop"
	case QueenType:
		return "queen"
	case KingType:
		return "king"
	case NoType:
		return "none"
	default:
		panic("cannot convert PieceType to string - unknown piece type")
	}
}

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

func PromotionFromSymbol(r rune) Promotion {
	switch r {
	case ' ':
		return NoPromotion
	case 'r', 'R':
		return RookPromotion
	case 'n', 'N':
		return KnightPromotion
	case 'b', 'B':
		return BishopPromotion
	case 'q', 'Q':
		return QueenPromotion
	default:
		panic("invalid promotion symbol")
	}
}

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

func (p Promotion) Symbol() rune {
	switch p {
	case NoPromotion:
		return ' '
	case RookPromotion:
		return 'r'
	case KnightPromotion:
		return 'n'
	case BishopPromotion:
		return 'b'
	case QueenPromotion:
		return 'q'
	default:
		panic("invalid promotion")
	}
}

type Capture uint32

const (
	// Capture              0b00000000000000000011100000000000
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

type CastleRights uint32

const (
	// CastleRights                 0b00000000000000000000001111000000
	NoCastleRights   CastleRights = 0b00000000000000000000000000000000
	AllCastleRights  CastleRights = 0b00000000000000000000001111000000
	WhiteKingCastle  CastleRights = 0b00000000000000000000001000000000
	WhiteQueenCastle CastleRights = 0b00000000000000000000000100000000
	BlackKingCastle  CastleRights = 0b00000000000000000000000010000000
	BlackQueenCastle CastleRights = 0b00000000000000000000000001000000
)

func NewCastleRights(whiteKing, whiteQueen, blackKing, blackQueen bool) CastleRights {
	cr := NoCastleRights
	if whiteKing {
		cr |= WhiteKingCastle
	}
	if whiteQueen {
		cr |= WhiteQueenCastle
	}
	if blackKing {
		cr |= BlackKingCastle
	}
	if blackQueen {
		cr |= BlackQueenCastle
	}
	return cr
}

func CastleRightsFromString(s string) CastleRights {
	if s == "-" {
		return NoCastleRights
	}

	cr := NoCastleRights

	for _, r := range s {
		switch r {
		case 'K':
			cr |= WhiteKingCastle
		case 'Q':
			cr |= WhiteQueenCastle
		case 'k':
			cr |= BlackKingCastle
		case 'q':
			cr |= BlackQueenCastle
		default:
			panic("unknown side")
		}
	}

	return cr
}

func (cr CastleRights) CanWhiteKing() bool {
	return cr&WhiteKingCastle != 0
}

func (cr CastleRights) CanWhiteQueen() bool {
	return cr&WhiteQueenCastle != 0
}

func (cr CastleRights) CanBlackKing() bool {
	return cr&BlackKingCastle != 0
}

func (cr CastleRights) CanBlackQueen() bool {
	return cr&BlackQueenCastle != 0
}

func (cr *CastleRights) LoseWhiteKing() {
	*cr &= ^WhiteKingCastle
}

func (cr *CastleRights) LoseWhiteQueen() {
	*cr &= ^WhiteQueenCastle
}

func (cr *CastleRights) LoseBlackKing() {
	*cr &= ^BlackKingCastle
}

func (cr *CastleRights) LoseBlackQueen() {
	*cr &= ^BlackQueenCastle
}

type Castle uint32

const (
	// Castle            0b00000000000000000000000000110000
	NoCastle    Castle = 0b00000000000000000000000000000000
	KingCastle  Castle = 0b00000000000000000000000000100000
	QueenCastle Castle = 0b00000000000000000000000000110000
)
