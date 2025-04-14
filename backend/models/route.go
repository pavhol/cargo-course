package models

import (
	"database/sql"
	"errors"
	"fmt"
)

type Route struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" form:"name"`
	Distance    float64 `json:"distance" form:"distance"`
	Days        int     `json:"days" form:"days"`
	BasePayment float64 `json:"base_payment" form:"base_payment"`
	CityFromID  *int    `json:"city_from_id,omitempty" form:"city_from_id"`
	CityToID    *int    `json:"city_to_id,omitempty" form:"city_to_id"`
}

// Validate проверяет корректность данных
func (r *Route) Validate() error {
	if r.Name == "" {
		return errors.New("название маршрута обязательно")
	}
	if r.Distance <= 0 {
		return errors.New("дистанция должна быть положительной")
	}
	if r.Days <= 0 {
		return errors.New("длительность маршрута должна быть положительной")
	}
	if r.BasePayment <= 0 {
		return errors.New("базовая оплата должна быть положительной")
	}
	return nil
}

func (r *Route) Create(db *sql.DB) error {
	if err := r.Validate(); err != nil {
		return err
	}
	query := `
		INSERT INTO routes (name, distance, days, base_payment, city_from_id, city_to_id)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	result, err := db.Exec(query, r.Name, r.Distance, r.Days, r.BasePayment, r.CityFromID, r.CityToID)
	if err != nil {
		return fmt.Errorf("ошибка при создании маршрута: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	r.ID = int(id)
	return nil
}

func (r *Route) Update(db *sql.DB) error {
	if r.ID == 0 {
		return errors.New("ID маршрута не указан")
	}
	if err := r.Validate(); err != nil {
		return err
	}
	query := `
		UPDATE routes SET name = ?, distance = ?, days = ?, base_payment = ?, city_from_id = ?, city_to_id = ?
		WHERE id = ?
	`
	_, err := db.Exec(query, r.Name, r.Distance, r.Days, r.BasePayment, r.CityFromID, r.CityToID, r.ID)
	return err
}

func (r *Route) Delete(db *sql.DB) error {
	if r.ID == 0 {
		return errors.New("ID маршрута не указан")
	}
	_, err := db.Exec("DELETE FROM routes WHERE id = ?", r.ID)
	return err
}

func (r *Route) ScanRow(row *sql.Row) error {
	return row.Scan(&r.ID, &r.Name, &r.Distance, &r.Days, &r.BasePayment, &r.CityFromID, &r.CityToID)
}

func (r *Route) ScanRows(rows *sql.Rows) error {
	return rows.Scan(&r.ID, &r.Name, &r.Distance, &r.Days, &r.BasePayment, &r.CityFromID, &r.CityToID)
}
