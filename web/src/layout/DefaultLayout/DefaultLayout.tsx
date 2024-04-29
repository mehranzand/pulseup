import React from "react";
import { Layout } from "antd";
import Sidebar from "../DefaultLayout/DefaultAside";
import Header from "../DefaultLayout/DefaultHeader";
import Footer from "../DefaultLayout/DefaultFooter";

const { Content } = Layout;

interface LayoutProps {
  noHeader?: boolean
  noFooter?: boolean
  children: React.ReactNode
}

function DefaultLayout(props: LayoutProps) {
  return (
    <Layout>
      {!props.noHeader && <Header />}
      <Layout >
        <Sidebar />
        <Content style={{ backgroundColor: "#1c262d"}}>{props.children}</Content>
        {!props.noFooter && <Footer />}
      </Layout>
    </Layout>
  );
}

export default DefaultLayout;