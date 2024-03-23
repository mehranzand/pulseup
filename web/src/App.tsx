import { useState, useEffect } from 'react'

import './App.css'

function App() {
  const [data, setData] = useState<any[]>([]);
  const [data1, setData1] = useState("")
  const [current, setCurrent] = useState<any>(null)
  let es: EventSource | null = null;

  function close() {
    if (es) {
      es.close();
      console.debug(`EventSource closed for ${current.id}`);
      es = null;
    }
  }

 
  useEffect(() =>{
    const fetchData = async () => {
      const res = await fetch("http://localhost:7070/api/localhost/containers")
      const data = await res.text()

      var json  = JSON.parse(data)
      setData(json)
    }

    fetchData().catch((err) => console.log(err))
  }, [])

  useEffect(() => {
    if (current == null) return;

    setData1("")

    es = new EventSource("http://localhost:7070/api/logs/stream/localhost/" + current.id)
    es.onmessage = (e) => {
      if (e.data) {
        setData1((prev) => e.data + '\n' + prev)
      }
    };
    es.onerror = () => setData1("");

    es.addEventListener('close', () => close()
  )
    return (() => close() )
  }, [current]);


  return (
    <div style={{height: "100vh", width: "100vw", textAlign: "left"}}>
      <h1 style={{color: 'green'}}>pulseUp</h1>
      <h2 style={{color: 'GrayText'}}>Seamless log monitoring for Docker containers with intelligent</h2>  
      <h2 style={{color: 'GrayText', marginTop: -15}}>action logs for next-level performance and insight.</h2> 
      <br />
      {data.length > 0 && <h2>containers list</h2>}
      <ul>
        {data?.map((c, i)=>(
          <li key={i}>
            <a style={{cursor: "pointer"}} onClick={() => {setCurrent(c)}}> {c.name}</a>   
          </li>
          ))}
      </ul>
      <br />
      {current && <h2>streaming log for {current?.name}</h2>}
      <h3><span style={{whiteSpace: "pre-line", textAlign: "left"}}>{data1}</span></h3>
    </div>
  )
}

export default App