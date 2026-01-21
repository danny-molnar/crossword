package domain

func GenerateEntries(grid Grid) []Entry {
	var entries []Entry
	num := 1

	isBlock := func(r, c int) bool {
		if r < 0 || c < 0 || r >= grid.Rows || c >= grid.Cols {
			return true
		}
		return grid.Cells[r][c].IsBlock
	}

	for r := 0; r < grid.Rows; r++ {
		for c := 0; c < grid.Cols; c++ {
			if grid.Cells[r][c].IsBlock {
				continue
			}

			startAcross := isBlock(r, c-1) && !isBlock(r, c+1)
			startDown   := isBlock(r-1, c) && !isBlock(r+1, c)

			if startAcross {
				var cells []CellRef
				cc := c
				for cc < grid.Cols && !isBlock(r, cc) {
					cells = append(cells, CellRef{R: r, C: cc})
					cc++
				}
				entries = append(entries, Entry{
					Dir:   Across,
					Num:   num,
					Cells: cells,
				})
			}

			if startDown {
				var cells []CellRef
				rr := r
				for rr < grid.Rows && !isBlock(rr, c) {
					cells = append(cells, CellRef{R: rr, C: c})
					rr++
				}
				entries = append(entries, Entry{
					Dir:   Down,
					Num:   num,
					Cells: cells,
				})
			}

			if startAcross || startDown {
				num++
			}
		}
	}
	return entries
}
