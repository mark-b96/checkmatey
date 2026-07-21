package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"backend/board"
)


func main() {
	initFen := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	fenRepArr := strings.Split(initFen,`/`)
	auxData := strings.Split(fenRepArr[len(fenRepArr)-1], ` `)

	parsedTurn, castlingData, enPasSqr, halfMoves, fullMoves := auxData[1], auxData[2], auxData[3], auxData[4], auxData[5]

	halfMovesInt, _ := strconv.Atoi(halfMoves)
	fullMovesInt, _ := strconv.Atoi(fullMoves)
	turnInt := 1
	if parsedTurn == "w"{
		turnInt = 0
	}

	initCB := board.InitChessboard(initFen)

	initFenRep := &board.Fenstate{
		FenRep: initFen, 
		Turn: turnInt, 
		CastlingStatus: castlingData, 
		EnPass: enPasSqr, 
		HalfMoves: halfMovesInt, 
		FullMoves: fullMovesInt,
		CB: initCB,
	} 
	
	
	http.HandleFunc("/getMoves", initFenRep.GetMoves)
	http.HandleFunc("/getInitState", initFenRep.GetInitState)
	log.Println("Starting server on port 5669...")
	http.ListenAndServe(":5669", nil)
}
