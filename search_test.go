package chess

import "testing"

func BenchmarkSearch(b *testing.B) {
	bd := NewBoard()
	for b.Loop() {
		mv := Search(bd, 5)
		b.Log(mv)
	}
}
