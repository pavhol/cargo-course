package repositories

import (
	"cargo-back/database"
	"cargo-back/models"
)

type DriverShipmentRepository struct{}

func NewDriverShipmentRepository() (*DriverShipmentRepository, error) {
	return &DriverShipmentRepository{}, nil
}

func (r *DriverShipmentRepository) AssignDriver(ds *models.DriverShipment) error {
	db, err := database.Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	return ds.Create(db)
}
