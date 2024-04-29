
import { BrowserRouter } from "react-router-dom";
import { AppRoutes } from "./routes/";
import useContinerSourceEvent from "./lib/hooks/useContinerSourceEvent";
import './App.css'

function App() {
  useContinerSourceEvent();

  return (
    <BrowserRouter children={AppRoutes} basename={"/"} />
  )
}

export default App