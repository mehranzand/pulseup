
import React from 'react';
import './dotsLoader.css'

interface DotsLoaderProps {
  style?: React.CSSProperties
}

const DotsLoader = (props : DotsLoaderProps) => {
  return (
    <div style={props.style} className="dots-loader">
      <div></div>
      <div></div>
      <div></div>
    </div>
  );
}

export default DotsLoader