package handlers

import (
	"cargo-back/models"
	"cargo-back/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AssignDriverToShipment(c *gin.Context) {
	var ds models.DriverShipment
	if err := c.ShouldBind(&ds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат"})
		return
	}

	repo, _ := repositories.NewDriverShipmentRepository()
	if err := repo.AssignDriver(&ds); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ds)
}
