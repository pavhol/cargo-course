package models

import (
	"database/sql"
	"errors"
	"fmt"
)

type City struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Region     string `json:"region,omitempty"`
	Country    string `json:"country,omitempty"`
	IsBaseCity bool   `json:"is_base_city"`
}

func (c *City) Create(db *sql.DB) error {
	if err := c.Validate(); err != nil {
		return err
	}

	query := `
		INSERT INTO cities (name, region, country, is_base_city)
		VALUES (?, ?, ?, ?)
	`
	result, err := db.Exec(query, c.Name, c.Region, c.Country, c.IsBaseCity)
	if err != nil {
		return fmt.Errorf("ошибка создания города: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	c.ID = int(id)
	return nil
}

func (c *City) GetByID(db *sql.DB, id int) error {
	query := `
		SELECT id, name, region, country, is_base_city
		FROM cities WHERE id = ?
	`
	row := db.QueryRow(query, id)
	return c.ScanRow(row)
}

func (c *City) Update(db *sql.DB) error {
	if c.ID == 0 {
		return errors.New("ID города не указан")
	}
	if err := c.Validate(); err != nil {
		return err
	}

	query := `
		UPDATE cities 
		SET name = ?, region = ?, country = ?, is_base_city = ?
		WHERE id = ?
	`
	_, err := db.Exec(query, c.Name, c.Region, c.Country, c.IsBaseCity, c.ID)
	return err
}

func (c *City) Delete(db *sql.DB) error {
	if c.ID == 0 {
		return errors.New("ID города не указан")
	}
	_, err := db.Exec("DELETE FROM cities WHERE id = ?", c.ID)
	return err
}

func (c *City) Validate() error {
	if c.Name == "" {
		return errors.New("название города обязательно")
	}
	return nil
}

func (c *City) ScanRow(row *sql.Row) error {
	return row.Scan(&c.ID, &c.Name, &c.Region, &c.Country, &c.IsBaseCity)
}

func ScanCityRows(rows *sql.Rows) ([]City, error) {
	var cities []City
	for rows.Next() {
		var c City
		err := rows.Scan(&c.ID, &c.Name, &c.Region, &c.Country, &c.IsBaseCity)
		if err != nil {
			return nil, err
		}
		cities = append(cities, c)
	}
	return cities, nil
}
