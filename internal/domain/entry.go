package domain

type Direction string

const (
	Across Direction = "across"
	Down   Direction = "down"
)

type Entry struct {
	ID    string
	Dir   Direction
	Num   int
	Cells []CellRef
	Enum  string
	Answer string
}

type CellRef struct {
	R int
	C int
}
