
import { Row, Col, Typography, Space, Tag } from 'antd'
import { useAppSelector } from "../../hooks"
import { Container } from '../../types/Container';
import './continerInfoBar.css'

interface ContinerInfoBarProps {
  continerId: string | undefined
}

const { Text } = Typography;

function ContinerInfoBar(props: ContinerInfoBarProps) {
  const continer = useAppSelector((state) => state.containers.data.find(a => a.id == props.continerId))

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
    <Row className='container-bar' align='middle'>
      <Col>
        {continer && <Space align='baseline'>
          <span className='container-name' >
            {continer?.name}
          </span>
          <Text className='continer-tag' keyboard>{continer?.image}</Text>
          <Text className='continer-tag' keyboard>{continer?.status}</Text>
          <Tag color={tagColor(continer)}>{continer?.state}</Tag>
        </Space>}
      </Col>
    </Row>
  )
}

export default ContinerInfoBar