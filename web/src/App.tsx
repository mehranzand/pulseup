import { useState, useEffect } from 'react'

import './App.css'

function App() {
  const [data, setData] = useState("")
  const [data1, setData1] = useState("")

  useEffect(() =>{
    const fetchData = async () => {
      const res = await fetch("http://localhost:7070/api/localhost/containers")
      const data = await res.text()

      setData(data)
    }

    fetchData().catch((err) => console.log(err))
  }, [])

  useEffect(() => {
    const sse = new EventSource('http://localhost:7070/api/logs/stream/localhost/5ab',
      { withCredentials: true });
    function getRealtimeData(data1 : any) {
      // process the data here,
      // then pass it to state to be rendered
      setData1(data1)
    }
    sse.onmessage = e => getRealtimeData(JSON.parse(e.data));
    sse.onerror = () => {
      // error log here 
      
      sse.close();
    }
    return () => {
      sse.close();
    };
  }, []);


  return (
    <div style={{height: "100vh"}}>
      <h1 style={{color: 'green'}}>pulseUp</h1>
      <h2 style={{color: 'GrayText'}}>Seamless log monitoring for Docker containers with intelligent</h2>  
      <h2 style={{color: 'GrayText', marginTop: -15}}>action logs for next-level performance and insight.</h2> 
      <br /><br /><br />
      <span style={{color: 'yellow'}}>{data}</span>
      <br /><br /><br />
      <h3>{data1}</h3>
    </div>
  )
}

export default App