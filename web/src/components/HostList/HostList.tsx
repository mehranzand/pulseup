
import './hostList.css'

interface HostListProps {
  host?: string
}

function HostList(props: HostListProps) {

  return (
    <>{props.host}</>
  )
}

export default HostList