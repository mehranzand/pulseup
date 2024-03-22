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

    const url = "http://localhost:7070/api/logs/stream/localhost/89e69b17c2ee"

    if ('EventSource' in window) {
      let source = new EventSource(url, {withCredentials: false})

      var evtSourceErrorHandler = function(event : any){
        var txt;
        switch( event.target.readyState ){
            case EventSource.CONNECTING:
                txt = 'Reconnecting...';
                break;
            case EventSource.CLOSED:
                txt = 'Reinitializing...';
                source = new EventSource("../sse.php");
                source.onerror = evtSourceErrorHandler;
                break;
        }
        console.log(txt);
    }

    source.addEventListener('message', function(e) {     
      
      setData1((prev) => e.data + '\n' + prev)
    }, false);
  }
    
  
  }, []);


  return (
    <div style={{height: "100vh", width: "100vw", textAlign: "left"}}>
      <h1 style={{color: 'green'}}>pulseUp</h1>
      <h2 style={{color: 'GrayText'}}>Seamless log monitoring for Docker containers with intelligent</h2>  
      <h2 style={{color: 'GrayText', marginTop: -15}}>action logs for next-level performance and insight.</h2> 
      <br /><br /><br />
      <span style={{color: 'yellow'}}>{data}</span>
      <br /><br /><br />
      <h3><span style={{whiteSpace: "pre-line", textAlign: "left"}}>{data1}</span></h3>
    </div>
  )
}

export default App