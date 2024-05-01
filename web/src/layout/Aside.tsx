import { Layout as AntLayout } from "antd";
import ContainerList from "../components/ContainerList";

function Aside() {
  return (
    <AntLayout.Sider
      breakpoint="lg"
      collapsedWidth="60"
      width={255}
      style={{ backgroundColor: "#141b1f", minHeight: 'calc(100vh - 40px)' }}
    >
      <ContainerList ></ContainerList>
    </AntLayout.Sider>
  );
}

export default Aside;