package domain

type PuzzlePublic struct {
	ID      string        `json:"id"`
	Title   string        `json:"title"`
	Type    PuzzleType    `json:"type"`
	Rows    int           `json:"rows"`
	Cols    int           `json:"cols"`
	Grid    GridPublic    `json:"grid"`
	Entries []EntryPublic `json:"entries"`
	Clues   []CluePublic  `json:"clues"`
}

type GridPublic struct {
	Rows  int            `json:"rows"`
	Cols  int            `json:"cols"`
	Cells [][]CellPublic `json:"cells"`
}

type CellPublic struct {
	R     int  `json:"r"`
	C     int  `json:"c"`
	Block bool `json:"block"`
	Given bool `json:"given"`
}

type EntryPublic struct {
	ID    string    `json:"id"`
	Dir   Direction `json:"dir"`
	Num   int       `json:"num"`
	Cells []CellRef `json:"cells"`
	Enum  string    `json:"enum"`
}

type CluePublic struct {
	EntryID     string   `json:"entryId"`
	Text        string   `json:"text"`
	Explanation *string  `json:"explanation,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}
