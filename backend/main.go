package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"unicode"
	"strconv"
)


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
	inputFenRep := fS.FenRep

	if r.URL.Query().Get("fenRep") != ""{
		inputFenRep = r.URL.Query().Get("fenRep")
	}
		
	var updatedFenRep = ""

	for index:= range inputFenRep{
		if unicode.IsDigit(rune(inputFenRep[index])) && inputFenRep[index]!='0'{
			num, _ := strconv.Atoi(string(rune(inputFenRep[index])))
			for range num{
				updatedFenRep += string(inputFenRep[index])
			}
		}else{
			updatedFenRep+=string(inputFenRep[index])
		}
	}
	
	fS.FenRep = updatedFenRep

	cb := fenToChessboard(updatedFenRep)

	turn := fS.Turn
	
	if fS.Turn == 1{
		fS.Turn = 0
	}else{
		fS.Turn = 1
	}
	
	moves := calculateMoves(cb, turn)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(moves); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}


}

type Piece struct {
	Row int
	Col int
	Colour int
	Symbol string
}

type Square struct {
	Row int
	Col int
	Piece Piece
}

type Chessboard struct {
	board [8][8]Square
}

type Fenstate struct {
	FenRep string
	Turn int
	CastlingStatus string
	EnPass string
	HalfMoves int
	FullMoves int
}

func fenToChessboard(fenRep string) Chessboard{
	fenRepArr := strings.Split(fenRep,`/`)
	var currPos [8][8]Square
	var width = 8
	var height = 8
	
	for row:= range height{
		fenRow := fenRepArr[row]

		for col:= range width{
			newSquare := Square{Row: row, Col: col}
			
			if col <= len(fenRow) -1 {
				pieceSymbol := fenRow[col]
				pieceColour := 1
				if unicode.IsUpper(rune(pieceSymbol)) {
					pieceColour = 0
				}
				if !unicode.IsDigit(rune(pieceSymbol)){
					newSquare.Piece = Piece{Row: row, Col: col, Colour: pieceColour, Symbol: string(pieceSymbol)}
				}
			}
			currPos[row][col] = newSquare
		}
	}
	newChessboard := Chessboard{board: currPos}

	return newChessboard
}

func chessboardToFen(cb Chessboard) string{
	fenRep := ""
	for row:= range len(cb.board){
		for col:= range len(cb.board[0]){
			cbSqr := cb.board[row][col]
			cbPiece := cbSqr.Piece
			fenRep += cbPiece.Symbol
		}
		fenRep += "/"
	}
	return fenRep
}

func calculateMoves(cb Chessboard, turn int) map[string][][]int{
	pieceMap := map[string][][2]int {
		"pb": [][2]int{{1,0}, {2, 0}, {1,1},{1,-1},},
		"pw": [][2]int{{-1,0}, {-2, 0}, {-1,1},{-1,-1},},
		"r": [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}},
		"b": [][2]int{{1, 1}, {-1, -1}, {-1, 1}, {1, -1}},
		"q": [][2]int{{1, 0}, {1, 1}, {-1, 0}, {-1, -1}, {-1, 1}, {1, -1}, {0, 1}, {0, -1}},
		"n": [][2]int{{1,2}, {2,1}, {-1,2}, {-2,1}, {1,-2}, {2,-1}, {-1,-2}, {-2,-1}},
		"k": [][2]int{{1, 0}, {1, 1}, {-1, 0}, {-1, -1}, {-1, 1}, {1, -1}, {0, 1}, {0, -1}},
	}

	multiMoves := map[string]bool {
		"r": true,
		"b": true,
		"q": true,
		"k": false,
		"n": false,
		"p": false,
	}

	outputMap := map[string][][]int{}

	for row:= range len(cb.board){
		for col:= range len(cb.board[0]){
			cbSqr := cb.board[row][col]
			cbPiece := cbSqr.Piece
			pieceSymbol := cbPiece.Symbol
			pieceColour := cbPiece.Colour
			pieceSymbolLower := strings.ToLower(pieceSymbol)

			if len(pieceSymbol)==0 || pieceColour != turn{
				continue
			}
			moves:= [][2]int{}

			pawnSymbol:=""
			if pieceSymbol == "P"{
				pawnSymbol="pw"
			}
			if pieceSymbol == "p"{
				pawnSymbol="pb"
			}
			
			if pawnSymbol != ""{
				moves = pieceMap[pawnSymbol]
			}else{
				moves = pieceMap[pieceSymbolLower]
			}

			var posSqrs [][]int
			for index:= range(moves){
				move := moves[index]
				posRow := row+move[0]
				posCol := col+move[1]
				posSqr := []int{posRow, posCol}

				if !multiMoves[pieceSymbolLower]{
					if posCol>=0 && posRow>=0 && posCol<=7 && posRow<=7{
						targetPiece := cb.board[posRow][posCol].Piece
						if targetPiece.Symbol==""{
							if !(pieceSymbolLower == "p" && move[1] != 0) &&
							!(pawnSymbol=="pw" && move[0]==-2 && move[1]==0 && row!=6) &&
							!(pawnSymbol=="pb" && move[0]==2 && move[1]==0 && row!=1) {
								posSqrs = append(posSqrs, posSqr)
							}
						}else{
							if targetPiece.Colour != cbPiece.Colour && 
							!(pieceSymbolLower == "p" && move[1] == 0){
								posSqrs = append(posSqrs, posSqr)
							}	
						}
					}
				}else{				
					for posCol>=0 && posRow>=0 && posCol<=7 && posRow<=7{
						targetPiece := cb.board[posRow][posCol].Piece
						newPosSqr := []int{posRow, posCol}
						if targetPiece.Symbol==""{
							posSqrs = append(posSqrs, newPosSqr)
						}else{
							if targetPiece.Colour != cbPiece.Colour{
								posSqrs = append(posSqrs, newPosSqr)
							}
							break
						}
						posRow += move[0]
						posCol += move[1]
					}
				}
			}
			rowColStr := fmt.Sprintf("%d,%d", row, col)
			outputMap[rowColStr] = posSqrs
		}			
	}
	return outputMap

}

func main() {
	initFen := "rnbqkbnr/pppppppp/8/2k3k1/2n6/6n1/PP1PPP1P/RNBQKBNR w KQkq - 0 1"
	fenRepArr := strings.Split(initFen,`/`)
	auxData := strings.Split(fenRepArr[len(fenRepArr)-1], ` `)

	parsedTurn, castlingData, enPasSqr, halfMoves, fullMoves := auxData[1], auxData[2], auxData[3], auxData[4], auxData[5]

	halfMovesInt, _ := strconv.Atoi(halfMoves)
	fullMovesInt, _ := strconv.Atoi(fullMoves)
	turnInt := 1
	if parsedTurn == "w"{
		turnInt = 0
	}


	initFenRep := &Fenstate{FenRep: initFen, Turn: turnInt, CastlingStatus: castlingData, EnPass: enPasSqr, HalfMoves: halfMovesInt, FullMoves: fullMovesInt} 
	http.HandleFunc("/getMoves", initFenRep.getMoves)
	http.HandleFunc("/getInitState", initFenRep.getInitState)
	log.Println("Starting server on port 5669...")
	http.ListenAndServe(":5669", nil)
}