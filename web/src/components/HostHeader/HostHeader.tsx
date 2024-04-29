
import './hostHeader.css'

interface HostHeaderProps {
  host?: string
}

function HostHeader(props: HostHeaderProps) {

  return (
    <>{props.host}</>
  )
}

export default HostHeader