package main

// PawnScore: pawn positional score for each square (A1..H1, A2..H2, ..., A8..H8)
// Values scaled to a 0..100-ish range; stronger toward center and advanced ranks.
var PawnScore = [...]int{
	// Rank 1
	0, 0, 0, 0, 0, 0, 0, 0,
	// Rank 2
	5, 10, 15, 2, 2, 15, 10, 5,
	// Rank 3
	10, 15, 25, 30, 30, 25, 15, 10,
	// Rank 4
	20, 30, 40, 50, 50, 40, 30, 20,
	// Rank 5
	30, 45, 60, 70, 70, 60, 45, 30,
	// Rank 6
	50, 70, 85, 95, 95, 85, 70, 50,
	// Rank 7
	80, 90, 95, 100, 100, 95, 90, 80,
	// Rank 8
	0, 0, 0, 0, 0, 0, 0, 0,
}

// KnightScore: typical knight-centralization table (negative on edges/corners)
var KnightScore = [...]int{
	// Rank 1
	-50, -40, -30, -30, -30, -30, -40, -50,
	// Rank 2
	-40, -20, 0, 5, 5, 0, -20, -40,
	// Rank 3
	-30, 5, 10, 15, 15, 10, 5, -30,
	// Rank 4
	-30, 0, 15, 20, 20, 15, 0, -30,
	// Rank 5
	-30, 5, 15, 20, 20, 15, 5, -30,
	// Rank 6
	-30, 0, 10, 15, 15, 10, 0, -30,
	// Rank 7
	-40, -20, 0, 0, 0, 0, -20, -40,
	// Rank 8
	-50, -40, -30, -30, -30, -30, -40, -50,
}

// BishopScore: favors long diagonals and center
var BishopScore = [...]int{
	// Rank 1
	-20, -10, -10, -10, -10, -10, -10, -20,
	// Rank 2
	-10, 0, 0, 0, 0, 0, 0, -10,
	// Rank 3
	-10, 0, 5, 10, 10, 5, 0, -10,
	// Rank 4
	-10, 5, 5, 10, 10, 5, 5, -10,
	// Rank 5
	-10, 0, 10, 10, 10, 10, 0, -10,
	// Rank 6
	-10, 10, 10, 10, 10, 10, 10, -10,
	// Rank 7
	-10, 5, 0, 0, 0, 0, 5, -10,
	// Rank 8
	-20, -10, -10, -10, -10, -10, -10, -20,
}

// RookScore: favors open files and ranks closer to opponent
var RookScore = [...]int{
	// Rank 1
	0, 0, 0, 5, 5, 0, 0, 0,
	// Rank 2
	5, 10, 10, 10, 10, 10, 10, 5,
	// Rank 3
	-5, 0, 0, 0, 0, 0, 0, -5,
	// Rank 4
	-5, 0, 0, 0, 0, 0, 0, -5,
	// Rank 5
	-5, 0, 0, 0, 0, 0, 0, -5,
	// Rank 6
	-5, 0, 0, 0, 0, 0, 0, -5,
	// Rank 7
	-5, 0, 0, 0, 0, 0, 0, -5,
	// Rank 8
	0, 0, 0, 0, 0, 0, 0, 0,
}

// QueenScore: combines mobility and centralization
var QueenScore = [...]int{
	// Rank 1
	-20, -10, -10, -5, -5, -10, -10, -20,
	// Rank 2
	-10, 0, 0, 0, 0, 0, 0, -10,
	// Rank 3
	-10, 0, 5, 5, 5, 5, 0, -10,
	// Rank 4
	-5, 0, 5, 5, 5, 5, 0, -5,
	// Rank 5
	0, 0, 5, 5, 5, 5, 0, -5,
	// Rank 6
	-10, 0, 5, 5, 5, 5, 0, -10,
	// Rank 7
	-10, 0, 0, 0, 0, 0, 0, -10,
	// Rank 8
	-20, -10, -10, -5, -5, -10, -10, -20,
}

// KingScore: simple middlegame table (prefer safety; encourage centralization slightly in endgame)
var KingScore = [...]int{
	// Rank 1
	-30, -40, -40, -50, -50, -40, -40, -30,
	// Rank 2
	-30, -40, -40, -50, -50, -40, -40, -30,
	// Rank 3
	-30, -40, -40, -50, -50, -40, -40, -30,
	// Rank 4
	-30, -40, -40, -50, -50, -40, -40, -30,
	// Rank 5
	-20, -30, -30, -40, -40, -30, -30, -20,
	// Rank 6
	-10, -20, -20, -20, -20, -20, -20, -10,
	// Rank 7
	20, 20, 0, 0, 0, 0, 20, 20,
	// Rank 8
	20, 30, 10, 0, 0, 10, 30, 20,
}

// KingEndgameScore: king positional table for the endgame where the king is stronger in the centre
var KingEndgameScore = [...]int{
	// Rank 1
	-40, -30, -20, -10, -10, -20, -30, -40,
	// Rank 2
	-30, -20, -10, 0, 0, -10, -20, -30,
	// Rank 3
	-20, -10, 10, 20, 20, 10, -10, -20,
	// Rank 4
	-10, 0, 25, 35, 35, 25, 0, -10,
	// Rank 5
	-10, 0, 25, 35, 35, 25, 0, -10,
	// Rank 6
	-20, -10, 10, 20, 20, 10, -10, -20,
	// Rank 7
	-30, -20, -10, 0, 0, -10, -20, -30,
	// Rank 8
	-40, -30, -20, -10, -10, -20, -30, -40,
}

// PassedPawnScore: bonus for passed pawns. More advanced = much more valuable;
// differences between successive ranks grow as pawns advance.
// Order: A1..H1, A2..H2, ..., A8..H8
var PassedPawnScore = [...]int{
	// Rank 1 (A1..H1) - irrelevant for passed pawn
	0, 0, 0, 0, 0, 0, 0, 0,
	// Rank 2 (A2..H2)
	0, 0, 0, 0, 0, 0, 0, 0,
	// Rank 3 (A3..H3)
	60, 60, 60, 60, 60, 60, 60, 60,
	// Rank 4 (A4..H4)
	30, 30, 30, 30, 30, 30, 30, 30,
	// Rank 5 (A5..H5)
	15, 15, 15, 15, 15, 15, 15, 15,
	// Rank 6 (A6..H6)
	8, 8, 8, 8, 8, 8, 8, 8,
	// Rank 7 (A7..H7) - very strong
	8, 8, 8, 8, 8, 8, 8, 8,
	// Rank 8 (A8..H8) - promotion square (handled separately)
	0, 0, 0, 0, 0, 0, 0, 0,
}
