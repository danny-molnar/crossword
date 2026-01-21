package domain

type PuzzleType string

const (
	PuzzleQuick   PuzzleType = "quick"
	PuzzleCryptic PuzzleType = "cryptic"
	PuzzleMixed   PuzzleType = "mixed"
)

type Puzzle struct {
	ID     string
	Title  string
	Type   PuzzleType
	Rows   int
	Cols   int
	Grid   Grid
	Entries []Entry
	Clues   []Clue
}
