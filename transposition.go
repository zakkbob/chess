package chess

type entryType int

const (
	ExactEntry entryType = iota
	LowerBoundEntry
	UpperBoundEntry
)

type Transposition struct {
	Key      uint64
	BestMove Move
	Depth    int
	Score    int
	Type     entryType
}

type TranspositionTable struct {
	entries []Transposition
	mask    uint64
}

// Creates a transposition table with 2^exp entries
func NewTranspositionTable(exp int) *TranspositionTable {
	length := 1 << exp
	mask := uint64(length - 1)
	return &TranspositionTable{
		entries: make([]Transposition, length),
		mask:    mask,
	}
}

func (tt *TranspositionTable) Get(key uint64) (Transposition, bool) {
	i := (key & tt.mask)
	t := tt.entries[i]
	return t, t.Key == key
}

// Always overwrites existing entry
func (tt *TranspositionTable) Save(t Transposition) {
	i := (t.Key & tt.mask)
	tt.entries[i] = t
}
