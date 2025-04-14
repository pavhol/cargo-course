package handlers

import (
	"cargo-back/database"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func GetRoles(c *gin.Context) {
	db, _ := database.Connect()
	defer db.Close()

	rows, err := db.Query("SELECT id, name FROM roles")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении ролей"})
		return
	}
	defer rows.Close()

	var roles []Role
	for rows.Next() {
		var role Role
		if err := rows.Scan(&role.ID, &role.Name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сканировании ролей"})
			return
		}
		roles = append(roles, role)
	}

	c.JSON(http.StatusOK, roles)
}
