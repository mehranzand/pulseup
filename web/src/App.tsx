import { useState, useEffect } from 'react'

import './App.css'

function App() {
  const [data, setData] = useState("")

  useEffect(() =>{
    const fetchData = async () => {
      const res = await fetch("http://localhost:7070/api")
      const data = await res.text()

      setData(data)
    }

    fetchData().catch((err) => console.log(err))
  }, [])

  return (
    <>
    FDHFDHFGHFGG
      {data}
      <h1>Hello Air!!!!!!!!!!!</h1>  
    </>
  )
}

export default App