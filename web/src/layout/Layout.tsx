import React, { useState, useEffect } from "react";
import { Layout as AntLayout, Affix, Button, Grid } from "antd";
import {
  MenuFoldOutlined,
  MenuUnfoldOutlined
} from '@ant-design/icons';
import Header from "./Header";
import Aside from "./Aside";

const { useBreakpoint } = Grid;

interface LayoutProps {
  noHeader?: boolean
  children: React.ReactNode
}

function Layout(props: LayoutProps) {
  const [collapsed, setCollapsed] = useState(false);
  const screens = useBreakpoint();

  useEffect(() => {
    if (screens.lg === undefined) return;

    if (screens.lg) {
      setCollapsed(false);
    } else {
      setCollapsed(true);
    }
  }, [screens.lg]);

  const toggleCollapsed = () => {
      setCollapsed(!collapsed);
  };

  return (
    <AntLayout>
      {!props.noHeader && 
      <Affix offsetTop={0}>
        <Header />
        <Button
            type="text"
            icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
            onClick={() => toggleCollapsed()}
            style={{
              fontSize: '16px',
              width: 40,
              height: 40,
              color: "white",
              backgroundColor: "#099cec",
              position:"absolute",
              left: 0,
              top: 0
            }}
          />
      </Affix>}
      <AntLayout >
        <Aside collapsed={collapsed} />
        <AntLayout.Content style={{ backgroundColor: "#1c262d"}}>{props.children}</AntLayout.Content>
      </AntLayout>
    </AntLayout>
  );
}

export default Layout;