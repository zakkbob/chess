package chess_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zakkbob/chess"
)

func perft(b *chess.Board, depth int) int {
	if depth == 0 {
		return 1
	}

	counter := 0

	ms := b.LegalMoves()
	for _, m := range ms {
		b.Move(m)
		counter += perft(b, depth-1)
		b.Unmove()
	}
	return counter
}

func TestPerft(t *testing.T) {
	perfts := []int{
		20,
		400,
		8902,
		197281,
		4865609,
		//	119060324,
	}

	for i, expected := range perfts {
		b := chess.NewBoard()

		got := perft(&b, i+1)

		assert.Equal(t, expected, got, "Incorrect perft result at depth %d", i+1)
	}
}

// generates move diagram for single piece
func moveDiagram(ms []chess.Move, t chess.Turn) string {
	var rs [64]rune

	for i := range 64 {
		rs[i] = '.'
	}

	for _, m := range ms {
		rs[63-m.To()] = '*'
		rs[63-m.From()] = m.PieceType().Symbol(t)
	}

	var b strings.Builder
	b.WriteString("+--------+\n")
	for i, r := range rs {
		if i%8 == 0 {
			b.WriteRune('|')
		}
		b.WriteRune(r)
		if i%8 == 7 {
			b.WriteRune('|')
			b.WriteRune('\n')
		}
	}
	b.WriteString("+--------+\n")

	return b.String()
}

func getRanksAndMoves(ranks [8]string) ([8]string, []chess.Move) {
	var start int
	var ends []int
	var moves []chess.Move
	var pieceType chess.PieceType
	for i, rank := range ranks {
		for j, cell := range rank {
			switch cell {
			case '*':
				ends = append(ends, 63-chess.Index(i, j))
			case 'P', 'p', 'R', 'r', 'N', 'n', 'B', 'b', 'Q', 'q', 'K', 'k':
				pieceType = chess.PieceTypeFromRune(cell)
				start = 63 - chess.Index(i, j)
			}
		}
	}

	for _, end := range ends {
		moves = append(moves, chess.NewMove(start, end, pieceType, chess.NoPromotion, chess.NoCapture, false, chess.NoCastle))
	}

	return ranks, moves
}

func TestSymmetricOnePieceMoveGen(t *testing.T) {
	tests := []struct {
		Ranks [8]string
	}{
		{
			[8]string{
				"  * *   ",
				" *   *  ",
				"   N    ",
				" *   *  ",
				"  * *   ",
				"        ",
				"        ",
				"        ",
			},
		},
		{
			[8]string{
				"N       ",
				"  *     ",
				" *      ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
			},
		},
		{
			[8]string{
				" N      ",
				"   *    ",
				"* *     ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
			},
		},
		{
			[8]string{
				"   *    ",
				" N      ",
				"   *    ",
				"* *     ",
				"        ",
				"        ",
				"        ",
				"        ",
			},
		},
		{
			[8]string{
				"       N",
				"     *  ",
				"      * ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
			},
		},
		{
			[8]string{
				"    *   ",
				"      N ",
				"    *   ",
				"     * *",
				"        ",
				"        ",
				"        ",
				"        ",
			},
		},
		{
			[8]string{
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
				"      * ",
				"     *  ",
				"       N",
			},
		},
		{
			[8]string{
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
				" *      ",
				"  *     ",
				"N       ",
			},
		},
		{
			[8]string{
				"   *    ",
				"   *    ",
				"   *    ",
				"***R****",
				"   *    ",
				"   *    ",
				"   *    ",
				"   *    ",
			},
		},
		{
			[8]string{
				"R*******",
				"*       ",
				"*       ",
				"*       ",
				"*       ",
				"*       ",
				"*       ",
				"*       ",
			},
		},
		{
			[8]string{
				"       *",
				"*******R",
				"       *",
				"       *",
				"       *",
				"       *",
				"       *",
				"       *",
			},
		},
		{
			[8]string{
				"*       ",
				"*       ",
				"*       ",
				"*       ",
				"*       ",
				"*       ",
				"*       ",
				"R*******",
			},
		},
		{
			[8]string{
				" *     *",
				"  *   * ",
				"   * *  ",
				"    B   ",
				"   * *  ",
				"  *   * ",
				" *     *",
				"*       ",
			},
		},
		{
			[8]string{
				"       *",
				"      * ",
				"     *  ",
				"    *   ",
				"   *    ",
				"  *     ",
				" *      ",
				"B       ",
			},
		},
		{
			[8]string{
				"*       ",
				" *      ",
				"  *     ",
				"   *    ",
				"    *   ",
				"     *  ",
				"      * ",
				"       B",
			},
		},
		{
			[8]string{
				"   *   *",
				"*  *  * ",
				" * * *  ",
				"  ***   ",
				"***Q****",
				"  ***   ",
				" * * *  ",
				"*  *  * ",
			},
		},
		{
			[8]string{
				" *      ",
				" *     *",
				" *    * ",
				" *   *  ",
				" *  *   ",
				" * *    ",
				"***     ",
				"*Q******",
			},
		},
		{
			[8]string{
				"        ",
				"        ",
				"        ",
				"  ***   ",
				"  *K*   ",
				"  ***   ",
				"        ",
				"        ",
			},
		},
		{
			[8]string{
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
				"**      ",
				"K*      ",
			},
		},
		{
			[8]string{
				"K*      ",
				"**      ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
			},
		},
		{
			[8]string{
				"      *K",
				"      **",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
			},
		},
		{
			[8]string{
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
				"        ",
				"      **",
				"      *K",
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("Test %d", i), func(t *testing.T) {
			t.Run("White", func(t *testing.T) {
				ranks, expected := getRanksAndMoves(tt.Ranks)
				b := chess.BoardFromRanks(ranks, chess.WhiteTurn)
				ms := b.LegalMoves()
				if !assert.ElementsMatch(t, expected, ms) {
					t.Logf("Expected \n%s\n", moveDiagram(expected, chess.WhiteTurn))
					t.Logf("Got \n%s\n", moveDiagram(ms, chess.WhiteTurn))
				}
			})
			//replace with black piece
			found := false
			for i, rank := range tt.Ranks {
				for j, cell := range rank {
					if cell == ' ' || cell == '*' {
						continue
					}
					symbol := chess.PieceTypeFromRune(cell).Symbol(chess.BlackTurn)
					bRank := []rune(rank)
					bRank[j] = symbol
					tt.Ranks[i] = string(bRank)
					found = true
					break
				}
				if found {
					break
				}
			}
			t.Run("Black", func(t *testing.T) {
				ranks, expected := getRanksAndMoves(tt.Ranks)
				b := chess.BoardFromRanks(ranks, chess.BlackTurn)
				ms := b.LegalMoves()
				if !assert.ElementsMatch(t, expected, ms) {
					t.Logf("Expected \n%s\n", moveDiagram(expected, chess.BlackTurn))
					t.Logf("Got \n%s\n", moveDiagram(ms, chess.BlackTurn))
				}
			})
		})
	}
}

func TestPinnedPiece(t *testing.T) {
	b := chess.BoardFromRanks([8]string{
		"        ",
		"        ",
		"q       ",
		" R      ",
		"   q    ",
		" n      ",
		" RqR    ",
		"K       ",
	}, chess.WhiteTurn)

	t.Log(b.LegalMoves())
	t.Fail()
}
