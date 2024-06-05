
import { Row, Col, Typography, Space, Tag, Button } from 'antd'
import { MoreOutlined} from '@ant-design/icons';
import { useAppSelector } from "../../hooks"
import { Container } from '../../types/Container';
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
    <Row className='container-bar' align='middle'>
      <Col span={18}>
        {container && <Space align='baseline'>
          <span className='container-name' >
            {container?.name}
          </span>
          <Text className='continer-tag' keyboard>{container?.image}</Text>
          <Text className='continer-tag' keyboard>{container?.status}</Text>
          <Tag color={tagColor(container)}>{container?.state}</Tag>

        </Space>}
      </Col>
      <Col span={2} offset={4} style={{textAlign:"right"}}>
      <Button  type="primary" size="small" className='more-button'><MoreOutlined style={{fontSize:20}}/></Button>
      </Col>
    </Row>
  )
}

export default ContinerInfoBar