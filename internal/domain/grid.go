package domain

type Cell struct {
	R       int
	C       int
	IsBlock bool
	Solution *rune
	IsGiven bool
}

type Grid struct {
	Rows  int
	Cols  int
	Cells [][]Cell
}
