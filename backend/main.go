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


func getMoves(w http.ResponseWriter, r *http.Request){
	inputFenRep := r.URL.Query().Get("fenRep")
	
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
	log.Println("updatedFenRep: ", updatedFenRep)

	cb := fenToChessboard(updatedFenRep)
	fenRepArr := strings.Split(updatedFenRep,`/`)
	auxData := strings.Split(fenRepArr[len(fenRepArr)-1], ` `)

	parsedTurn, castlingData, enPasSqr, halfMoves, fullMoves := auxData[1], auxData[2], auxData[3], auxData[4], auxData[5]
	log.Println("DATA: ", parsedTurn, castlingData, enPasSqr, halfMoves, fullMoves)
	
	turn := 1
	if parsedTurn == "w"{
		turn = 0
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

	// enPasSqr := []int{2, 3}

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
			log.Println(pieceSymbol, pieceColour)
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
							if pieceSymbolLower == "p" && move[1] != 0  ||  //&& slices.Equal(enPasSqr, posSqr)
							pawnSymbol=="pw" && move[0]==-2 && move[1]==0 && row!=6 || 
							pawnSymbol=="pb" && move[0]==2 && move[1]==0 && row!=1 {
								log.Println("Pawn")
							}else{
								posSqrs = append(posSqrs, posSqr)
							}				

						}else{
							if pieceSymbolLower == "p" && move[1] == 0{
								log.Println("Pawn")
							}else{
								if targetPiece.Colour != cbPiece.Colour{
									posSqrs = append(posSqrs, posSqr)
									log.Println("CAPTURE MOVE")
								}
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
								log.Println("CAPTURE MOVE")
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
			
			if len(posSqrs)>0{
				log.Printf("%s: %v ", rowColStr, outputMap[rowColStr])
			}
		}			
	}
	return outputMap

}

func main() {
	http.HandleFunc("/getMoves", getMoves)
	log.Println("Starting server on port 5669...")
	http.ListenAndServe(":5669", nil)
}