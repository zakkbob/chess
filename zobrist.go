package chess

//go:generate go run zobrist_generate.go

var (
	blackToMoveZobrist    uint64
	castlingRightsZobrist [16]uint64
	enPassantZobrist      [8]uint64

	whitePawnZobrist   [64]uint64
	whiteRookZobrist   [64]uint64
	whiteBishopZobrist [64]uint64
	whiteKnightZobrist [64]uint64
	whiteQueenZobrist  [64]uint64
	whiteKingZobrist   [64]uint64

	blackPawnZobrist   [64]uint64
	blackRookZobrist   [64]uint64
	blackBishopZobrist [64]uint64
	blackKnightZobrist [64]uint64
	blackQueenZobrist  [64]uint64
	blackKingZobrist   [64]uint64
)
