import './App.css';
import { useState, useEffect } from 'react'
import Chessboard from './components/Chessboard/Chessboard';

function App() {
  const [fenRep, setFenRep] = useState("")
  const [legalMoves, setLegalMoves] = useState({});

  useEffect(() => {
    const fetchData = async () => {
      try {
        const response1 = await fetch(`http://localhost:5669/getInitState`, {
          method: 'GET',    
      })
      const fenRepResult = await response1.json()
      setFenRep(fenRepResult)
      
      const response2 = await fetch(`http://localhost:5669/getMoves?userMove=`, {
        method: 'GET',    
      })
      const result = await response2.json()          
      const resultMap: Record<any, any> = {}
          
      for (const [key, value] of Object.entries(result)) {
          resultMap[key] = value
        }
      updateLegalMoves(resultMap)
    }
    catch (err) {
        console.error(err)
    }
    };
  
    fetchData();
  }, []);

  const updateLegalMoves = (moves: Record<any, any>) => {
    setLegalMoves(moves)
  };

  const updateFenRep = (userMove: string) => {

    const fetchData = async () => {
      try{
          const response = await fetch(`http://localhost:5669/getMoves?userMove=${userMove}`, {
              method: 'GET',    
          })
          const result = await response.json()          
          const resultMap: Record<any, any> = {}
              
          for (const [key, value] of Object.entries(result)) {
              resultMap[key] = value
            }
          console.log("Successfully fetched moves...")

          updateLegalMoves(resultMap)
      }
      catch (err) {
          console.error(err)
      }
    }
    fetchData()
  };


  if (!fenRep)
  {
    return <div>Loading...</div>;
  }
  else{
    return (
      <div id="app">
        <Chessboard fenRep={fenRep} legalMoves={legalMoves} updateFenRep={updateFenRep}/>
      </div>
    );
  }
}

export default App;
