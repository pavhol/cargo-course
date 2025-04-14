import React, { useEffect, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import "./navbar.css";

const Navbar = () => {
  const [logoUrl, setLogoUrl] = useState("");
  const navigate = useNavigate();

  useEffect(() => {
    fetch("http://localhost:8080/api/logo")
      .then(res => res.json())
      .then(data => setLogoUrl("http://localhost:8080" + data.logo))
      .catch(err => console.error("Ошибка загрузки лого:", err));
  }, []);

  const handleLogout = () => {
    localStorage.removeItem("token");
    navigate("/login");
  };

  return (
    <nav className="navbar">
      <div className="navbar-logo" onClick={() => navigate("/dashboard")}>
        {logoUrl && <img src={logoUrl} alt="logo" style={{ height: "40px" }} />}
      </div>
      <ul className="navbar-links">
        <li><Link to="/drivers">Водители</Link></li>
        <li><Link to="/routes">Маршруты</Link></li>
        <li><Link to="/shipments">Отправления</Link></li>
        <li><Link to="/users">Пользователи</Link></li>
      </ul>
      <button className="logout-button" onClick={handleLogout}>Выйти</button>
    </nav>
  );
};

export default Navbar;
