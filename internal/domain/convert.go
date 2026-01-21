package domain

func ToPublic(p Puzzle) PuzzlePublic {
	pub := PuzzlePublic{
		ID:      p.ID,
		Title:   p.Title,
		Type:    p.Type,
		Rows:    p.Rows,
		Cols:    p.Cols,
		Grid:    GridPublic{Rows: p.Grid.Rows, Cols: p.Grid.Cols},
		Entries: make([]EntryPublic, 0, len(p.Entries)),
		Clues:   make([]CluePublic, 0, len(p.Clues)),
	}

	// Grid cells
	if p.Grid.Rows > 0 && p.Grid.Cols > 0 && len(p.Grid.Cells) == p.Grid.Rows {
		pub.Grid.Cells = make([][]CellPublic, p.Grid.Rows)
		for r := 0; r < p.Grid.Rows; r++ {
			pub.Grid.Cells[r] = make([]CellPublic, p.Grid.Cols)
			for c := 0; c < p.Grid.Cols; c++ {
				cell := p.Grid.Cells[r][c]
				pub.Grid.Cells[r][c] = CellPublic{
					R:     r,
					C:     c,
					Block: cell.IsBlock,
					Given: cell.IsGiven,
				}
			}
		}
	}

	// Entries (no answers)
	for _, e := range p.Entries {
		pub.Entries = append(pub.Entries, EntryPublic{
			ID:    e.ID,
			Dir:   e.Dir,
			Num:   e.Num,
			Cells: e.Cells,
			Enum:  e.Enum,
		})
	}

	// Clues (explanations optional)
	for _, c := range p.Clues {
		pub.Clues = append(pub.Clues, CluePublic{
			EntryID:     c.EntryID,
			Text:        c.Text,
			Explanation: c.Explanation,
			Tags:        c.Tags,
		})
	}

	return pub
}
