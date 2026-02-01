interface Props {
    updateMoves: (moves: Record<any, any>) => void 
    fenRep: string
}

export default function SubmitButton({updateMoves, fenRep}: Props) {
    function submitRequest(): void {
        const fetchData = async () => {
            try{
                console.log("FENREP:", fenRep)
                const response = await fetch(`http://localhost:5669/getMoves?fenRep=${fenRep.replaceAll(" ", "+")}`, {
                    method: 'GET',    
                })
                const result = await response.json()          
                const resultMap: Record<any, any> = {}
                    
                for (const [key, value] of Object.entries(result)) {
                    resultMap[key] = value
                  }
                console.log(resultMap)

                updateMoves(resultMap)
            }
            catch (err) {
                console.error(err)
            }
        }
        fetchData()
    }
    return (
    <div> 
        <button onClick={submitRequest}>
            Submit
        </button>
    </div>
    )
}