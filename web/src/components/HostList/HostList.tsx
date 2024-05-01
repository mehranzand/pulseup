
import { Row, Col } from "antd";
import './hostList.css'

interface HostListProps {
  hosts: { name: string; id: string }[];
}

function HostList(props: HostListProps) {
  return (
    <Row>
      <Col span={24}>
        {props.hosts?.length}
      {props.hosts?.map(host => 
        <a>{host.name}</a>
      )}
      </Col>
    </Row>
  )
}

export default HostList