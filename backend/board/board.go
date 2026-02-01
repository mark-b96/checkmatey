package board


type Piece struct {
	Row int
	Col int
	Colour int
	Symbol string
	Promoted bool
}

type square struct {
	Row int
	Col int
	Colour int
	Piece Piece
}

var pieces []Piece

type chessboard struct{
	cb [8][8]square
}


type AllPieces struct{
	Pieces []Piece
}

func GetBoard() AllPieces {
	var width = 8
	var height = 8
	var board [8][8]square

	for col:= range height{
		for row:= range width{

			newPiece := Piece{}
			switch {
				case (row==0 || row==7) && (col==2 || col==5):
					newPiece = Piece{Symbol: "B"}
				case (row==0 || row==7) && (col==1 || col==6):
					newPiece = Piece{Symbol: "N"}
				case (row==0 || row==7) && (col==0 || col==7):
					newPiece = Piece{Symbol: "R"}
				case (row==0 || row==7) && col==3:
					newPiece = Piece{Symbol: "Q"}
				case (row==0 || row==7) && col==4:
					newPiece = Piece{Symbol: "K"}
				case row==1:
					newPiece = Piece{Symbol: "P"}
				case row==6:
					newPiece = Piece{Symbol: "P"}
			}

			if row < 4{
				newPiece.Colour = 0
			}else {
				newPiece.Colour = 1
			}

			newPiece.Promoted = false
			newPiece.Col = col
			newPiece.Row = row

			pieces = append(pieces, newPiece)

			if (row+col)%2 ==0{
				pBlack := square{Row: row, Col: col, Colour: 0, Piece: newPiece}
				board[row][col] = pBlack
			}else {
				pWhite := square{Row: row, Col: col, Colour: 1, Piece: newPiece}
				board[row][col] = pWhite
			}
		}
	}
	
	newPieces := AllPieces{Pieces: pieces}
	return newPieces
}

func GetPieces() []Piece{
	return pieces
}




