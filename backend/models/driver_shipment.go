package models

import (
	"database/sql"
	"fmt"
)

type DriverShipment struct {
	ShipmentID        int     `json:"shipment_id" form:"shipment_id"`
	DriverID          int     `json:"driver_id" form:"driver_id"`
	CalculatedPayment float64 `json:"calculated_payment" form:"calculated_payment"`
}

func (ds *DriverShipment) Create(db *sql.DB) error {
	query := `
		INSERT INTO driver_shipments (shipment_id, driver_id, calculated_payment)
		VALUES (?, ?, ?)
	`
	_, err := db.Exec(query, ds.ShipmentID, ds.DriverID, ds.CalculatedPayment)
	if err != nil {
		return fmt.Errorf("ошибка при привязке водителя к перевозке: %v", err)
	}
	return nil
}
