package models

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
	RoleID   int    `json:"role_id"`
	Role     string `json:"role"` // только для отображения роли по имени (JOIN)
}

// Хеширует пароль
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Проверяет пароль по хэшу
func CheckPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// Получить всех пользователей
func GetAllUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query(`
		SELECT u.id, u.username, u.role_id, r.name AS role
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username, &u.RoleID, &u.Role); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

// Получить пользователя по ID
func GetUserByID(db *sql.DB, id int) (*User, error) {
	row := db.QueryRow(`
		SELECT u.id, u.username, u.role_id, r.name AS role
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.id = ?
	`)
	var u User
	if err := row.Scan(&u.ID, &u.Username, &u.RoleID, &u.Role); err != nil {
		return nil, err
	}
	return &u, nil
}

// Получить пользователя по username (для логина)
func GetUserByUsername(db *sql.DB, username string) (*User, error) {
	row := db.QueryRow(`
		SELECT u.id, u.username, u.password_hash, u.role_id, r.name AS role
		FROM users u
		LEFT JOIN roles r ON u.role_id = r.id
		WHERE u.username = ?
	`)
	var u User
	if err := row.Scan(&u.ID, &u.Username, &u.Password, &u.RoleID, &u.Role); err != nil {
		return nil, err
	}
	return &u, nil
}

// Создание пользователя с хешированием пароля
func (u *User) Create(db *sql.DB) error {
	if u.Password == "" {
		return errors.New("пароль не может быть пустым")
	}
	hashedPassword, err := hashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = hashedPassword

	result, err := db.Exec("INSERT INTO users (username, password_hash, role_id) VALUES (?, ?, ?)", u.Username, u.Password, u.RoleID)
	if err != nil {
		return err
	}
	insertedID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	u.ID = int(insertedID)
	return nil
}

// Обновление пользователя, с условием на хеширование нового пароля
func (u *User) Update(db *sql.DB) error {
	if u.Password != "" {
		hashedPassword, err := hashPassword(u.Password)
		if err != nil {
			return err
		}
		u.Password = hashedPassword
		_, err = db.Exec("UPDATE users SET username = ?, password_hash = ?, role_id = ? WHERE id = ?", u.Username, u.Password, u.RoleID, u.ID)
		return err
	}
	_, err := db.Exec("UPDATE users SET username = ?, role_id = ? WHERE id = ?", u.Username, u.RoleID, u.ID)
	return err
}

// Удаление пользователя
func DeleteUser(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}
