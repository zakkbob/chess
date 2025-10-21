package chess

import "testing"

func BenchmarkSearch(b *testing.B) {
	bd := NewBoard()
	for b.Loop() {
		mv := Search(bd, 4)
		b.Log(mv)
	}
}
