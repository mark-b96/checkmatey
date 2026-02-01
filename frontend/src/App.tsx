import './App.css';
import { useState } from 'react'
import Chessboard from './components/Chessboard/Chessboard';
import SubmitButton from './components/SubmitButton/submitButton';

function App() {
  const [fenRep, setFenRep] = useState("rnbqkbnr/pppppppp/8/2k3k1/2n6/6n1/PP1PPP1P/RNBQKBNR w KQkq - 0 1")
  const [legalMoves, setLegalMoves] = useState({});

  const updateLegalMoves = (moves: Record<any, any>) => {
    setLegalMoves(moves)
  };

  const updateFenRep = (fenRepStr: string) => {
    setFenRep(fenRepStr)
  };


  return (
    <div id="app">
      <SubmitButton updateMoves={updateLegalMoves} fenRep={fenRep}></SubmitButton>
      <Chessboard fenRep={fenRep} legalMoves={legalMoves} updateFenRep={updateFenRep}/>
    </div>
  );
}

export default App;
