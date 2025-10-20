package chess_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zakkbob/chess"
)

func perft(b *chess.Board, depth int) (int, int, int) {
	if depth == 0 {
		captures, enPassant := 0, 0
		if len(b.Moves) != 0 && b.Moves[len(b.Moves)-1].Capture() != chess.NoCapture {
			captures = 1
		}
		if len(b.Moves) != 0 && b.Moves[len(b.Moves)-1].EnPassant() {
			enPassant = 1
		}
		return 1, captures, enPassant
	}

	counter := 0
	captureCounter := 0
	enPassantCounter := 0

	ms := b.LegalMoves()
	for _, m := range ms {
		b.Move(m)
		nodes, captures, enPassant := perft(b, depth-1)
		counter += nodes
		captureCounter += captures
		enPassantCounter += enPassant
		b.Unmove()
	}
	return counter, captureCounter, enPassantCounter
}

func TestPerft(t *testing.T) {
	tests := []struct {
		Name   string
		Board  chess.Board
		Depths []int
	}{
		{
			Name: "Initial Position",
			Depths: []int{
				20,
				400,
				8902,
				197281,
				4865609,
				//119060324,
			},
			Board: chess.NewBoard(),
		},
		{
			Name: "Kiwipete",
			Depths: []int{
				48,
				2039,
				97862,
				4085603,
				//193690690,
			},
			Board: chess.BoardFromRanks(
				[8]string{
					"r   k  r",
					"p ppqpb ",
					"bn  pnp ",
					"   PN   ",
					" p  P   ",
					"  N  Q p",
					"PPPBBPPP",
					"R   K  R",
				},
				chess.WhiteTurn,
				chess.AllCastleRights,
			),
		},
		{
			Name: "Position 3",
			Depths: []int{
				14,
				191,
				2812,
				43238,
				674624,
				11030083,
			},
			Board: chess.BoardFromRanks(
				[8]string{
					"        ",
					"  p     ",
					"   p    ",
					"KP     r",
					" R   p k",
					"        ",
					"    P P ",
					"        ",
				},
				chess.WhiteTurn,
				chess.NoCastleRights,
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			for i := range len(tt.Depths) {
				t.Run(fmt.Sprintf("Depth %d", i+1), func(t *testing.T) {
					got, captures, enPassants := perft(&tt.Board, i+1)

					assert.Equal(t, tt.Depths[i], got, "Incorrect perft result at depth %d", i+1)
					t.Log("Captures", captures)
					t.Log("En Passants", enPassants)
				})
			}
		})
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
		moves = append(moves, chess.NewMove(start, end, pieceType, chess.NoPromotion, chess.NoCapture, false, chess.AllCastleRights, chess.NoCastle))
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
				b := chess.BoardFromRanks(ranks, chess.WhiteTurn, chess.AllCastleRights)
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
				b := chess.BoardFromRanks(ranks, chess.BlackTurn, chess.AllCastleRights)
				ms := b.LegalMoves()
				if !assert.ElementsMatch(t, expected, ms) {
					t.Logf("Expected \n%s\n", moveDiagram(expected, chess.BlackTurn))
					t.Logf("Got \n%s\n", moveDiagram(ms, chess.BlackTurn))
				}
			})
		})
	}
}
