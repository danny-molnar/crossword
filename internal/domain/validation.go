package domain

import (
	"fmt"
)

type ValidationError struct {
	Problems []string
}

func (e ValidationError) Error() string {
	if len(e.Problems) == 0 {
		return "validation failed"
	}
	// Keep it readable.
	out := "validation failed:"
	for _, p := range e.Problems {
		out += "\n- " + p
	}
	return out
}

func (e *ValidationError) add(format string, args ...any) {
	e.Problems = append(e.Problems, fmt.Sprintf(format, args...))
}

func (e *ValidationError) ok() bool {
	return len(e.Problems) == 0
}

// ValidatePuzzle performs structural validation suitable for both creator and API ingestion.
// It does not enforce "good crossword" rules (like rotational symmetry), only correctness.
func ValidatePuzzle(p Puzzle) error {
	var verr ValidationError

	if p.Rows <= 0 || p.Cols <= 0 {
		verr.add("puzzle dimensions must be > 0, got %dx%d", p.Rows, p.Cols)
	}

	// Grid must match declared dimensions.
	verrGrid := validateGrid(p.Grid, p.Rows, p.Cols)
	if verrGrid != nil {
		for _, msg := range verrGrid.Problems {
			verr.add("%s", msg)
		}
	}

	// Entries can be generated from grid; but if provided, validate them.
	if len(p.Entries) > 0 {
		verrEntries := validateEntries(p.Grid, p.Entries)
		if verrEntries != nil {
			for _, msg := range verrEntries.Problems {
				verr.add("%s", msg)
			}
		}
	}

	// Clues must map to entries if both are present.
	if len(p.Entries) > 0 && len(p.Clues) > 0 {
		verrClues := validateClues(p.Entries, p.Clues)
		if verrClues != nil {
			for _, msg := range verrClues.Problems {
				verr.add("%s", msg)
			}
		}
	}

	if verr.ok() {
		return nil
	}
	return verr
}

func validateGrid(g Grid, rows, cols int) *ValidationError {
	var verr ValidationError

	if g.Rows != 0 && g.Rows != rows {
		verr.add("grid.Rows (%d) does not match puzzle rows (%d)", g.Rows, rows)
	}
	if g.Cols != 0 && g.Cols != cols {
		verr.add("grid.Cols (%d) does not match puzzle cols (%d)", g.Cols, cols)
	}

	if len(g.Cells) != rows {
		verr.add("grid has %d rows of cells, expected %d", len(g.Cells), rows)
		return &verr
	}
	for r := 0; r < rows; r++ {
		if len(g.Cells[r]) != cols {
			verr.add("grid row %d has %d cols, expected %d", r, len(g.Cells[r]), cols)
			continue
		}
		for c := 0; c < cols; c++ {
			cell := g.Cells[r][c]
			if cell.R != 0 || cell.C != 0 {
				// If caller populated R/C, ensure it matches.
				if cell.R != r || cell.C != c {
					verr.add("cell coords mismatch at [%d,%d]: has R=%d C=%d", r, c, cell.R, cell.C)
				}
			}
			if cell.IsBlock {
				if cell.Solution != nil {
					verr.add("block cell [%d,%d] must not have a solution letter", r, c)
				}
				if cell.IsGiven {
					verr.add("block cell [%d,%d] must not be marked as given", r, c)
				}
			}
		}
	}

	if verr.ok() {
		return nil
	}
	return &verr
}

func validateEntries(g Grid, entries []Entry) *ValidationError {
	var verr ValidationError

	seen := map[string]bool{}
	for i, e := range entries {
		if e.Dir != Across && e.Dir != Down {
			verr.add("entry[%d] has invalid direction %q", i, e.Dir)
		}
		if e.Num <= 0 {
			verr.add("entry[%d] has invalid number %d", i, e.Num)
		}
		if len(e.Cells) == 0 {
			verr.add("entry[%d] has no cells", i)
			continue
		}

		// Validate contiguity and non-block membership.
		for j, cr := range e.Cells {
			if cr.R < 0 || cr.C < 0 || cr.R >= g.Rows || cr.C >= g.Cols {
				verr.add("entry[%d] cell[%d] out of bounds: (%d,%d)", i, j, cr.R, cr.C)
				continue
			}
			if g.Cells[cr.R][cr.C].IsBlock {
				verr.add("entry[%d] includes block cell (%d,%d)", i, cr.R, cr.C)
			}
		}

		// Check contiguity by direction.
		for j := 1; j < len(e.Cells); j++ {
			prev := e.Cells[j-1]
			curr := e.Cells[j]
			switch e.Dir {
			case Across:
				if curr.R != prev.R || curr.C != prev.C+1 {
					verr.add("entry[%d] across cells not contiguous at index %d: (%d,%d)->(%d,%d)",
						i, j, prev.R, prev.C, curr.R, curr.C)
				}
			case Down:
				if curr.C != prev.C || curr.R != prev.R+1 {
					verr.add("entry[%d] down cells not contiguous at index %d: (%d,%d)->(%d,%d)",
						i, j, prev.R, prev.C, curr.R, curr.C)
				}
			}
		}

		// Enum validation (if provided).
		if e.Enum != "" {
			en, err := ParseEnum(e.Enum)
			if err != nil {
				verr.add("entry[%d] enum invalid: %v", i, err)
			} else {
				if en.Total != len(e.Cells) {
					verr.add("entry[%d] enum total %d does not match cell count %d", i, en.Total, len(e.Cells))
				}
			}
		}

		// Answer validation (if provided).
		if e.Answer != "" {
			n := NormalizedAnswerLen(e.Answer)
			if n != len(e.Cells) {
				verr.add("entry[%d] answer length %d does not match cell count %d (answer=%q)",
					i, n, len(e.Cells), e.Answer)
			}
		}

		// Uniqueness: number+dir should be unique.
		key := fmt.Sprintf("%s:%d", e.Dir, e.Num)
		if seen[key] {
			verr.add("duplicate entry number for %s %d", e.Dir, e.Num)
		}
		seen[key] = true
	}

	// Cell ownership sanity (optional but useful):
	// Ensure no cell belongs to 2 across entries or 2 down entries.
	type ownerKey struct {
		r, c int
		dir  Direction
	}
	owners := map[ownerKey]string{} // value: entry id/index
	for i, e := range entries {
		for _, cr := range e.Cells {
			ok := ownerKey{r: cr.R, c: cr.C, dir: e.Dir}
			if prev, exists := owners[ok]; exists {
				verr.add("cell (%d,%d) belongs to multiple %s entries (%s and entry[%d])", cr.R, cr.C, e.Dir, prev, i)
			} else {
				owners[ok] = fmt.Sprintf("entry[%d]", i)
			}
		}
	}

	if verr.ok() {
		return nil
	}
	return &verr
}

func validateClues(entries []Entry, clues []Clue) *ValidationError {
	var verr ValidationError

	entryIDs := map[string]bool{}
	// If IDs are empty in early MVP, map by (dir,num) is alternative;
	// but here we assume creator flow will set Entry.ID.
	for _, e := range entries {
		if e.ID != "" {
			entryIDs[e.ID] = true
		}
	}

	for i, c := range clues {
		if c.EntryID == "" {
			verr.add("clue[%d] missing entryId", i)
			continue
		}
		if len(entryIDs) > 0 && !entryIDs[c.EntryID] {
			verr.add("clue[%d] references unknown entryId %q", i, c.EntryID)
		}
		if stringsTrim(c.Text) == "" {
			verr.add("clue[%d] has empty text", i)
		}
	}

	if verr.ok() {
		return nil
	}
	return &verr
}

func stringsTrim(s string) string {
	// tiny helper to avoid importing strings twice at top-level
	i := 0
	j := len(s)
	for i < j && (s[i] == ' ' || s[i] == '\t' || s[i] == '\n' || s[i] == '\r') {
		i++
	}
	for j > i && (s[j-1] == ' ' || s[j-1] == '\t' || s[j-1] == '\n' || s[j-1] == '\r') {
		j--
	}
	return s[i:j]
}
