import { Layout as AntLayout, Row, Col, Typography, Button, Space, Tooltip } from "antd";
import { LeftOutlined } from '@ant-design/icons';
import ContainerList from "../components/ContainerList";
import HostList from "../components/HostList";
import config from "../stores/config";
import './aside.css'
import { useState } from "react";

const { Title } = Typography

function Aside() {
  let [viewMode] = useState<"host" | "container">("container")
  let [selectedHost] = useState("localhost")

  return (
    <AntLayout.Sider
      breakpoint="lg"
      collapsedWidth="60"
      width={255}
      style={{ backgroundColor: "#141b1f", minHeight: 'calc(100vh - 40px)' }}
    >
      <Row>
        <Col span={24}>
          <Space align="baseline" style={{marginLeft: 10}}>
           {config.hosts.length > 1 && <Tooltip title="go to hosts list">
              <Button type="primary"size="small" shape="circle" className="back-button" icon={<LeftOutlined />} />
            </Tooltip>}
            <Title level={4} style={{color: '#fff'}}>{config.hosts[1].name}</Title>
          </Space>
        </Col>
      </Row>
      {viewMode == "host" && <HostList hosts={config.hosts}></HostList>}
      {viewMode == "container" && <ContainerList host={selectedHost} ></ContainerList>}
    </AntLayout.Sider>
  );
}

export default Aside;