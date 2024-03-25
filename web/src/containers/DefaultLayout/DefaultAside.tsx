
import { Layout } from "antd";

const { Sider } = Layout;

function DefaultAside() {
  return (
    <Sider
      breakpoint="lg"
      collapsedWidth="60"
      width={255}
      style={{backgroundColor: "#141b1f"}}
    >

    </Sider>
  );
}

export default DefaultAside;