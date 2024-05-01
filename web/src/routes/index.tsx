import React from "react";
import { Route, Routes } from "react-router-dom";

const Layout = React.lazy(() => import("../layout"));
const LogViewer = React.lazy(() => import("../pages/LogViewer"));

export const AppRoutes = (
  <React.Suspense>
    <Routes>
      <Route path="/" element={<Layout noFooter children={null}/>} />
      <Route path="/container/:id" element={<Layout noFooter children={<LogViewer />}/>} />
    </Routes>
  </React.Suspense>
);