import { Row, Col, Typography, Space, Button } from 'antd'
import { useAppSelector } from "../../hooks"
import { Container } from '../../types/Container';
import { MenuOutlined } from '@ant-design/icons';
import './continerInfoBar.css'

interface ContinerInfoBarProps {
  continerId: string | undefined
}

const { Text } = Typography;

function ContinerInfoBar(props: ContinerInfoBarProps) {
  const container = useAppSelector((state) => state.containers.data.find(a => a.id == props.continerId))

  const getType = (c: Container | undefined) => {
    if (c == undefined) return

    if (c.state == 'running') {
      return "success"
    } else if (c.state == 'exited') {
      return "danger"
    }
    else {
      return "secondary"
    }
  }

  return (
    <>
      <Row className='container-bar' align='middle'>
        <Col flex="auto">
          {container && <Space align='baseline'>
            <span className='container-name' >
              {container?.name}
            </span>
            <Text className='container-tag' keyboard>{container?.image}</Text>
            <Text className='container-tag' keyboard>{container?.status}</Text>
            <Text  type={getType(container)} >{container?.state}</Text>
          </Space>}
        </Col>
        <Col>
          <Button type="text" size="small" style={{color: "white",  marginRight: "10px" }} icon={<MenuOutlined />}></Button>
        </Col>
      </Row>
    </>
  )
}

export default ContinerInfoBar