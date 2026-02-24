package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"unicode"
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
	cb Chessboard
}


func initChessboard(fenRep string) Chessboard{
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


func updateChessboard(fS * Fenstate, userMove string){
	if len(userMove)==0{
		return
	}

	cb := &fS.cb.board
	
	parsedMove := strings.Split(userMove, ":")

	srcPiece, moveSqrs := parsedMove[0], parsedMove[1]

	parsedSqrs := strings.Split(moveSqrs, "-")

	srcSquare, dstSquare := parsedSqrs[0], parsedSqrs[1]

	srcCoords := strings.Split(srcSquare, ",")
	dstCoords := strings.Split(dstSquare, ",")

	srcSquareRow, srcSquareCol := srcCoords[0], srcCoords[1]
	dstSquareRow, dstSquareCol := dstCoords[0], dstCoords[1]
	
	srcSquareRowInt, _ := strconv.Atoi(srcSquareRow)
	srcSquareColInt, _ := strconv.Atoi(srcSquareCol)
	dstSquareRowInt, _ := strconv.Atoi(dstSquareRow)
	dstSquareColInt, _ := strconv.Atoi(dstSquareCol)

	pieceColour := 1
	if unicode.IsUpper(rune(srcPiece[0])) {
		pieceColour = 0
	}


	cb[srcSquareRowInt][srcSquareColInt].Piece = Piece{}
	cb[dstSquareRowInt][dstSquareColInt].Piece = Piece{Row: dstSquareRowInt, Col: dstSquareColInt, Colour: pieceColour, Symbol: srcPiece}

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

func getCastlingStatus(cb Chessboard, fS *Fenstate) string{
	if fS.CastlingStatus == ""{
		return ""
	}

	castlingStr := ""

	if fS.Turn == 0{
		if (cb.board[7][4].Piece.Symbol == "K"){
			if (cb.board[7][7].Piece.Symbol == "R") && 
			cb.board[7][6].Piece.Symbol == "" && 
			cb.board[7][5].Piece.Symbol == ""{
				castlingStr += "K"
			}
			if (cb.board[7][0].Piece.Symbol == "R") && 
			cb.board[7][1].Piece.Symbol == "" && 
			cb.board[7][2].Piece.Symbol == "" &&
			cb.board[7][3].Piece.Symbol == "" {
				castlingStr += "Q"
			}
		}
	}else{
		if (cb.board[0][4].Piece.Symbol == "k"){
			if (cb.board[0][7].Piece.Symbol == "r") && 
			cb.board[0][6].Piece.Symbol == "" && 
			cb.board[0][5].Piece.Symbol == ""{
				castlingStr += "k"
			}
			if (cb.board[0][0].Piece.Symbol == "r") && 
			cb.board[0][1].Piece.Symbol == "" && 
			cb.board[0][2].Piece.Symbol == "" &&
			cb.board[0][3].Piece.Symbol == "" {
				castlingStr += "q"
			}
		}
	}

	return castlingStr

}


func canMovePawn(pawnSymbol string, move [2]int, row int, col int, cb Chessboard, enPassSqr Square, posRow int, posCol int) bool{
	if move[1] == 0{
		if pawnSymbol=="pw" && move[0]==-2 && row==6 && cb.board[row-1][col].Piece.Symbol==""{
			return true
		}
		if pawnSymbol=="pb" && move[0]==2 && row==1 && cb.board[row+1][col].Piece.Symbol==""{
			return true
		}
		if move[0] == 1 || move[0]==-1{
			return true
		}
	}else{
		if move[1] == 1 || move[1]==-1{		
			if enPassSqr.Row == posRow && enPassSqr.Col == posCol{
				return true
			}
		}
	}
	return false
}

func canCastle(pieceSymbol string, move [2]int, castlingStatus string) bool{
	if (pieceSymbol == "K"){
		if move[1]==2{
			if strings.Contains(castlingStatus, "K"){
				return true
			}else{
				return false
			}
		}
		if move[1]==-2{
			if strings.Contains(castlingStatus, "Q"){
				return true
			}else{
				return false
			}
		}
	}

	if (pieceSymbol == "k"){
		if move[1]==2{
			if strings.Contains(castlingStatus, "k"){
				return true
			}else{
				return false
			}
		}
		if move[1]==-2{
			if strings.Contains(castlingStatus, "q"){
				return true
			}else{
				return false
			}
		}
	}
	return false
}

func calculateMoves(fS *Fenstate) map[string][][]int{
	cb := fS.cb
	// log.Println(cb.board)
	pieceMap := map[string][][2]int {
		"pb": [][2]int{{1,0}, {2, 0}, {1,1},{1,-1},},
		"pw": [][2]int{{-1,0}, {-2, 0}, {-1,1},{-1,-1},},
		"r": [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}},
		"b": [][2]int{{1, 1}, {-1, -1}, {-1, 1}, {1, -1}},
		"q": [][2]int{{1, 0}, {1, 1}, {-1, 0}, {-1, -1}, {-1, 1}, {1, -1}, {0, 1}, {0, -1}},
		"n": [][2]int{{1,2}, {2,1}, {-1,2}, {-2,1}, {1,-2}, {2,-1}, {-1,-2}, {-2,-1}},
		"k": [][2]int{{1, 0}, {1, 1}, {-1, 0}, {-1, -1}, {-1, 1}, {1, -1}, {0, 1}, {0, -1}, {0, 2}, {0, -2}},
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

	castlingStatus := getCastlingStatus(cb, fS)
	fS.EnPass = "e3"

	enPassSqr := cb.board[5][4]

	for row:= range len(cb.board){
		for col:= range len(cb.board[0]){
			cbSqr := cb.board[row][col]
			cbPiece := cbSqr.Piece
			pieceSymbol := cbPiece.Symbol
			pieceColour := cbPiece.Colour
			pieceSymbolLower := strings.ToLower(pieceSymbol)

			if len(pieceSymbol)==0 || pieceColour != fS.Turn{
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
							if pieceSymbolLower=="p"{
								if posRow-row != move[0] || posCol-col != move[1]{
									log.Println("HERE")
								}
								canMovePawn := canMovePawn(pawnSymbol, move, row, col, cb, enPassSqr, posRow, posCol)
								if canMovePawn{
									posSqrs = append(posSqrs, posSqr)
								}
							}else if pieceSymbolLower=="k"{
								if (math.Abs(float64(move[1]))>1){
									canCastle := canCastle(pieceSymbol, move, castlingStatus)		
									if canCastle{
										posSqrs = append(posSqrs, posSqr)
									}
								}else{
									posSqrs = append(posSqrs, posSqr)
								} 
							}else{
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