import './Chessboard.css'
import React, { useState } from 'react'

interface Props {
    fenRep: string
    legalMoves: Record<string, number[][]>
    updateFenRep: (userMove: string) => void 
}

const pieceMap: Record<string, string> = {
    "p": "assets/black-pawn.png",
    "n": "assets/black-knight.png",
    "b": "assets/black-bishop.png",
    "r": "assets/black-rook.png",
    "q": "assets/black-queen.png",
    "k": "assets/black-king.png",
    "P": "assets/white-pawn.png",
    "N": "assets/white-knight.png",
    "B": "assets/white-bishop.png",
    "R": "assets/white-rook.png",
    "Q": "assets/white-queen.png",
    "K": "assets/white-king.png"
}

export default function Chessboard({fenRep, legalMoves, updateFenRep}: Props) {
    let prevMoves: number[][] = [[]]
    let prevSquareId: string = ""
    let board: React.ReactElement[] = []

    var fenRepVisTmp = replaceDigits(fenRep)

    function replaceDigits(originalString: string): string {

        const digitWords = {
          '1': '1', '2': '22', '3': '333', '4': '4444',
          '5': '55555', '6': '666666', '7': '7777777', '8': '88888888'
        };
      
        const newString = originalString.replace(/\d/g, (match) => {
            let result = digitWords[match as keyof typeof digitWords];
            
            if (result === undefined){
                return match 
            }
            return result
          
        });
        return newString
    }

    const [fenRepVis, setFenRepVis] = useState(fenRepVisTmp)
   
    function updateBoard() {
        board = []
        let fenRepList: string[] = fenRepVis.split("/")
        for (let row = 0; row < 8; row++){
            for (let col = 0; col < 8; col++){
                let fenRow: string = fenRepList[row]
                let square: string = fenRow[col]
                let img = undefined

                if (typeof square !== 'undefined') {
                    img = pieceMap[square as keyof typeof pieceMap]
                } 
            
                if ((row+col) %2===0){
                    board.push(
                        <div className="square white-square" key={`${row},${col}`} id={`${row},${col}`}>
                            {img && <div style={{ backgroundImage: `url(${img})` }} className="chess-piece"  id={`${row},${col}`}></div>}
                        </div>
                    )
                }else {
                    board.push(
                        <div className="square black-square" key={`${row},${col}`} id={`${row},${col}`}>
                        {img && <div style={{ backgroundImage: `url(${img})` }} className="chess-piece"  id={`${row},${col}`}></div>}
                        </div>
                    )
                } 
            }
        }
        return board
    }

    function updateFenState(srcSquare: string, dstSquare: string) {
        let fenRepList: string[] = fenRepVis.split("/")

        let updatedFenRepVis: string = ""
        const [srcRow, srcCol] = srcSquare.split(",")
        var srcPiece: string = fenRepList[parseInt(srcRow)][parseInt(srcCol)]


        var userMove = `${srcPiece}:${srcSquare}-${dstSquare}`
        updateFenRep(userMove)
     
        for (let row = 0; row < 8; row++){
            let emptySquares: number = 0
            for (let col = 0; col < 8; col++){
                let fenRow: string = fenRepList[row]
                let fenSqr: string = fenRow[col]
                let currSqr: string = `${row},${col}`
                
                if (fenSqr && pieceMap.hasOwnProperty(fenSqr) && currSqr!==dstSquare){   
                    if (currSqr===srcSquare){
                        emptySquares += 1
                        if (emptySquares > 0){
                            updatedFenRepVis += emptySquares.toString().repeat(emptySquares)
                            emptySquares = 0
                        }  
                        srcPiece = fenSqr
                    }
                    else {
                        if (emptySquares > 0){
                            updatedFenRepVis += emptySquares.toString().repeat(emptySquares)
                            emptySquares = 0
                        } 
                        updatedFenRepVis+=fenSqr
                    }
                } 
                else if (currSqr===dstSquare){
                    if (emptySquares > 0){
                        updatedFenRepVis += emptySquares.toString().repeat(emptySquares)
                        emptySquares = 0
                    }
                    updatedFenRepVis += srcPiece
                }
                else {
                    emptySquares +=1
                }
            }
            if (emptySquares > 0){
                updatedFenRepVis += emptySquares.toString().repeat(emptySquares)
            }
            updatedFenRepVis += "/"
        }
 
        setFenRepVis(updatedFenRepVis)
    }

    function getLegalMoves(e: React.MouseEvent): void {
        const element: HTMLElement = e.target as HTMLElement;
        if (prevMoves) {
            var currSquare: string = element.id
            for (let i = 0; i < prevMoves.length; i++) {
                if (currSquare === prevMoves[i].toString()) {
                    updateFenState(prevSquareId, currSquare)
                }
            }
        }
        
        let blackSquares: HTMLCollectionOf<Element> = document.getElementsByClassName('square black-square');
        for (let i = 0; i < blackSquares.length; i++) {
            let id: string = blackSquares[i].id
            const originalSquare: HTMLElement|null = document.getElementById(id)
            if (originalSquare) {
                originalSquare.style.backgroundColor = "#779556"
            }
        }

        let whiteSquares: HTMLCollectionOf<Element> = document.getElementsByClassName('square white-square');
        for (let i = 0; i < whiteSquares.length; i++) {
            let id: string = whiteSquares[i].id
            const originalSquare: HTMLElement|null = document.getElementById(id)
            if (originalSquare) {
                originalSquare.style.backgroundColor = "#ebecd0"
            }
        }
        
        e.preventDefault()
        prevMoves = [[]]

        if (element.classList.contains("chess-piece")) {
            let pieceData: string[] = element.id.split(",")
            const [row, col] = pieceData;
     
            var rowCol: string = `${row},${col}`
            var pieceMoves: number[][] = legalMoves[rowCol as keyof typeof legalMoves]

            if (pieceMoves) {
                for (let i = 0; i<pieceMoves.length; i++){
                    var pSquare: string = pieceMoves[i].toString()
                    const square: HTMLElement|null = document.getElementById(pSquare)
                    if (square){
                        square.style.backgroundColor = "lightblue";
                    }
                    prevMoves.push(pieceMoves[i])  
                }
                prevSquareId = element.id
            }
        }
    }

    return <div onMouseDown={(e) => getLegalMoves(e)}id="chessboard">{updateBoard()}</div>
}