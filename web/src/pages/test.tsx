import { useState, useEffect } from 'react'
import { Typography, Row, Col, Tag } from "antd";
import './test.css'

const { Text } = Typography;

function Test() {
  const [data, setData] = useState<any[]>([])
  const [logs, setLogs] = useState<any[]>([])
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

    setLogs([])

    es = new EventSource("http://localhost:7070/api/logs/stream/localhost/" + current.id)
    es.onmessage = (e) => {
      if (e.data) {
        var m = JSON.parse(e.data)

        setLogs(oldLogs => [m, ...oldLogs])
      }
    };
    es.onerror = () => setLogs([]);

    es.addEventListener('close', () => close()
  )
    return (() => close() )
  }, [current]);

  return (
    <div style={{height: "100%", width: "100%", textAlign: "left", padding: 40}}>
      {data.length > 0 && <h2 style={{color: 'white'}}>container list</h2>}
      <ul style={{listStyleType: 'none'}}>
        {data?.map((c, i)=>(
          <li key={i}>
           <Tag color='#108ee9'  style={{cursor: "pointer", margin: 5}}  onClick={() => {setCurrent(c)}}> {c.name}</Tag>  
          </li>
          ))}
      </ul>
      <br />
      {current && <h2 style={{color: 'white'}}>streaming log for {current?.name}</h2>}
      <h3><span style={{whiteSpace: "pre-line", textAlign: "left", color: 'white'}}>{logs.map(item => {
       return <Row style={{marginBottom: 5}}>
                <Col span={2}>  <Tag>{item.t}</Tag></Col>
                <Col span={22}>  <Text type='success'>{item.m}</Text></Col>
              </Row>
      })}</span></h3>
    </div>
  )
}

export default Test