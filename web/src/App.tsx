import { BrowserRouter } from "react-router-dom";
import { AppRoutes } from "./routes/";

import './App.css'

function App() {
  return (
    <BrowserRouter children={AppRoutes} basename={"/"} />
  )
}

export default App