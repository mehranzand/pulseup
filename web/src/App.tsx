
import { BrowserRouter } from "react-router-dom";
import { useEffect } from "react";
import { AppRoutes } from "./routes/";
import useContinerSourceEvent from "./lib/hooks/useContinerSourceEvent";
import { useAppDispatch } from "./hooks"
import { setCurrent } from "./stores/slices/hostSlice";
import config from "./stores/config";
import './App.css'

function App() {
  const dispatch = useAppDispatch();
  useContinerSourceEvent();
  useEffect(() => {
    dispatch(setCurrent(config.hosts[0]))
  }, [])

  return (
    <BrowserRouter children={AppRoutes} basename={"/"} />
  )
}

export default App