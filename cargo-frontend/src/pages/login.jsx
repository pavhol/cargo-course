import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import '../../styles/common.css';
import { FaUser, FaLock, FaSignInAlt } from 'react-icons/fa';

const Login = () => {
  const [formData, setFormData] = useState({ username: '', password: '' });
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const handleChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');

    try {
      const response = await fetch('http://localhost:8080/api/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(formData),
      });

      if (!response.ok) {
        throw new Error('Неверный логин или пароль');
      }

      const data = await response.json();

      // Сохраняем токен и роль
      localStorage.setItem('token', data.token);
      localStorage.setItem('role_id', data.role_id); // <--- ключевая строка
      localStorage.setItem('username', data.username);

      navigate('/drivers'); // или другой путь
    } catch (err) {
      setError(err.message || 'Ошибка входа');
    }
  };

  return (
    <div className="login-container fade-in">
      <form className="login-form" onSubmit={handleSubmit}>
        <h2><FaSignInAlt /> Вход в систему</h2>
        {error && <div className="error">{error}</div>}

        <div className="form-group">
          <label><FaUser /> Логин:</label>
          <input
            type="text"
            name="username"
            value={formData.username}
            onChange={handleChange}
            required
          />
        </div>

        <div className="form-group">
          <label><FaLock /> Пароль:</label>
          <input
            type="password"
            name="password"
            value={formData.password}
            onChange={handleChange}
            required
          />
        </div>

        <div className="form-actions">
          <button type="submit">
            <FaSignInAlt /> Войти
          </button>
        </div>
      </form>
    </div>
  );
};

export default Login;
