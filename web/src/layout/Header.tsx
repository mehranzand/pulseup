import { Layout as AntLayout, Space } from "antd";
import  { DockerOutlined } from '@ant-design/icons';

function Header() {
  return (
    <AntLayout.Header style={{backgroundColor: "#099cec", height: 40, color: "#fff", fontSize: 23, lineHeight:"35px", fontFamily:"monospace", fontWeight: 700}}>
      <Space>
        <DockerOutlined />
        <span>pulseUp</span>
      </Space>
    </AntLayout.Header>
  );
}

export default Header;
