package models

import (
	"database/sql"
	"errors"
	"fmt"
)

type Driver struct {
	ID         int    `json:"id"`
	LastName   string `json:"last_name"`
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name,omitempty"`
	Experience int    `json:"experience"`
}

// Create - Создает нового водителя в БД
func (d *Driver) Create(db *sql.DB) error {
	if err := d.Validate(); err != nil {
		return err
	}

	query := `
		INSERT INTO drivers 
		(last_name, first_name, middle_name, experience) 
		VALUES (?, ?, ?, ?)
	`
	result, err := db.Exec(
		query,
		d.LastName,
		d.FirstName,
		d.MiddleName,
		d.Experience,
	)
	if err != nil {
		return fmt.Errorf("ошибка создания водителя: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("ошибка получения ID: %v", err)
	}

	d.ID = int(id)
	return nil
}

// GetByID - Получает водителя по ID
func (d *Driver) GetByID(db *sql.DB, id int) error {
	query := `
		SELECT id, last_name, first_name, middle_name, experience 
		FROM drivers 
		WHERE id = ?
	`
	row := db.QueryRow(query, id)
	return d.ScanRow(row)
}

// Update - Обновляет данные водителя
func (d *Driver) Update(db *sql.DB) error {
	if d.ID == 0 {
		return errors.New("ID водителя не указан")
	}

	if err := d.Validate(); err != nil {
		return err
	}

	query := `
		UPDATE drivers 
		SET last_name = ?, first_name = ?, middle_name = ?, experience = ? 
		WHERE id = ?
	`
	_, err := db.Exec(
		query,
		d.LastName,
		d.FirstName,
		d.MiddleName,
		d.Experience,
		d.ID,
	)
	return err
}

// Delete - Удаляет водителя
func (d *Driver) Delete(db *sql.DB) error {
	if d.ID == 0 {
		return errors.New("ID водителя не указан")
	}

	query := `DELETE FROM drivers WHERE id = ?`
	_, err := db.Exec(query, d.ID)
	return err
}

// Validate - Проверяет корректность данных
func (d *Driver) Validate() error {
	if d.LastName == "" {
		return errors.New("фамилия обязательна для заполнения")
	}
	if d.FirstName == "" {
		return errors.New("имя обязательно для заполнения")
	}
	if d.Experience < 0 {
		return errors.New("опыт не может быть отрицательным")
	}
	return nil
}

// GetAllDrivers - Получает всех водителей
func GetAllDrivers(db *sql.DB) ([]Driver, error) {
	query := `
		SELECT id, last_name, first_name, middle_name, experience 
		FROM drivers
	`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var drivers []Driver
	for rows.Next() {
		var d Driver
		if err := d.ScanRows(rows); err != nil {
			return nil, err
		}
		drivers = append(drivers, d)
	}
	return drivers, nil
}

// scanRow - Сканирует одну строку результата
func (d *Driver) ScanRow(row *sql.Row) error {
	return row.Scan(
		&d.ID,
		&d.LastName,
		&d.FirstName,
		&d.MiddleName,
		&d.Experience,
	)
}

// scanRows - Сканирует несколько строк результатов
func (d *Driver) ScanRows(rows *sql.Rows) error {
	return rows.Scan(
		&d.ID,
		&d.LastName,
		&d.FirstName,
		&d.MiddleName,
		&d.Experience,
	)
}

// ToMap - Преобразует структуру в map (для обновления)
func (d *Driver) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"last_name":   d.LastName,
		"first_name":  d.FirstName,
		"middle_name": d.MiddleName,
		"experience":  d.Experience,
	}
}
