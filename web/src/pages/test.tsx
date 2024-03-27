import { useState, useEffect } from 'react'
import { Typography, Row, Col, Tag } from "antd";
import './test.css'
import { useParams } from 'react-router-dom';

const { Text } = Typography;

function Test() {
  const [logs, setLogs] = useState<any[]>([])
  let es: EventSource | null = null;
  let params = useParams()

  function close() {
    if (es) {
      es.close();
      console.debug(`EventSource closed for ${params.id}`);
      es = null;
    }
  }

  useEffect(() => {
    if (params.id == null) return;

    setLogs([])

    es = new EventSource("http://localhost:7070/api/logs/stream/localhost/" + params.id)
    es.onmessage = (e) => {
      if (e.data) {
        var m = JSON.parse(e.data)

        setLogs(oldLogs => [m, ...oldLogs])
      }
    };
    es.onerror = () => setLogs([]);

    es.addEventListener('close', () => close()
    )
    return (() => close())
  }, [params]);

  return (
    <div style={{ height: "100%", width: "100%", textAlign: "left", padding: "20px 40px" }}>
      {params.id && <h2 style={{ color: 'white' }}>streaming log for {params?.id}</h2>}
      <h3><span style={{ whiteSpace: "pre-line", textAlign: "left", color: 'white' }}>{logs.map(item => {
        return <Row style={{ marginBottom: 5 }}>
          <Col span={1.5}>
            <Tag color='volcano-inverse'>{item.t}</Tag></Col>
          <Col span={4.5}>
            <Tag color='geekblue'>{new Date(item.ts).toUTCString()}</Tag></Col>
          <Col span={18}>
            <Text type='success'>{item.m}</Text></Col>
        </Row>
      })}</span></h3>
    </div>
  )
}

export default Test