import { Layout as AntLayout, Row, Col, Typography, Button, Space, Tooltip } from "antd";
import { MenuOutlined } from '@ant-design/icons';
import ContainerList from "../components/ContainerList";
import HostList from "../components/HostList";
import config from "../stores/config";
import { useState } from "react";
import './aside.css';


const { Title } = Typography

function Aside() {
  let [viewMode, setViewMode] = useState<"host" | "container">("container")
  let [selectedHost, setSelectedHost] = useState("localhost")

  const handleClickHostList = () => {
    if (viewMode == 'container') {
      setViewMode('host')
    }
  }

  const handleSelectHostCallback=(host : any) =>{
    setSelectedHost(host.name)
    setViewMode("container")
  }

  return (
    <AntLayout.Sider
      width={255}
      className="aside"
    >
      {viewMode == "container" && <><Row>
        <Col span={24}>
          <Space align="baseline">
            <Tooltip title={selectedHost} >
              <Title level={4} className="host-title">{selectedHost}</Title>
            </Tooltip>
            {config.hosts.length > 1 && <Tooltip title="host list" placement="top">
              <Button type="primary" size="small" className="back-button" icon={<MenuOutlined />} onClick={handleClickHostList} />
            </Tooltip>}
          </Space>
        </Col>
      </Row>
        <ContainerList host={selectedHost} ></ContainerList></>}
      {viewMode == "host" && <HostList hosts={config.hosts} onSelect={handleSelectHostCallback}></HostList>}

    </AntLayout.Sider>
  );
}

export default Aside;