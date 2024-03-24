import React from "react";
import { Route, Routes } from "react-router-dom";

const DefaultLayout = React.lazy(() => import("../containers/DefaultLayout"));
const Test = React.lazy(() => import("../pages/test"));

export const AppRoutes = (
  <React.Suspense fallback={null}>
    <Routes>
    <Route path="/" element={<DefaultLayout noFooter noHeader children={<Test />}/>} />
    </Routes>
  </React.Suspense>
);