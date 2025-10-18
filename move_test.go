package chess_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zakkbob/chess"
)

func TestIsDoublePush(t *testing.T) {
	for file := range 8 {
		t.Run(fmt.Sprintf("White pawn single push (home rank, file %d)", file), func(t *testing.T) {
			m := chess.NewMove(chess.Index(1, file), chess.Index(2, file), chess.PawnType, chess.NoPromotion, chess.NoCapture, false, chess.NoCastle)
			assert.Equal(t, false, m.IsDoublePush())
		})

		t.Run(fmt.Sprintf("White pawn double push (file %d)", file), func(t *testing.T) {
			m := chess.NewMove(chess.Index(1, file), chess.Index(3, file), chess.PawnType, chess.NoPromotion, chess.NoCapture, false, chess.NoCastle)
			assert.Equal(t, true, m.IsDoublePush())
		})

		t.Run(fmt.Sprintf("Black pawn single push (home rank, file %d)", file), func(t *testing.T) {
			m := chess.NewMove(chess.Index(6, file), chess.Index(5, file), chess.PawnType, chess.NoPromotion, chess.NoCapture, false, chess.NoCastle)
			assert.Equal(t, false, m.IsDoublePush())
		})

		t.Run(fmt.Sprintf("Black pawn double push (file %d)", file), func(t *testing.T) {
			m := chess.NewMove(chess.Index(6, file), chess.Index(4, file), chess.PawnType, chess.NoPromotion, chess.NoCapture, false, chess.NoCastle)
			assert.Equal(t, true, m.IsDoublePush())
		})
	}

	t.Run("Non-pawn move", func(t *testing.T) {
		m := chess.NewMove(1, 10, chess.BishopType, chess.NoPromotion, chess.NoCapture, false, chess.NoCastle)
		assert.Equal(t, false, m.IsDoublePush())
	})
}
