import { useParams, useLocation} from 'react-router-dom';
import ContinerInfoBar from '../components/ContinerInfoBar';
import DotsLoader from '../components/DotsLoader';
import useLogSourceEvent from '../lib/hooks/useLogSourceEvent';
import { Row, Col, Tag, Tooltip, Typography } from 'antd'
import { useEffect, useState, CSSProperties, createRef, forwardRef } from 'react'
import { FixedSizeList as List } from 'react-window';
import AutoSizer from "react-virtualized-auto-sizer";
import Ansi from 'ansi-to-react';
import moment from 'moment';
import './logViewer.css'

const { Text } = Typography
const PADDING_SIZE = 20;
const ITEM_SIZE = 32;

function LogViewer() {
  let params = useParams()
  let location = useLocation()
  const { messages, loading, pause, resume } = useLogSourceEvent()
  const listRef = createRef<HTMLDivElement>();
  const [paused, setPaused] = useState<boolean>(false)

  function scrollToEnd() {
   listRef.current?.scrollIntoView({ behavior: 'instant', block: 'end' });
  }

  useEffect(() => {
    setPaused(resume())
    scrollToEnd()
  }, [location, !loading])

  useEffect(() => {
    if (!paused) {
     scrollToEnd()
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

    return <Row className='log-row'
      style={{
        ...style,
        bottom: `${parseFloat(Number(style.bottom).toString()) + PADDING_SIZE}px`
      }}>
      <Col span={1.5}>
        <Tooltip title={record.type} placement='top'>
          <span className={`type-dot ${record.type}`}></span>
        </Tooltip>
      </Col>
      <Col span={4.5}>
        <Tag color='blue-inverse'>{moment(record.date).format('MM/DD/YYYY h:mm:ss a')}</Tag></Col>
      <Col span={18}>
        <Text className='log-message'><Ansi>{record.message}</Ansi></Text>
      </Col>
    </Row>
  };

  const innerElementType = forwardRef<HTMLDivElement, React.HTMLAttributes<HTMLDivElement>>(
    ({ style, ...rest }, ref) => (
      <div
        ref={ref}
        style={{
          ...style,
          height: `${parseFloat(Number(style?.height).toString()) + PADDING_SIZE * 2}px`
        }}
        {...rest}
      />
    )
  );

  return (
    <>
      {!loading && <ContinerInfoBar continerId={params.id} ></ContinerInfoBar>}
      <AutoSizer>
        {({ height, width }) => (
          <List
            innerRef={listRef}
            className='log-wrapper'
            height={height - 50}
            itemCount={messages.length}
            innerElementType={innerElementType}
            itemSize={ITEM_SIZE}
            width={width}
          >
            {LogRow}
          </List>
        )}
      </AutoSizer>
      {loading && <DotsLoader style={{ margin: '45vh' }} />}
      {!loading && messages.length > 0 && <button className={'scroll-down' + (paused ? ' bounce' : '')} onClick={scrollToEnd}></button>}
    </>
  )
}

export default LogViewer