import { Link, useNavigate } from "react-router-dom";
import { useEffect, useState } from "react";

export default function Dashboard() {
  const navigate = useNavigate();
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  // Проверяем, есть ли токен в localStorage, чтобы убедиться, что пользователь авторизован
  useEffect(() => {
    const token = localStorage.getItem("token");
    if (token) {
      setIsAuthenticated(true);
    } else {
      navigate("/login"); // Если токена нет, перенаправляем на страницу логина
    }
  }, [navigate]);

  if (!isAuthenticated) {
    return <div>Загрузка...</div>; // Показать загрузку, пока не проверяется токен
  }

  return (
    <div className="flex">
      {/* Боковое меню */}
      <div className="w-64 bg-gray-800 text-white p-6">
        <h2 className="text-2xl font-bold mb-6">Панель управления</h2>
        <ul className="space-y-4">
          <li>
            <Link to="/drivers" className="text-xl hover:text-blue-500">
              Водители
            </Link>
          </li>
          <li>
            <Link to="/routes" className="text-xl hover:text-blue-500">
              Маршруты
            </Link>
          </li>
          <li>
            <Link to="/shipments" className="text-xl hover:text-blue-500">
              Отправления
            </Link>
          </li>
        </ul>
      </div>

      {/* Контент панели */}
      <div className="flex-1 p-6">
        <h1 className="text-3xl font-semibold mb-4">Добро пожаловать в панель управления!</h1>
        <p className="text-lg">
          Выберите одну из вкладок в боковом меню, чтобы начать работать с системой.
        </p>
      </div>
    </div>
  );
}
