// src/pages/Drivers.jsx
import React, { useEffect, useState } from 'react';
import '../../styles/common.css';
import { useNavigate } from 'react-router-dom';
import { FaTrash, FaEdit, FaSave, FaTimes } from 'react-icons/fa';

const Drivers = () => {
  const [drivers, setDrivers] = useState([]);
  const [formData, setFormData] = useState({
    id: null,
    first_name: '',
    last_name: '',
    middle_name: '',
    experience: ''
  });
  const [editing, setEditing] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const [sortField, setSortField] = useState('id');
  const [sortAsc, setSortAsc] = useState(true);
  const navigate = useNavigate();

  const role = localStorage.getItem('role_id');
  const isDriver = role === '3';

  useEffect(() => {
    fetchDrivers();
  }, []);

  const fetchDrivers = async () => {
    try {
      const res = await fetch('http://localhost:8080/api/drivers');
      const data = await res.json();
      setDrivers(data);
    } catch (err) {
      console.error('Ошибка при загрузке водителей:', err);
    }
  };

  const handleChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    if(isDriver) return;
    e.preventDefault();
    try {
      const url = editing ? `http://localhost:8080/api/drivers/${formData.id}` : 'http://localhost:8080/api/drivers';
      const method = editing ? 'PUT' : 'POST';

      const res = await fetch(url, {
        method,
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          first_name: formData.first_name,
          last_name: formData.last_name,
          middle_name: formData.middle_name,
          experience: Number(formData.experience),
        })
      });

      if (!res.ok) throw new Error('Ошибка при сохранении');

      fetchDrivers();
      resetForm();
    } catch (err) {
      console.error('Ошибка при сохранении:', err);
    }
  };

  const handleEdit = (driver) => {
    if (isDriver) return;
    setEditing(true);
    setFormData({ ...driver });
  };

  const handleDelete = async (id) => {
    if (isDriver||!window.confirm('Удалить водителя?')) return;
    try {
      const res = await fetch(`http://localhost:8080/api/drivers/${id}`, { method: 'DELETE' });
      if (!res.ok) throw new Error('Ошибка при удалении');
      fetchDrivers();
    } catch (err) {
      console.error('Ошибка при удалении:', err);
    }
  };

  const resetForm = () => {
    setEditing(false);
    setFormData({ id: null, first_name: '', last_name: '', middle_name: '', experience: '' });
  };

  const filtered = drivers.filter(d =>
    `${d.last_name} ${d.first_name} ${d.middle_name}`.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const sorted = [...filtered].sort((a, b) => {
    const aVal = a[sortField];
    const bVal = b[sortField];
    return (aVal > bVal ? 1 : -1) * (sortAsc ? 1 : -1);
  });

  return (
    <div className="container">
      <h1>Водители</h1>
      <button className="back-button" onClick={() => navigate('/dashboard')}>⬅ На главную</button>
      {!isDriver && (
      <form onSubmit={handleSubmit} className="shipment-form">
        {editing && (
          <div className="form-group">
            <label>ID:</label>
            <input type="text" name="id" value={formData.id} disabled />
          </div>
        )}
        <div className="form-group">
          <label>Фамилия:</label>
          <input type="text" name="last_name" value={formData.last_name} onChange={handleChange} required />
        </div>
        <div className="form-group">
          <label>Имя:</label>
          <input type="text" name="first_name" value={formData.first_name} onChange={handleChange} required />
        </div>
        <div className="form-group">
          <label>Отчество:</label>
          <input type="text" name="middle_name" value={formData.middle_name} onChange={handleChange} />
        </div>
        <div className="form-group">
          <label>Стаж (лет):</label>
          <input type="number" name="experience" value={formData.experience} onChange={handleChange} required />
        </div>
        <div className="form-actions">
          <button type="submit"><FaSave /> {editing ? 'Сохранить' : 'Добавить'}</button>
          {editing && <button type="button" onClick={resetForm}><FaTimes /> Отмена</button>}
        </div>
      </form>
      )}
      <div className="table-controls">
        <input
          type="text"
          placeholder="Поиск по имени..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
        />
        <select onChange={(e) => setSortField(e.target.value)} value={sortField}>
          <option value="id">ID</option>
          <option value="experience">Стаж</option>
          <option value="last_name">Фамилия</option>
        </select>
        <button onClick={() => setSortAsc(!sortAsc)}>{sortAsc ? '↑' : '↓'}</button>
      </div>

      <table className="shipments-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Фамилия</th>
            <th>Имя</th>
            <th>Отчество</th>
            <th>Стаж</th>
            <th>Действия</th>
          </tr>
        </thead>
        <tbody>
          {sorted.length > 0 ? (
            sorted.map(driver => (
              <tr key={driver.id}>
                <td>{driver.id}</td>
                <td>{driver.last_name}</td>
                <td>{driver.first_name}</td>
                <td>{driver.middle_name}</td>
                <td>{driver.experience}</td>
                <td>{!isDriver && (
                  <>
                  <button onClick={() => handleEdit(driver)}><FaEdit /></button>
                  <button onClick={() => handleDelete(driver.id)}><FaTrash /></button>
                  </>
                )}
                  </td>
              </tr>
            ))
          ) : (
            <tr>
              <td colSpan="6" style={{ textAlign: 'center' }}>Нет данных</td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
};

export default Drivers;
