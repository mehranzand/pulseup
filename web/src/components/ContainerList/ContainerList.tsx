import { useState, useEffect } from 'react'
import {Link} from 'react-router-dom';
import './containerList.css'

interface ContainerListProps {
  host?: string
}

function ContainerList(props: ContainerListProps) {
  const [containers, setContainers] = useState<any[]>([])
  useEffect(() => {

    const fetchData = async () => {
      const res = await fetch(`http://localhost:7070/api/${props.host ? props.host : 'localhost'}/containers`)
      const data = await res.text()

      var json = JSON.parse(data)
      setContainers(json)
    }

    fetchData().catch((err) => console.log(err))
  }, [])
  return (<>
    <h1 style={{marginLeft: 20}}>{props.host ? props.host : 'localhost'}</h1>
    <ul style={{ listStyleType: 'none' }}>
      {containers?.map((c, i) => (
        <Link to={'/container/' + c.id}>
          <li key={i} >
            <span style={{ cursor: "pointer", margin: 5, fontFamily: "monospace", fontSize: 16, color: "#fff", fontWeight: 'bold' }} onClick={() => { }}> {c.name}</span>
          </li>
        </Link>
      ))}
    </ul>
  </>
  )
}

export default ContainerList