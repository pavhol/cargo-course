import React, { useState, useEffect } from 'react';
import '../../styles/common.css';
import { useNavigate } from 'react-router-dom';
import { FaSave, FaTimes, FaEdit, FaTrash } from 'react-icons/fa';

const Routes = () => {
  const [routes, setRoutes] = useState([]);
  const [formData, setFormData] = useState({ id: null, name: '', base_payment: '' });
  const [editing, setEditing] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const [sortField, setSortField] = useState('id');
  const [sortAsc, setSortAsc] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    fetchRoutes();
  }, []);

  const fetchRoutes = async () => {
    try {
      const response = await fetch('http://localhost:8080/api/routes');
      const data = await response.json();
      setRoutes(data);
    } catch (error) {
      console.error('Ошибка при получении маршрутов:', error);
    }
  };

  const handleChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const url = editing
        ? `http://localhost:8080/api/routes/${formData.id}`
        : 'http://localhost:8080/api/routes';
      const method = editing ? 'PUT' : 'POST';

      const response = await fetch(url, {
        method,
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name: formData.name, base_payment: Number(formData.base_payment) }),
      });

      if (!response.ok) throw new Error('Ошибка при сохранении данных');

      await fetchRoutes();
      resetForm();
    } catch (error) {
      console.error('Ошибка при сохранении маршрута:', error);
    }
  };

  const handleEdit = (route) => {
    setEditing(true);
    setFormData({ id: route.id, name: route.name, base_payment: route.base_payment });
  };

  const handleDelete = async (id) => {
    if (window.confirm('Удалить этот маршрут?')) {
      try {
        const response = await fetch(`http://localhost:8080/api/routes/${id}`, { method: 'DELETE' });
        if (!response.ok) throw new Error('Ошибка при удалении маршрута');
        await fetchRoutes();
      } catch (error) {
        console.error('Ошибка при удалении маршрута:', error);
      }
    }
  };

  const resetForm = () => {
    setEditing(false);
    setFormData({ id: null, name: '', base_payment: '' });
  };

  const filteredRoutes = routes
    .filter((r) => r.name.toLowerCase().includes(searchTerm.toLowerCase()))
    .sort((a, b) => {
      const valueA = a[sortField];
      const valueB = b[sortField];
      return sortAsc ? valueA - valueB : valueB - valueA;
    });

  return (
    <div className="container">
      <h1>Маршруты</h1>
      <button className="back-button" onClick={() => navigate('/dashboard')}>
        ⬅ На главную
      </button>
      <form onSubmit={handleSubmit} className="shipment-form">
        {editing && (
          <div className="form-group">
            <label>ID:</label>
            <input type="text" name="id" value={formData.id} disabled />
          </div>
        )}
        <div className="form-group">
          <label>Название:</label>
          <input type="text" name="name" value={formData.name} onChange={handleChange} required />
        </div>
        <div className="form-group">
          <label>Базовая оплата:</label>
          <input type="number" name="base_payment" value={formData.base_payment} onChange={handleChange} required />
        </div>
        <div className="form-actions">
          <button type="submit">
            <FaSave /> {editing ? 'Сохранить' : 'Добавить'}
          </button>
          {editing && (
            <button type="button" onClick={resetForm}>
              <FaTimes /> Отмена
            </button>
          )}
        </div>
      </form>

      <div className="table-controls">
        <input
          type="text"
          placeholder="Поиск по названию маршрута..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
        />
        <select onChange={(e) => setSortField(e.target.value)} value={sortField}>
          <option value="id">ID</option>
          <option value="base_payment">Базовая оплата</option>
        </select>
        <button onClick={() => setSortAsc(!sortAsc)}>
          {sortAsc ? 'По возрастанию' : 'По убыванию'}
        </button>
      </div>

      <table className="shipments-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Название</th>
            <th>Базовая оплата</th>
            <th>Действия</th>
          </tr>
        </thead>
        <tbody>
          {filteredRoutes.length > 0 ? (
            filteredRoutes.map((route) => (
              <tr key={route.id}>
                <td>{route.id}</td>
                <td>{route.name}</td>
                <td>{route.base_payment}</td>
                <td>
                  <button onClick={() => handleEdit(route)}><FaEdit /></button>
                  <button onClick={() => handleDelete(route.id)}><FaTrash /></button>
                </td>
              </tr>
            ))
          ) : (
            <tr>
              <td colSpan="4" style={{ textAlign: 'center' }}>Нет данных</td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
};

export default Routes;
