package main


type piece struct {
	row int
	col int
	colour string
	moves [][]int
}

func newPawn(row int, col int, colour string) *piece {
	p := piece{row: row, col: col, colour: colour}
	// p.moves := [1, -1], [1, 0], [1, 1]
	return &p
}
