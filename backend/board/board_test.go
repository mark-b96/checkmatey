package board_test

import (
	"backend/board"
	"strings"
	"testing"
)


func TestChessboardToFen(t *testing.T){

	initFen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	expectedFen := strings.Split(initFen,` `)[0]+"/"

	actualCb := board.InitChessboard(initFen).Board

	actualFen := board.ChessboardToFen(&actualCb)

	if actualFen != expectedFen{
		t.Fatal("Conversion failure") 
	}
}
