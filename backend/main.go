package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)


type BoardState struct {
	Moves   map[string][][]int `json:"moves"`
	FenRep  string  `json:"fenrep"`
}

func (fS *Fenstate) getInitState(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(fS.FenRep); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}


func (fS *Fenstate) getMoves(w http.ResponseWriter, r *http.Request){
	userMove := ""
	if r.URL.Query().Get("userMove") != ""{
		userMove = r.URL.Query().Get("userMove")
	}

	updateChessboard(fS, userMove)

	moves := calculateMoves(fS)

	if fS.Turn == 1{
		fS.Turn = 0
	}else{
		fS.Turn = 1
	}

	fS.EnPass = "-"

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.WriteHeader(http.StatusOK)

	payload := BoardState{Moves: moves, FenRep: fS.FenRep}

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

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

	initCB := initChessboard(initFen)

	initFenRep := &Fenstate{
		FenRep: initFen, 
		Turn: turnInt, 
		CastlingStatus: castlingData, 
		EnPass: enPasSqr, 
		HalfMoves: halfMovesInt, 
		FullMoves: fullMovesInt,
		cb: initCB,
	} 
	
	
	http.HandleFunc("/getMoves", initFenRep.getMoves)
	http.HandleFunc("/getInitState", initFenRep.getInitState)
	log.Println("Starting server on port 5669...")
	http.ListenAndServe(":5669", nil)
}