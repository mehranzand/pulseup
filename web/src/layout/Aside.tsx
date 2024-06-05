import { Layout as AntLayout, Row, Col, Typography, Button, Space, Tooltip } from "antd";
import { MenuOutlined } from '@ant-design/icons';
import ContainerList from "../components/ContainerList";
import HostList from "../components/HostList";
import config from "../stores/config";
import { useState } from "react";
import { useAppSelector, useAppDispatch } from "../hooks"
import { setCurrent } from "../stores/slices/hostSlice"
import './aside.css';

const { Title } = Typography

function Aside() {
  let [viewMode, setViewMode] = useState<"host" | "container">("container")
  const { current } = useAppSelector((state) => state.host)
  const dispatch = useAppDispatch()

  const handleClickHostList = () => {
    if (viewMode == 'container') {
      setViewMode('host')
    }
  }

  const handleSelectHostCallback = (host: string) => {
    setViewMode("container")

    dispatch(setCurrent(host))
  }

  return (
    <AntLayout.Sider
      width={255}
      className="aside"
      breakpoint={"lg"}
      theme="light"
      collapsedWidth={0}
      trigger={null}
    >
      {viewMode == "container" && <><Row>
        <Col span={24}>
          <Space align="baseline">
            <Tooltip title={current} >
              <Title level={4} className="host-title">{current}</Title>
            </Tooltip>
            {config.hosts.length > 1 && <Tooltip title="host list" placement="top">
              <Button type="primary" size="small" className="back-button" icon={<MenuOutlined />} onClick={handleClickHostList} />
            </Tooltip>}
          </Space>
        </Col>
      </Row>
        <ContainerList></ContainerList></>}
      {viewMode == "host" && <HostList hosts={config.hosts} onSelect={handleSelectHostCallback}></HostList>} 
    </AntLayout.Sider>
  );
}

export default Aside;