package main

import "fmt"

// pawnMovesWhite defines the possible forward moves for a white pawn (relative offsets)
// White pawns move towards rank 8 (increasing rank)
var pawnMovesWhite = [...]int{
	8, // one square forward
}

// pawnCapturesWhite defines the possible capture moves for a white pawn (relative offsets)
// White pawns capture diagonally forward
var pawnCapturesWhite = [...]int{
	7, // diagonally forward-left
	9, // diagonally forward-right
}

// pawnMovesBlack defines the possible forward moves for a black pawn (relative offsets)
// Black pawns move towards rank 1 (decreasing rank)
var pawnMovesBlack = [...]int{
	-8, // one square forward
}

// pawnCapturesBlack defines the possible capture moves for a black pawn (relative offsets)
// Black pawns capture diagonally forward
var pawnCapturesBlack = [...]int{
	-7, // diagonally forward-left
	-9, // diagonally forward-right
}

// 2D arrays with all possible target squares for each piece type and source square
// Index: [sourceSquare][targetSquare index within max possible moves]
// Each piece can have up to 27 possible moves (Queen has most)

// PawnTargetsWhite[sourceSquare] contains all legal target squares for a white pawn from that source
var PawnTargetsWhite [64][]int

// PawnTargetsBlack[sourceSquare] contains all legal target squares for a black pawn from that source
var PawnTargetsBlack [64][]int

// KnightTargets[sourceSquare] contains all legal target squares for a knight from that source
var KnightTargets [64][]int

// BishopTargets[sourceSquare] contains all legal target squares for a bishop from that source
var BishopTargets [64][]int

// RookTargets[sourceSquare] contains all legal target squares for a rook from that source
var RookTargets [64][]int

// QueenTargets[sourceSquare] contains all legal target squares for a queen from that source
var QueenTargets [64][]int

// KingTargets[sourceSquare] contains all legal target squares for a king from that source
var KingTargets [64][]int

// knightMoves defines all possible moves for a knight (relative offsets)
// Knights move in an L-shape: 2 squares in one direction, 1 square perpendicular
var knightMoves = [...]int{
	6,   // up 2, right 1
	10,  // up 2, left 1 (across file boundary - needs check)
	15,  // up 1, right 2
	17,  // up 1, left 2 (across file boundary - needs check)
	-6,  // down 2, left 1 (across file boundary - needs check)
	-10, // down 2, right 1
	-15, // down 1, left 2 (across file boundary - needs check)
	-17, // down 1, right 2
}

// bishopMoves defines all possible move directions for a bishop (relative offsets)
// Bishops move diagonally - represented as directions that can be repeated
var bishopDirections = [...]int{
	7,  // up-left diagonal
	9,  // up-right diagonal
	-7, // down-right diagonal
	-9, // down-left diagonal
}

// rookMoves defines all possible move directions for a rook (relative offsets)
// Rooks move horizontally or vertically - represented as directions that can be repeated
var rookDirections = [...]int{
	1,  // right
	-1, // left
	8,  // up
	-8, // down
}

// queenMoves defines all possible move directions for a queen (relative offsets)
// Queens move like both bishops and rooks - horizontally, vertically, and diagonally
var queenDirections = [...]int{
	1,  // right
	-1, // left
	8,  // up
	-8, // down
	7,  // up-left diagonal
	9,  // up-right diagonal
	-7, // down-right diagonal
	-9, // down-left diagonal
}

// kingMoves defines all possible moves for a king (relative offsets)
// Kings move exactly one square in any direction
var kingMoves = [...]int{
	1,  // right
	-1, // left
	8,  // up
	-8, // down
	7,  // up-left diagonal
	9,  // up-right diagonal
	-7, // down-right diagonal
	-9, // down-left diagonal
}

// castlingKingSideWhite represents the king-side castling move for white (king E1 to G1)
const castlingKingSideWhite = 2 // relative offset from E1

// castlingQueenSideWhite represents the queen-side castling move for white (king E1 to C1)
const castlingQueenSideWhite = -2 // relative offset from E1

// castlingKingSideBlack represents the king-side castling move for black (king E8 to G8)
const castlingKingSideBlack = 2 // relative offset from E8

// castlingQueenSideBlack represents the queen-side castling move for black (king E8 to C8)
const castlingQueenSideBlack = -2 // relative offset from E8

// GetPieceMoveDirections returns the move directions for a piece type
// For sliding pieces (Bishop, Rook, Queen), returns directions that can be repeated
// For non-sliding pieces (Pawn, Knight, King), call GetPieceMoves instead
func GetPieceMoveDirections(pieceType int) []int {
	switch pieceType {
	case Bishop:
		return bishopDirections[:]
	case Rook:
		return rookDirections[:]
	case Queen:
		return queenDirections[:]
	default:
		return nil
	}
}

// GetPieceMoves returns all possible moves for a piece type
// For sliding pieces, this returns individual move offsets
func GetPieceMoves(pieceType int, color int) []int {
	switch pieceType {
	case Pawn:
		if color == White {
			// Could combine regular moves and captures
			// For now, return separately as they have different validation rules
			return pawnMovesWhite[:]
		} else {
			return pawnMovesBlack[:]
		}
	case Knight:
		return knightMoves[:]
	case King:
		return kingMoves[:]
	default:
		return nil
	}
}

// GetPawnCaptures returns all possible capture moves for a pawn
func GetPawnCaptures(color int) []int {
	if color == White {
		return pawnCapturesWhite[:]
	} else {
		return pawnCapturesBlack[:]
	}
}

// isSquareValid checks if a square index is valid (within board boundaries)
func isSquareValid(sq int) bool {
	return sq >= 0 && sq < 64
}

// isFileWrappingMove checks if a move wraps around file boundaries (invalid L-shape moves)
// This is used for knight moves to prevent wrapping from h-file to a-file
func isFileWrappingMove(from, to int) bool {
	fromFile := from % 8
	toFile := to % 8
	// Detect wrapping (e.g., from H-file to A-file or vice versa)
	if fromFile == 0 && toFile >= 6 { // A-file wrapping right
		return true
	}
	if fromFile == 7 && toFile <= 1 { // H-file wrapping left
		return true
	}
	if fromFile <= 1 && toFile >= 6 { // Left files wrapping right
		return true
	}
	if fromFile >= 6 && toFile <= 1 { // Right files wrapping left
		return true
	}
	return false
}

// initMoveTargets generates all target squares for each piece type from each source square
func initMoveTargets() {
	// Initialize Pawn targets
	for sq := 0; sq < 64; sq++ {
		rank := sq / 8
		file := sq % 8

		// White pawns
		var whiteMoves []int
		// Forward move (only if not on rank 8)
		if rank < 7 {
			whiteMoves = append(whiteMoves, sq+8)
		}
		// Captures (diagonal moves - only if not on rank 8)
		if rank < 7 {
			if file > 0 {
				whiteMoves = append(whiteMoves, sq+7) // left-up diagonal
			}
			if file < 7 {
				whiteMoves = append(whiteMoves, sq+9) // right-up diagonal
			}
		}
		PawnTargetsWhite[sq] = whiteMoves

		// Black pawns
		var blackMoves []int
		// Forward move (only if not on rank 1)
		if rank > 0 {
			blackMoves = append(blackMoves, sq-8)
		}
		// Captures (diagonal moves - only if not on rank 1)
		if rank > 0 {
			if file > 0 {
				blackMoves = append(blackMoves, sq-9) // left-down diagonal
			}
			if file < 7 {
				blackMoves = append(blackMoves, sq-7) // right-down diagonal
			}
		}
		PawnTargetsBlack[sq] = blackMoves
	}

	// Initialize Knight targets
	for sq := 0; sq < 64; sq++ {
		var targets []int
		for _, offset := range knightMoves {
			target := sq + offset
			if isSquareValid(target) && !isFileWrappingMove(sq, target) {
				targets = append(targets, target)
			}
		}
		KnightTargets[sq] = targets
	}

	// Initialize Bishop targets (sliding piece - all diagonals until edge)
	for sq := 0; sq < 64; sq++ {
		var targets []int
		for _, direction := range bishopDirections {
			// Slide in this direction until we hit the edge
			target := sq + direction
			for isSquareValid(target) && !isFileWrappingMove(target-direction, target) {
				targets = append(targets, target)
				target += direction
			}
		}
		BishopTargets[sq] = targets
	}

	// Initialize Rook targets (sliding piece - all straight lines until edge)
	for sq := 0; sq < 64; sq++ {
		var targets []int
		for _, direction := range rookDirections {
			// Slide in this direction until we hit the edge
			target := sq + direction
			for isSquareValid(target) && !isFileWrappingMove(target-direction, target) {
				targets = append(targets, target)
				target += direction
			}
		}
		RookTargets[sq] = targets
	}

	// Initialize Queen targets (sliding piece - all directions until edge)
	for sq := 0; sq < 64; sq++ {
		var targets []int
		for _, direction := range queenDirections {
			// Slide in this direction until we hit the edge
			target := sq + direction
			for isSquareValid(target) && !isFileWrappingMove(target-direction, target) {
				targets = append(targets, target)
				target += direction
			}
		}
		QueenTargets[sq] = targets
	}

	// Initialize King targets
	for sq := 0; sq < 64; sq++ {
		var targets []int
		for _, offset := range kingMoves {
			target := sq + offset
			if isSquareValid(target) && !isFileWrappingMove(sq, target) {
				targets = append(targets, target)
			}
		}
		KingTargets[sq] = targets
	}
}

// IndexToAlgebraic converts a board index (0..63) to algebraic notation ("a1".."h8").
// Returns an empty string for invalid indices.
func IndexToAlgebraic(idx int) string {
	if !isSquareValid(idx) {
		return ""
	}
	file := idx % 8
	rank := idx / 8
	return fmt.Sprintf("%c%d", 'a'+byte(file), rank+1)
}
