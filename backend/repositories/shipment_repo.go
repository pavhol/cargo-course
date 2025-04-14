package repositories

import (
	"cargo-back/database"
	"cargo-back/models"
)

type ShipmentRepository struct{}

func NewShipmentRepository() (*ShipmentRepository, error) {
	return &ShipmentRepository{}, nil
}

func (r *ShipmentRepository) GetAll() ([]models.Shipment, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, route_id, start_date, end_date, bonus FROM shipments")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shipments []models.Shipment
	for rows.Next() {
		var s models.Shipment
		if err := s.ScanRows(rows); err != nil {
			return nil, err
		}
		shipments = append(shipments, s)
	}
	return shipments, nil
}

func (r *ShipmentRepository) GetByID(id int) (*models.Shipment, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	row := db.QueryRow("SELECT id, route_id, start_date, end_date, bonus FROM shipments WHERE id = ?", id)
	var s models.Shipment
	if err := s.ScanRow(row); err != nil {
		return nil, err
	}
	return &s, nil
}
