import { useEffect } from 'react'
import { Row, Col } from 'antd'
import { Link } from 'react-router-dom';
import { useAppDispatch, useAppSelector } from "../../hooks";
import { fetchContainers } from '../../stores/slices/containerSlice';
import ContainerRow from '../ContainerRow';
import './containerList.css'


function ContainerList() {
  const dispatch = useAppDispatch();
  const { data, loading, error } = useAppSelector(state => state.containers)
  const { current } = useAppSelector((state) => state.host)

  useEffect(() => {
    let title = 'pulseUp for Docker'
    const running = data.filter(d => d.state == 'running')

    if (running.length > 0) {
      title = `${running.length} containers | pulseUp`
    }

    document.title = title
  }, [data])

  useEffect(() => {
    if (current === undefined) return

    dispatch(fetchContainers(current))
  }, [dispatch, current]);

  return (<>
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