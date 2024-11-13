import { Row, Col, Badge, Space } from 'antd';
import { Container } from '../../types/Container';
import './containerRow.css'

interface ContainerRowProps {
  continer: Container
}

function ContainerRow(props: ContainerRowProps) {
  const badge = () => {
    if (props.continer.state == 'running') {
      return <Badge status="success"></Badge>
    } else if (props.continer.state == 'exited') {
      return <Badge status="warning"></Badge>
    } else {
      return <Badge status="processing"></Badge>
    }
  }

  return (
    <>
      <Row className='continer-row'>
        <Col>
          <Space >
            {badge()}
            <div className='continer-name' style={{color: props.continer.state == 'exited'? 'gray' : 'white'}}>{props.continer.name}</div>
          </Space>
        </Col>
      </Row>
    </>
  )
}

export default ContainerRow