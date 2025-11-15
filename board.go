package chess

import (
	"errors"

	"math/bits"
	"slices"
	"strconv"
	"strings"
)

var (
	ErrInvalidFEN = errors.New("Invalid FEN")
)

func Index(rank, file int) int {
	return rank*8 + file
}

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

	Turn          Turn
	HalfMoves     int
	CastleRights  CastleRights
	Moves         []Move
	noisyMoves    []int
	CanEnPassant  bool
	EnPassantFile int // if last move was a double push, holds file
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

		Turn:         WhiteTurn,
		HalfMoves:    0,
		Moves:        []Move{},
		noisyMoves:   []int{},
		CastleRights: AllCastleRights,
	}
}

func BoardFromFEN(fen string) (Board, error) {
	parts := strings.Split(fen, " ")
	if len(parts) != 6 && len(parts) != 4 {
		return Board{}, ErrInvalidFEN
	}

	b := Board{
		Moves:      []Move{},
		noisyMoves: []int{},
	}

	pieces := parts[0]

	i := 0
	for _, symbol := range pieces {
		posMask := uint64(1) << (63 - i)
		switch symbol {
		case '/':
			continue
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
		case '1':

		case '2':
			i += 1
		case '3':
			i += 2
		case '4':
			i += 3
		case '5':
			i += 4
		case '6':
			i += 5
		case '7':
			i += 6
		case '8':
			i += 7
		default:
			return Board{}, ErrInvalidFEN
		}
		i++
	}

	active := parts[1]
	switch active {
	case "w", "W":
		b.Turn = WhiteTurn
	case "b", "B":
		b.Turn = BlackTurn
	default:
		return Board{}, ErrInvalidFEN
	}

	b.CastleRights = CastleRightsFromString(parts[2])

	enPassantTarget := parts[3]
	if enPassantTarget != "-" {
		i, err := IndexFromAlgebraic(enPassantTarget)
		if err != nil {
			return Board{}, ErrInvalidFEN
		}
		b.CanEnPassant = true
		b.EnPassantFile = i % 8
	}

	if len(parts) == 6 {
		fullMoveClock, err := strconv.Atoi(parts[5])
		if err != nil {
			return Board{}, ErrInvalidFEN
		}
		if b.Turn == WhiteTurn {
			b.HalfMoves = (fullMoveClock - 1) * 2
		} else {
			b.HalfMoves = (fullMoveClock-1)*2 + 1
		}

		halfMoveClock, err := strconv.Atoi(parts[4])
		if err != nil {
			return Board{}, ErrInvalidFEN
		}
		b.noisyMoves = append(b.noisyMoves, b.HalfMoves-halfMoveClock)
	}

	return b, nil
}

func BoardFromRanks(rs [8]string, turn Turn, castleRights CastleRights) Board {
	b := Board{
		Turn:         turn,
		HalfMoves:    0,
		Moves:        []Move{},
		noisyMoves:   []int{},
		CastleRights: castleRights,
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

func (b *Board) Zobrist() uint64 {
	var z uint64

	applyBitboard := func(bb uint64, zobristVals [64]uint64) {
		for bb != 0 {
			i := bits.LeadingZeros64(bb)
			z ^= zobristVals[i]
			bb &= bb - 1
		}
	}

	if b.Turn == BlackTurn {
		z ^= blackTurnZobrist
	}

	if b.CanEnPassant {
		z ^= enPassantZobrist[b.EnPassantFile]
	}

	z ^= castlingRightsZobrist[b.CastleRights.Uint64()]

	applyBitboard(b.whitePawns, whitePawnZobrist)
	applyBitboard(b.whiteRooks, whiteRookZobrist)
	applyBitboard(b.whiteBishops, whiteBishopZobrist)
	applyBitboard(b.whiteKnights, whiteKnightZobrist)
	applyBitboard(b.whiteQueens, whiteQueenZobrist)
	applyBitboard(b.whiteKings, whiteKingZobrist)

	applyBitboard(b.blackPawns, blackPawnZobrist)
	applyBitboard(b.blackRooks, blackRookZobrist)
	applyBitboard(b.blackBishops, blackBishopZobrist)
	applyBitboard(b.blackKnights, blackKnightZobrist)
	applyBitboard(b.blackQueens, blackQueenZobrist)
	applyBitboard(b.blackKings, blackKingZobrist)

	return z
}

func (b *Board) Copy() Board {
	c := *b
	c.Moves = slices.Clone(b.Moves)
	c.noisyMoves = slices.Clone(b.noisyMoves)
	return c
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
		if i%8 == 0 {
			s.WriteString(strconv.FormatInt(int64(8-i/8), 10))
			s.WriteString("|")
		}

		switch {
		case b.whitePawns&m != 0:
			s.WriteRune('♙')
		case b.whiteRooks&m != 0:
			s.WriteRune('♖')
		case b.whiteKnights&m != 0:
			s.WriteRune('♘')
		case b.whiteBishops&m != 0:
			s.WriteRune('♗')
		case b.whiteQueens&m != 0:
			s.WriteRune('♕')
		case b.whiteKings&m != 0:
			s.WriteRune('♔')
		case b.blackPawns&m != 0:
			s.WriteRune('♟')
		case b.blackRooks&m != 0:
			s.WriteRune('♜')
		case b.blackKnights&m != 0:
			s.WriteRune('♞')
		case b.blackBishops&m != 0:
			s.WriteRune('♝')
		case b.blackQueens&m != 0:
			s.WriteRune('♛')
		case b.blackKings&m != 0:
			s.WriteRune('♚')
		default:
			s.WriteByte('.')
		}
		s.WriteByte(' ')

		if i%8 == 7 && i < 63 {
			s.WriteByte('\n')
		}

		m >>= 1
	}

	s.WriteByte('\n')

	s.WriteString("  ---------------\n")
	s.WriteString("  a b c d e f g h")

	return s.String()
}
