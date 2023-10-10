package tetris

// BoardRows is the height of the game board in terms of blocks
const BoardRows = 22

// BoardCols is the width of the game board in terms of blocks
const BoardCols = 10

// Point represents a coordinate on the game board with Point{row:0, col:0}
// representing the bottom left
type Point struct {
	row int
	col int
}

// Block represents the color of the block
type Block int

// Different values a point on the grid can hold
const (
	Empty Block = iota
	Goluboy
	Siniy
	Pink
	Purple
	Red
	Yellow
	Green
	Gray
	GoluboySpecial
	SiniySpecial
	PinkSpecial
	PurpleSpecial
	RedSpecial
	YellowSpecial
	GreenSpecial
	GraySpecial
)

// Piece is a constant for a shape of piece. There are 7 classic pieces like L, and O
type Piece int

// Various values that the pieces can be
const (
	IPiece Piece = iota
	JPiece
	LPiece
	OPiece
	SPiece
	TPiece
	ZPiece
)

// Shape is a type containing four points, which represents the four points
// making a contiguous 'piece'.
type Shape [4]Point

const levelLength = 60.0 // Time it takes for game to speed up
const speedUpRate = 0.1  // Every new level, the amount the game speeds up by
