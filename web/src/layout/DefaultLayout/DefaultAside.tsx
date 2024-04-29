
import { Layout } from "antd";
import ContainerList from "../../components/ContainerList";

const { Sider } = Layout;

function DefaultAside() {
  return (
    <Sider
      breakpoint="lg"
      collapsedWidth="60"
      width={255}
      style={{ backgroundColor: "#141b1f", minHeight: 'calc(100vh - 40px)' }}
    >
      <ContainerList></ContainerList>
    </Sider>
  );
}

export default DefaultAside;