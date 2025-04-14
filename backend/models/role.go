package models

import "database/sql"

type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func GetAllRoles(db *sql.DB) ([]Role, error) {
	rows, err := db.Query("SELECT id, name FROM roles")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []Role
	for rows.Next() {
		var r Role
		if err := rows.Scan(&r.ID, &r.Name); err != nil {
			return nil, err
		}
		roles = append(roles, r)
	}
	return roles, nil
}
