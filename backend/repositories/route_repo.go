package repositories

import (
	"cargo-back/database"
	"cargo-back/models"
	"database/sql"
)

type RouteRepository struct {
	db *sql.DB
}

func NewRouteRepository() (*RouteRepository, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}
	return &RouteRepository{db: db}, nil
}

func (r *RouteRepository) GetAll() ([]models.Route, error) {
	rows, err := r.db.Query("SELECT * FROM routes")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var routes []models.Route
	for rows.Next() {
		var route models.Route
		if err := route.ScanRows(rows); err != nil {
			return nil, err
		}
		routes = append(routes, route)
	}
	return routes, nil
}

func (r *RouteRepository) GetByID(id int) (*models.Route, error) {
	row := r.db.QueryRow("SELECT * FROM routes WHERE id = ?", id)
	var route models.Route
	if err := route.ScanRow(row); err != nil {
		return nil, err
	}
	return &route, nil
}
