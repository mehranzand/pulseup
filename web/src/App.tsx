import { useState, useEffect } from 'react'

import './App.css'

function App() {
  const [data, setData] = useState("")
  const [data1, setData1] = useState("")

  useEffect(() =>{
    const fetchData = async () => {
      const res = await fetch("http://localhost:7070/base/api/localhost/containers")
      const data = await res.text()

      setData(data)
    }

    fetchData().catch((err) => console.log(err))
  }, [])

  useEffect(() =>{
    const fetchData = async () => {
      const res = await fetch("http://localhost:7070/base/api/logs/:host/:id")
      const data = await res.text()

      setData1(data)
    }

    fetchData().catch((err) => console.log(err))
  }, [])

  return (
    <>
      <h1>{data}</h1>
      <h1>{data1}</h1>
      <h2>Seamless log monitoring for Docker containers with intelligent </h2>  
      <h2> action logs for next-level performance and insight.</h2>  
    </>
  )
}

export default App