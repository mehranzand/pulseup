import { Layout, Space } from "antd";
import  { DockerOutlined } from '@ant-design/icons';

const { Header } = Layout;

function DefaultHeader() {
  return (
    <Header style={{backgroundColor: "#099cec", height: 40, color: "#fff", fontSize: 23, lineHeight:"35px", fontFamily:"monospace", fontWeight: 700}}>
      <Space>
        <DockerOutlined />
        <span>pulseUp</span>
      </Space>
    </Header>
  );
}

export default DefaultHeader;
