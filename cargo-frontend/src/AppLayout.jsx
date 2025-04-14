// src/layouts/AppLayout.jsx
import Navbar from "../src/pages/navbar";
import { Outlet } from "react-router-dom";

const AppLayout = () => {
  return (
    <>
      <Navbar />
      <main className="main-content">
        <Outlet />
      </main>
    </>
  );
};

export default AppLayout;
