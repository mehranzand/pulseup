import { Layout as AntLayout, Space } from "antd";
import { DockerOutlined } from '@ant-design/icons';
import './header.css';

function Header() {
  return (
    <AntLayout.Header className="header" >
      <Space>
        <DockerOutlined />
        <span className="logo-type">pulseUp</span>
      </Space>
    </AntLayout.Header>
  );
}

export default Header;
