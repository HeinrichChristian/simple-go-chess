package main

import "fmt"

const (
	Pawn = iota
	Knight
	Bishop
	Rook
	Queen
	King
	Empty
)

const (
	White = iota
	Black
)

// piece values indexed by piece constants: Pawn..King, Empty
var pieceValues = [...]int{
	100,   // Pawn
	300,   // Knight
	300,   // Bishop
	500,   // Rook
	900,   // Queen
	10000, // King
	0,     // Empty
}

// initial colors for each square at game start:
// order matches Square enum (A1..H1, A2..H2, ..., A8..H8)
// 0 = White, 1 = Black, 6 = Empty
var initColors = [...]int{
	// Rank 1 (A1..H1) - White pieces
	White, White, White, White, White, White, White, White,
	// Rank 2 (A2..H2) - White pawns
	White, White, White, White, White, White, White, White,
	// Rank 3 (A3..H3) - empty
	Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty,
	// Rank 4 (A4..H4) - empty
	Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty,
	// Rank 5 (A5..H5) - empty
	Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty,
	// Rank 6 (A6..H6) - empty
	Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty,
	// Rank 7 (A7..H7) - Black pawns
	Black, Black, Black, Black, Black, Black, Black, Black,
	// Rank 8 (A8..H8) - Black pieces
	Black, Black, Black, Black, Black, Black, Black, Black,
}

// initial piece types for each square at game start:
// order matches Square enum (A1..H1, A2..H2, ..., A8..H8)
// values are Pawn..King, Empty
var initPieces = [...]int{
	// Rank 1 (A1..H1) - White back rank
	Rook, Knight, Bishop, Queen, King, Bishop, Knight, Rook,
	// Rank 2 (A2..H2) - White pawns
	Pawn, Pawn, Pawn, Pawn, Pawn, Pawn, Pawn, Pawn,
	// Rank 3 (A3..H3) - empty
	Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty,
	// Rank 4 (A4..H4) - empty
	Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty,
	// Rank 5 (A5..H5) - empty
	Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty,
	// Rank 6 (A6..H6) - empty
	Empty, Empty, Empty, Empty, Empty, Empty, Empty, Empty,
	// Rank 7 (A7..H7) - Black pawns
	Pawn, Pawn, Pawn, Pawn, Pawn, Pawn, Pawn, Pawn,
	// Rank 8 (A8..H8) - Black back rank
	Rook, Knight, Bishop, Queen, King, Bishop, Knight, Rook,
}

// Square represents a board square (a1..h8).
type Square int

const (
	A1 Square = iota
	B1
	C1
	D1
	E1
	F1
	G1
	H1

	A2
	B2
	C2
	D2
	E2
	F2
	G2
	H2

	A3
	B3
	C3
	D3
	E3
	F3
	G3
	H3

	A4
	B4
	C4
	D4
	E4
	F4
	G4
	H4

	A5
	B5
	C5
	D5
	E5
	F5
	G5
	H5

	A6
	B6
	C6
	D6
	E6
	F6
	G6
	H6

	A7
	B7
	C7
	D7
	E7
	F7
	G7
	H7

	A8
	B8
	C8
	D8
	E8
	F8
	G8
	H8
)

// file (column) index for each square: 0 = file a, 1 = file b, ... 7 = file h
var squareFile = [...]int{
	// Rank 1 (A1..H1)
	0, 1, 2, 3, 4, 5, 6, 7,
	// Rank 2 (A2..H2)
	0, 1, 2, 3, 4, 5, 6, 7,
	// Rank 3 (A3..H3)
	0, 1, 2, 3, 4, 5, 6, 7,
	// Rank 4 (A4..H4)
	0, 1, 2, 3, 4, 5, 6, 7,
	// Rank 5 (A5..H5)
	0, 1, 2, 3, 4, 5, 6, 7,
	// Rank 6 (A6..H6)
	0, 1, 2, 3, 4, 5, 6, 7,
	// Rank 7 (A7..H7)
	0, 1, 2, 3, 4, 5, 6, 7,
	// Rank 8 (A8..H8)
	0, 1, 2, 3, 4, 5, 6, 7,
}

// rank (row) index for each square: 0 = rank 1, 1 = rank 2, ... 7 = rank 8
var squareRank = [...]int{
	// Rank 1 (A1..H1)
	0, 0, 0, 0, 0, 0, 0, 0,
	// Rank 2 (A2..H2)
	1, 1, 1, 1, 1, 1, 1, 1,
	// Rank 3 (A3..H3)
	2, 2, 2, 2, 2, 2, 2, 2,
	// Rank 4 (A4..H4)
	3, 3, 3, 3, 3, 3, 3, 3,
	// Rank 5 (A5..H5)
	4, 4, 4, 4, 4, 4, 4, 4,
	// Rank 6 (A6..H6)
	5, 5, 5, 5, 5, 5, 5, 5,
	// Rank 7 (A7..H7)
	6, 6, 6, 6, 6, 6, 6, 6,
	// Rank 8 (A8..H8)
	7, 7, 7, 7, 7, 7, 7, 7,
}

// squareScoreTable holds for each color and piece type a precomputed score for every square.
// index 0 = White, 1 = Black; piece index = Pawn..King, Empty
var squareScoreTable [2][7][64]int

// kingOpeningScore and kingEndgameScore per color
var kingOpeningScore [2][64]int
var kingEndgameScore [2][64]int

// flipSquare maps a square to the vertically flipped square (same file, mirrored rank).
// e.g. A1 -> A8, B2 -> B7, ...
var flipSquare [64]int

// current board arrays (mutable game state)
var boardPieces [64]int // piece type on each square (Pawn..King, Empty)
var boardColors [64]int // color on each square (White, Black, Empty)

// global half-move ply counter
var hply int

// initSquareScoreTable fills squareScoreTable and kingOpeningScore/kingEndgameScore
// by combining material value and positional tables defined in eval.go.
// Black tables are computed by flipping the positional tables via flipSquare.
func initSquareScoreTable() {
	// build flip table
	for sq := 0; sq < 64; sq++ {
		file := sq % 8
		rank := sq / 8
		flipSquare[sq] = file + (7-rank)*8
	}

	// initialize move targets for all pieces
	initMoveTargets()

	for sq := 0; sq < 64; sq++ {
		// White: use positional tables as-is
		squareScoreTable[White][Pawn][sq] = pieceValues[Pawn] + PawnScore[sq]
		squareScoreTable[White][Knight][sq] = pieceValues[Knight] + KnightScore[sq]
		squareScoreTable[White][Bishop][sq] = pieceValues[Bishop] + BishopScore[sq]
		squareScoreTable[White][Rook][sq] = pieceValues[Rook] + RookScore[sq]
		squareScoreTable[White][Queen][sq] = pieceValues[Queen] + QueenScore[sq]
		squareScoreTable[White][Empty][sq] = pieceValues[Empty]
		kingOpeningScore[White][sq] = pieceValues[King] + KingScore[sq]
		kingEndgameScore[White][sq] = pieceValues[King] + KingEndgameScore[sq]

		// Black: mirror positional tables via flipSquare
		fs := flipSquare[sq]
		squareScoreTable[Black][Pawn][sq] = pieceValues[Pawn] + PawnScore[fs]
		squareScoreTable[Black][Knight][sq] = pieceValues[Knight] + KnightScore[fs]
		squareScoreTable[Black][Bishop][sq] = pieceValues[Bishop] + BishopScore[fs]
		squareScoreTable[Black][Rook][sq] = pieceValues[Rook] + RookScore[fs]
		squareScoreTable[Black][Queen][sq] = pieceValues[Queen] + QueenScore[fs]
		squareScoreTable[Black][Empty][sq] = pieceValues[Empty]
		kingOpeningScore[Black][sq] = pieceValues[King] + KingScore[fs]
		kingEndgameScore[Black][sq] = pieceValues[King] + KingEndgameScore[fs]
	}
}

// initBoard initializes boardPieces and boardColors from the starting tables
// and resets the half-move ply counter.
func initBoard() {
	for i := 0; i < 64; i++ {
		boardPieces[i] = initPieces[i]
		boardColors[i] = initColors[i]
	}
	hply = 0
}

// printBoard prints a simple ASCII representation of the current board.
// White pieces are uppercase, Black pieces lowercase, empty squares shown as '.'.
func printBoard() {
	whiteChar := map[int]byte{
		Pawn: 'P', Knight: 'N', Bishop: 'B', Rook: 'R', Queen: 'Q', King: 'K', Empty: '.',
	}
	blackChar := map[int]byte{
		Pawn: 'p', Knight: 'n', Bishop: 'b', Rook: 'r', Queen: 'q', King: 'k', Empty: '.',
	}

	for r := 7; r >= 0; r-- {
		fmt.Printf("%d ", r+1)
		for f := 0; f < 8; f++ {
			idx := r*8 + f
			col := boardColors[idx]
			p := boardPieces[idx]
			if col == White {
				fmt.Printf("%c ", whiteChar[p])
			} else if col == Black {
				fmt.Printf("%c ", blackChar[p])
			} else {
				fmt.Printf(". ")
			}
		}
		fmt.Println()
	}
	fmt.Println("  a b c d e f g h")
}

// PrintMoveTargets zeigt alle Zielfelder für eine Position und Figur
func PrintMoveTargets(square int, targets []int) {
	file := square % 8
	rank := square / 8
	fmt.Printf("Square: %c%d (%d) -> Targets: ", 'a'+byte(file), rank+1, square)
	for _, t := range targets {
		fmt.Printf("%c%d ", 'a'+byte(t%8), t/8+1)
	}
	fmt.Println()
}

func main() {
	// initialize precomputed square score tables
	initSquareScoreTable()

	// initialize board state
	initBoard()

	// show board
	printBoard()

	// example usage
	fmt.Println("squareRank[10] =", squareRank[10])
	fmt.Println("PawnScore[C3] =", PawnScore[10])
	fmt.Println("KnightScore[C3] =", KnightScore[10])

	// white pawn on C3
	fmt.Println("white pawn score on C3 =", squareScoreTable[White][Pawn][10])
	// black pawn on C3 (uses flipped positional table)
	fmt.Println("black pawn score on C3 =", squareScoreTable[Black][Pawn][10])

	// king opening vs endgame on E4 (square index 28)
	fmt.Println("kingOpeningScore white E4 =", kingOpeningScore[White][28])
	fmt.Println("kingEndgameScore black E4 =", kingEndgameScore[Black][28])

	fmt.Print("Debug: Knight from D4 can move to: ")
	for _, t := range KnightTargets[D4] {
		fmt.Print(IndexToAlgebraic(t), " ")
	}
	fmt.Println()
	fmt.Printf("Debug: Square E4 = %d\n", E4)
}
