package domain

import "time"

type SolveSession struct {
	ID        string    `json:"id"`
	PuzzleID  string    `json:"puzzleId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	// GridState stores the user's current fill.
	// Key format: "r,c" -> cell value (single letter or empty string)
	GridState map[string]string `json:"gridState"`

	// Pencil marks (optional MVP)
	Pencil map[string]bool `json:"pencil"`

	ChecksUsed  int `json:"checksUsed"`
	RevealsUsed int `json:"revealsUsed"`
}
