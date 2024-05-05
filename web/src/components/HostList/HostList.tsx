
import { Row, Col, Typography, Space } from "antd";
import { HddOutlined } from '@ant-design/icons';
import './hostList.css'

interface HostListProps {
  hosts: { name: string; id: string }[];
  onSelect: (host: any) => void;
}

function HostList(props: HostListProps) {
  return (
    <>
      {props.hosts?.map(host =>
        <Row>
          <Col span={24}>
            <Typography.Title level={5} className="title">
              <Space size={"small"} onClick={() => props.onSelect(host)}>
                <HddOutlined />
                {host.name}
              </Space>
            </Typography.Title>
          </Col>
        </Row>
      )}</>
  )
}

export default HostList