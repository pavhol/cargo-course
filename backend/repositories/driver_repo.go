package repositories

import (
	"cargo-back/database"
	"cargo-back/models"
	"database/sql"
)

type DriverRepository struct {
	db *sql.DB
}

func NewDriverRepository() (*DriverRepository, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}
	return &DriverRepository{db: db}, nil
}

// GetAll возвращает всех водителей
func (r *DriverRepository) GetAll() ([]models.Driver, error) {
	rows, err := r.db.Query("SELECT * FROM drivers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var drivers []models.Driver
	for rows.Next() {
		var d models.Driver
		err := rows.Scan(
			&d.ID,
			&d.LastName,
			&d.FirstName,
			&d.MiddleName,
			&d.Experience,
		)
		if err != nil {
			return nil, err
		}
		drivers = append(drivers, d)
	}
	return drivers, nil
}

// GetByID возвращает водителя по ID
func (r *DriverRepository) GetByID(id int) (*models.Driver, error) {
	row := r.db.QueryRow("SELECT * FROM drivers WHERE id = ?", id)
	var d models.Driver
	err := d.ScanRow(row)
	return &d, err
}
