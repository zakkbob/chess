package chess

type entryType int

const (
	exact entryType = iota
	lowerBound
	upperBound
)

type Transposition struct {
	Key      uint64
	BestMove Move
	Depth    int
	Score    int
	Type     entryType
}

type TranspositionTable struct {
	Entries []Transposition
	Mask uint64
}

func NewTranspositionTable
