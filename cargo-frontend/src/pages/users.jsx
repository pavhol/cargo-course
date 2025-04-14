import React, { useState, useEffect } from 'react';
import '../../styles/common.css';
import { useNavigate } from 'react-router-dom';
import { FaSave, FaTimes, FaEdit, FaTrash } from 'react-icons/fa';

const Users = () => {
  const [users, setUsers] = useState([]);
  const [roles, setRoles] = useState([]);
  const [formData, setFormData] = useState({
    id: null,
    username: '',
    password: '',
    role_id: ''
  });
  const [editing, setEditing] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const [sortField, setSortField] = useState('id');
  const [sortAsc, setSortAsc] = useState(true);
  const navigate = useNavigate();

  const role = localStorage.getItem('role_id');
  const isAdmin = role === '1';

  useEffect(() => {
    fetchUsers();
    fetchRoles();
  }, []);

  const fetchUsers = async () => {
    try {
      const res = await fetch('http://localhost:8080/api/users');
      const data = await res.json();
      setUsers(data);
    } catch (error) {
      console.error('Ошибка при получении пользователей:', error);
    }
  };

  const fetchRoles = async () => {
    try {
      const res = await fetch('http://localhost:8080/api/roles');
      const data = await res.json();
      setRoles(data);
    } catch (error) {
      console.error('Ошибка при получении ролей:', error);
    }
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: name === 'role_id' ? parseInt(value) : value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (!isAdmin) return;

    try {
      if (editing) {
        const url = `http://localhost:8080/api/users/${formData.id}`;
        const method = 'PUT';
        const payload = {
          username: formData.username,
          role_id: formData.role_id
        };
        if (formData.password.trim() !== '') {
          payload.password = formData.password;
        }
        const res = await fetch(url, {
          method,
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(payload),
        });
        if (!res.ok) throw new Error('Ошибка при сохранении данных');
      } else {
        const url = 'http://localhost:8080/api/users';
        const res = await fetch(url, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            username: formData.username,
            password: formData.password,
            role_id: formData.role_id,
          }),
        });
        if (!res.ok) throw new Error('Ошибка при сохранении данных');
      }
      await fetchUsers();
      resetForm();
    } catch (error) {
      console.error('Ошибка при сохранении:', error);
    }
  };

  const handleEdit = (user) => {
    if (!isAdmin) return;
    const role = roles.find(role => role.name === user.role);
    setEditing(true);
    setFormData({
      id: user.id,
      username: user.username,
      password: '',
      role_id: role ? role.id : ''
    });
  };

  const handleDelete = async (id) => {
    if (!isAdmin || !window.confirm('Удалить пользователя?')) return;

    try {
      const res = await fetch(`http://localhost:8080/api/users/${id}`, {
        method: 'DELETE',
      });
      if (!res.ok) throw new Error('Ошибка при удалении пользователя');
      await fetchUsers();
    } catch (error) {
      console.error('Ошибка при удалении пользователя:', error);
    }
  };

  const resetForm = () => {
    setEditing(false);
    setFormData({ id: null, username: '', password: '', role_id: '' });
  };

  const filteredUsers = users.filter(user =>
    user.username.toLowerCase().includes(searchTerm.toLowerCase())
  );
  const sortedUsers = [...filteredUsers].sort((a, b) => {
    const valA = a[sortField];
    const valB = b[sortField];
    if (typeof valA === 'number' && typeof valB === 'number') {
      return sortAsc ? valA - valB : valB - valA;
    } else {
      return sortAsc
        ? ('' + valA).localeCompare('' + valB)
        : ('' + valB).localeCompare('' + valA);
    }
  });

  return (
    <div className="container">
      <h1>Пользователи</h1>
      <button className="back-button" onClick={() => navigate('/dashboard')}>
        ⬅ На главную
      </button>

      {isAdmin && (
        <form onSubmit={handleSubmit} className="shipment-form">
          {editing && (
            <div className="form-group">
              <label>ID:</label>
              <input type="text" name="id" value={formData.id} disabled />
            </div>
          )}
          <div className="form-group">
            <label>Логин:</label>
            <input type="text" name="username" value={formData.username} onChange={handleChange} required />
          </div>
          <div className="form-group">
            <label>{editing ? 'Новый пароль (оставьте пустым):' : 'Пароль:'}</label>
            <input
              type="password"
              name="password"
              value={formData.password}
              onChange={handleChange}
              required={!editing}
              placeholder={editing ? 'Оставьте пустым' : ''}
            />
          </div>
          <div className="form-group">
            <label>Роль:</label>
            <select name="role_id" value={formData.role_id} onChange={handleChange} required>
              <option value="">Выберите роль</option>
              {roles.map(role => (
                <option key={role.id} value={role.id}>{role.name}</option>
              ))}
            </select>
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
          placeholder="Поиск по логину..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
        />
        <select onChange={(e) => setSortField(e.target.value)} value={sortField}>
          <option value="id">ID</option>
          <option value="username">Логин</option>
          <option value="role">Роль</option>
        </select>
        <button onClick={() => setSortAsc(!sortAsc)}>{sortAsc ? '↑' : '↓'}</button>
      </div>

      <table className="shipments-table">
        <thead>
          <tr>
            <th>ID</th>
            <th>Логин</th>
            <th>Роль</th>
            <th>Действия</th>
          </tr>
        </thead>
        <tbody>
          {sortedUsers.length > 0 ? (
            sortedUsers.map((user) => (
              <tr key={user.id}>
                <td>{user.id}</td>
                <td>{user.username}</td>
                <td>{user.role}</td>
                <td>
                  {isAdmin && (
                    <>
                      <button onClick={() => handleEdit(user)} title="Редактировать"><FaEdit /></button>
                      <button onClick={() => handleDelete(user.id)} title="Удалить"><FaTrash /></button>
                    </>
                  )}
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

export default Users;
