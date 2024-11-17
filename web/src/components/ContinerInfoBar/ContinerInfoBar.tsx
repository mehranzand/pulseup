import { Row, Col, Typography, Space, Tag, Button } from 'antd'
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

  const tagColor = (c: Container | undefined) => {
    if (c == undefined) return

    if (c.state == 'running') {
      return "#55acee"
    } else if (c.state == 'exited') {
      return "#f50"
    }
    else {
      return "#108ee9"
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
            <Tag className='state' color={tagColor(container)}>{container?.state}</Tag>

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