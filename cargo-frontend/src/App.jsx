import { BrowserRouter as Router, Routes, Route, Navigate } from "react-router-dom";
import Login from "./pages/Login";
import AppLayout from "./AppLayout";
import Dashboard from "./pages/Dashboard";
import Drivers from "./pages/Drivers";
import RoutesPage from "./pages/Routes";
import Shipments from "./pages/Shipments";
import Users from "./pages/Users";

function App() {
  const token = localStorage.getItem("token");

  return (
    <Router>
      <Routes>
        <Route path="/login" element={<Login />} />

        {/* Все защищённые маршруты внутри AppLayout */}
        {token ? (
          <Route element={<AppLayout />}>
            <Route path="/dashboard" element={<Dashboard />} />
            <Route path="/drivers" element={<Drivers />} />
            <Route path="/routes" element={<RoutesPage />} />
            <Route path="/shipments" element={<Shipments />} />
            <Route path="/users" element={<Users />} />
            <Route path="*" element={<Navigate to="/dashboard" />} />
          </Route>
        ) : (
          <Route path="*" element={<Navigate to="/login" />} />
        )}
      </Routes>
    </Router>
  );
}

export default App;
