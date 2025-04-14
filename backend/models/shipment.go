package models

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type ShipmentDetail struct {
	ID                int      `json:"id"`
	StartDate         string   `json:"start_date"`
	EndDate           string   `json:"end_date"`
	Bonus             float64  `json:"bonus"`
	RouteName         string   `json:"route_name"`
	DriverID          *int     `json:"driver_id,omitempty"`
	FirstName         *string  `json:"first_name,omitempty"`
	LastName          *string  `json:"last_name,omitempty"`
	MiddleName        *string  `json:"middle_name,omitempty"`
	Experience        *int     `json:"experience,omitempty"`
	CalculatedPayment *float64 `json:"calculated_payment,omitempty"`
}

type Shipment struct {
	ID        int       `json:"id" form:"id"`
	RouteID   int       `json:"route_id" form:"route_id"`
	StartDate time.Time `json:"start_date" form:"start_date" time_format:"2006-01-02"`
	EndDate   time.Time `json:"end_date" form:"end_date" time_format:"2006-01-02"`
	Bonus     *float64  `json:"bonus,omitempty" form:"bonus"`
}

func (s *Shipment) Validate() error {
	if s.RouteID <= 0 {
		return errors.New("маршрут обязателен")
	}
	if s.EndDate.Before(s.StartDate) {
		return errors.New("дата окончания не может быть раньше начала")
	}
	return nil
}

func (s *Shipment) Create(db *sql.DB) error {
	if err := s.Validate(); err != nil {
		return err
	}

	query := `INSERT INTO shipments (route_id, start_date, end_date, bonus) VALUES (?, ?, ?, ?)`
	result, err := db.Exec(query, s.RouteID, s.StartDate, s.EndDate, s.Bonus)
	if err != nil {
		return fmt.Errorf("ошибка при создании перевозки: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	s.ID = int(id)
	return nil
}

// Update обновляет перевозку в базе данных
func (s *Shipment) Update(db *sql.DB) error {
	query := `
        UPDATE shipments
        SET route_id = ?, start_date = ?, end_date = ?, bonus = ?
        WHERE id = ?
    `
	result, err := db.Exec(query, s.RouteID, s.StartDate, s.EndDate, s.Bonus, s.ID)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("перевозка с id %d не найдена", s.ID)
	}
	return nil
}

func (s *Shipment) ScanRow(row *sql.Row) error {
	return row.Scan(&s.ID, &s.RouteID, &s.StartDate, &s.EndDate, &s.Bonus)
}

func (s *Shipment) ScanRows(rows *sql.Rows) error {
	return rows.Scan(&s.ID, &s.RouteID, &s.StartDate, &s.EndDate, &s.Bonus)
}
