import { useParams } from 'react-router-dom';
import ContinerInfoBar from '../components/ContinerInfoBar';
import DotsLoader from '../components/DotsLoader';
import useLogSourceEvent from '../lib/hooks/useLogSourceEvent';
import { Row, Col, Tag, Tooltip, Typography } from 'antd'
import { useEffect, useState, CSSProperties, createRef } from 'react'
import { FixedSizeList as List } from 'react-window';
import AutoSizer from "react-virtualized-auto-sizer";
import Ansi from 'ansi-to-react';
import moment from 'moment';
import './logViewer.css'

const { Text } = Typography

function LogViewer() {
  let params = useParams()
  const { messages, loading, pause, resume } = useLogSourceEvent()
  const listRef = createRef<List<any>>();
  const [paused, setPaused] = useState<boolean>(false)

  function handleScroll() {
    listRef.current?.scrollToItem(messages.length)
  }

  useEffect(() => {
    if (!paused) {
      listRef.current?.scrollToItem(messages.length)
    }
  }, [messages])

  const LogRow = ({ index, style }: { index: number; style: CSSProperties; }) => {
    const record = messages[index]

    if (index != messages.length - 1) {
      setPaused(pause())
    }
    else {
      setPaused(resume())
    }

    return <Row style={style} className='log-row'>
      <Col span={1.5}>
        <Tooltip title={record.type} placement='top'>
          <span className={`type-dot ${record.type}`}></span>
        </Tooltip>
      </Col>
      <Col span={4.5}>
        <Tag color='blue-inverse'>{moment(record.date).format('MM/DD/YYYY h:mm:ss a')}</Tag></Col>
      <Col span={18}>
        <Text className='log-message' copyable><Ansi>{record.message}</Ansi></Text>
      </Col>
    </Row>
  };

  return (
    <>
      {!loading && <ContinerInfoBar continerId={params.id} ></ContinerInfoBar>}
      <AutoSizer>
        {({ height, width }) => (
          <List
            ref={listRef}
            className='log-wrapper'
            height={height - 50 - 20}
            itemCount={messages.length}
            itemSize={32}
            width={width}
          >
            {LogRow}
          </List>
        )}
      </AutoSizer>
      {loading && <DotsLoader style={{ margin: '45vh' }} />}
      {!loading && messages.length > 0 && <button className={'scroll-down' + (paused ? ' bounce' : '')} onClick={handleScroll}></button>}
    </>
  )
}

export default LogViewer