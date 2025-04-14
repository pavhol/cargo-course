import React, { useState, useEffect } from 'react';
import '../../styles/common.css';
import { useNavigate } from 'react-router-dom';
import { FaSave, FaTimes, FaEdit } from 'react-icons/fa'; // иконки

const Shipments = () => {
  const [shipments, setShipments] = useState([]);
  const [routes, setRoutes] = useState([]);
  const [drivers, setDrivers] = useState([]);
  const [formData, setFormData] = useState({
    id: null,
    route_id: '',
    start_date: '',
    end_date: '',
    bonus: '',
    driver_id: ''
  });
  const [editing, setEditing] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const [sortField, setSortField] = useState('id');
  const [sortAsc, setSortAsc] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    fetchAllData();
  }, []);

  const fetchAllData = async () => {
    await Promise.all([fetchShipments(), fetchRoutes(), fetchDrivers()]);
  };

  const fetchShipments = async () => {
    try {
      const response = await fetch('http://localhost:8080/api/shipments/detailed');
      const data = await response.json();
      setShipments(data);
    } catch (error) {
      console.error('Ошибка при получении перевозок:', error);
    }
  };

  const fetchRoutes = async () => {
    try {
      const response = await fetch('http://localhost:8080/api/routes');
      const data = await response.json();
      setRoutes(data);
    } catch (error) {
      console.error('Ошибка при получении маршрутов:', error);
    }
  };

  const fetchDrivers = async () => {
    try {
      const response = await fetch('http://localhost:8080/api/drivers');
      const data = await response.json();
      setDrivers(data);
    } catch (error) {
      console.error('Ошибка при получении водителей:', error);
    }
  };

  const handleChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      const url = editing
        ? `http://localhost:8080/api/shipments/${formData.id}`
        : 'http://localhost:8080/api/shipments/detailed';
      const method = editing ? 'PUT' : 'POST';

      const response = await fetch(url, {
        method,
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          route_id: Number(formData.route_id),
          start_date: formData.start_date,
          end_date: formData.end_date,
          bonus: Number(formData.bonus),
          driver_id: Number(formData.driver_id),
        }),
        
      });

      if (!response.ok) throw new Error('Ошибка при сохранении данных');

      await fetchShipments();
      resetForm();
    } catch (error) {
      console.error('Ошибка при сохранении:', error);
    }
  };

  const handleEdit = (shipment) => {
    setEditing(true);
    setFormData({
      id: shipment.id,
      route_id: routes.find(r => r.name === shipment.route_name)?.id || '',
      start_date: shipment.start_date.slice(0, 10),
      end_date: shipment.end_date.slice(0, 10),
      bonus: shipment.bonus,
      driver_id: shipment.driver_id || ''
    });
  };

  const resetForm = () => {
    setEditing(false);
    setFormData({
      id: null,
      route_id: '',
      start_date: '',
      end_date: '',
      bonus: '',
      driver_id: ''
    });
  };

  const formatDate = (dateStr) => {
    const d = new Date(dateStr);
    return d.toLocaleDateString('ru-RU');
  };

  return (
    <div className="container">
      <h1>Перевозки</h1>
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
          <label>Маршрут:</label>
          <select name="route_id" value={formData.route_id} onChange={handleChange} required>
            <option value="">Выберите маршрут</option>
            {routes.map((route) => (
              <option key={route.id} value={route.id}>{route.name}</option>
            ))}
          </select>
        </div>
        <div className="form-group">
          <label>Водитель:</label>
          <select name="driver_id" value={formData.driver_id} onChange={handleChange} required>
            <option value="">Выберите водителя</option>
            {drivers.map((driver) => (
              <option key={driver.id} value={driver.id}>
                {`${driver.last_name} ${driver.first_name} ${driver.middle_name}`}
              </option>
            ))}
          </select>
        </div>
        <div className="form-group">
          <label>Дата начала:</label>
          <input type="date" name="start_date" value={formData.start_date} onChange={handleChange} required />
        </div>
        <div className="form-group">
          <label>Дата окончания:</label>
          <input type="date" name="end_date" value={formData.end_date} onChange={handleChange} required />
        </div>
        <div className="form-group">
          <label>Бонус:</label>
          <input type="number" name="bonus" value={formData.bonus} onChange={handleChange} />
        </div>
        <div className="form-actions">
          <button type="submit">
            <FaSave />
            {editing ? 'Сохранить' : 'Добавить'}
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
          placeholder="Поиск по маршруту или водителю..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
        />
        <select onChange={(e) => setSortField(e.target.value)} value={sortField}>
          <option value="id">ID</option>
          <option value="bonus">Бонус</option>
          <option value="calculated_payment">Оплата</option>
        </select>
        <button onClick={() => setSortAsc(!sortAsc)}>
          {sortAsc ? 'По возрастанию' : 'По убыванию'}
        </button>
      </div>

      <table className="shipments-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Маршрут</th>
            <th>Дата начала</th>
            <th>Дата окончания</th>
            <th>Бонус</th>
            <th>Водитель</th>
            <th>Оплата</th>
            <th>Действия</th>
          </tr>
        </thead>
        <tbody>
          {shipments.length > 0 ? (
            shipments.map((shipment, index) => (
              <tr key={`${shipment.id}-${shipment.driver_id}-${index}`}>
                <td>{shipment.id}</td>
                <td>{shipment.route_name}</td>
                <td>{formatDate(shipment.start_date)}</td>
                <td>{formatDate(shipment.end_date)}</td>
                <td>{shipment.bonus}</td>
                <td>{`${shipment.last_name} ${shipment.first_name} ${shipment.middle_name}`}</td>
                <td>{shipment.calculated_payment}</td>
                <td>
                  <button onClick={() => handleEdit(shipment)}>Редактировать</button>
                </td>
              </tr>
            ))
          ) : (
            <tr>
              <td colSpan="8" style={{ textAlign: 'center' }}>Нет данных</td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
};

export default Shipments;
