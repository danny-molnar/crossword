package domain

import "testing"

func makeGrid(rows, cols int, blocks map[[2]int]bool) Grid {
	cells := make([][]Cell, rows)
	for r := 0; r < rows; r++ {
		cells[r] = make([]Cell, cols)
		for c := 0; c < cols; c++ {
			isBlock := blocks[[2]int{r, c}]
			cells[r][c] = Cell{R: r, C: c, IsBlock: isBlock}
		}
	}
	return Grid{Rows: rows, Cols: cols, Cells: cells}
}

func TestValidatePuzzle_OK_Minimal(t *testing.T) {
	g := makeGrid(5, 5, map[[2]int]bool{
		{0, 1}: true,
		{1, 1}: true,
		{3, 3}: true,
	})

	p := Puzzle{
		ID:    "p1",
		Title: "Test",
		Type:  PuzzleQuick,
		Rows:  5,
		Cols:  5,
		Grid:  g,
	}

	if err := ValidatePuzzle(p); err != nil {
		t.Fatalf("expected OK, got %v", err)
	}
}

func TestValidatePuzzle_BadGridShape(t *testing.T) {
	g := makeGrid(5, 5, nil)
	// Break shape: drop a column in row 2
	g.Cells[2] = g.Cells[2][:4]

	p := Puzzle{ID: "p1", Title: "Test", Type: PuzzleQuick, Rows: 5, Cols: 5, Grid: g}

	if err := ValidatePuzzle(p); err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestValidateEntries_EnumAndAnswerMatch(t *testing.T) {
	g := makeGrid(5, 5, map[[2]int]bool{
		{0, 1}: true,
	})

	entries := []Entry{
		{
			ID:     "e1",
			Dir:    Across,
			Num:    1,
			Cells:  []CellRef{{0, 0}, {0, 2}, {0, 3}}, // intentionally non-contiguous
			Enum:   "3",
			Answer: "CAT",
		},
	}

	p := Puzzle{ID: "p1", Title: "Test", Type: PuzzleQuick, Rows: 5, Cols: 5, Grid: g, Entries: entries}

	if err := ValidatePuzzle(p); err == nil {
		t.Fatalf("expected contiguity error, got nil")
	}
}

func TestValidateEntries_OK_ContiguousAcross(t *testing.T) {
	g := makeGrid(5, 5, map[[2]int]bool{
		{0, 1}: true,
	})

	entries := []Entry{
		{
			ID:     "e1",
			Dir:    Across,
			Num:    1,
			Cells:  []CellRef{{0, 2}, {0, 3}, {0, 4}},
			Enum:   "3",
			Answer: "CAT",
		},
	}

	p := Puzzle{ID: "p1", Title: "Test", Type: PuzzleQuick, Rows: 5, Cols: 5, Grid: g, Entries: entries}

	if err := ValidatePuzzle(p); err != nil {
		t.Fatalf("expected OK, got %v", err)
	}
}
