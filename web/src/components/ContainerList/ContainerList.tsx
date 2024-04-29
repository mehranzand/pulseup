import { useEffect } from 'react'
import { Row, Col, Typography } from 'antd'
import { Link } from 'react-router-dom';
import { useAppDispatch, useAppSelector } from "../../hooks";
import { fetchContainers } from '../../stores/slices/containerSlice';
import ContainerRow from '../ContainerRow';
import DotsLoader from '../DotsLoader';
import './containerList.css'

interface ContainerListProps {
  host?: string
}

const { Title } = Typography

function ContainerList(props: ContainerListProps) {
  const dispatch = useAppDispatch();
  const { data, loading, error } = useAppSelector(state => state.containers)

  useEffect(() => {
    let title = 'pulseUp for Docker'
    const running = data.filter(d => d.state == 'running')

    if (running.length > 0) {
      title = `${running.length} containers | pulseUp`
    }

    document.title = title
  }, [data])

  useEffect(() => {
    dispatch(fetchContainers(props.host ? props.host : 'localhost'))
  }, [dispatch]);

  return (<>
    <Row>
      <Col span={24}>
        {loading ? <DotsLoader style={{ paddingTop: 40, paddingBottom: 20 }} /> :
          <Title className='host-title' level={4}>{props.host ? props.host : 'localhost'}</Title>}
      </Col>
    </Row>
    <Row>
      <Col>
        <ul className='container-ul'>
          {!loading && data?.map((c, i) => (
            <Link key={i} to={'/container/' + c.id} state={c}>
              <ContainerRow continer={c}></ContainerRow>
            </Link>
          ))}

          {!loading && error && <li>{error}</li>}
        </ul>
      </Col>
    </Row>
  </>
  )
}

export default ContainerList