import React from "react";
import { Layout as AntLayout } from "antd";
import Header from "./Header";
import Aside from "./Aside";

interface LayoutProps {
  noHeader?: boolean
  children: React.ReactNode
}
function Layout(props: LayoutProps) {
  return (
    <AntLayout>
      {!props.noHeader && <Header />}
      <AntLayout >
        <Aside />
        <AntLayout.Content style={{ backgroundColor: "#1c262d"}}>{props.children}</AntLayout.Content>
      </AntLayout>
    </AntLayout>
  );
}

export default Layout;