package repositories

import (
	"cargo-back/database"
	"cargo-back/models"
	"database/sql"
)

type CityRepository struct {
	db *sql.DB
}

func NewCityRepository() (*CityRepository, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}
	return &CityRepository{db: db}, nil
}

func (r *CityRepository) GetAll() ([]models.City, error) {
	rows, err := r.db.Query("SELECT * FROM cities")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return models.ScanCityRows(rows)
}

func (r *CityRepository) GetByID(id int) (*models.City, error) {
	row := r.db.QueryRow("SELECT * FROM cities WHERE id = ?", id)
	var city models.City
	err := city.ScanRow(row)
	return &city, err
}
