package handlers

import (
	"cargo-back/database"
	"cargo-back/models"
	"cargo-back/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetDrivers(c *gin.Context) {
	repo, err := repositories.NewDriverRepository()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	drivers, err := repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, drivers)
}
func GetDriverByID(c *gin.Context) {
	repo, _ := repositories.NewDriverRepository()

	id, _ := strconv.Atoi(c.Param("id"))
	driver, err := repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "водитель не найден"})
		return
	}

	c.JSON(http.StatusOK, driver)
}

func CreateDriver(c *gin.Context) {
	db, _ := database.Connect()
	defer db.Close()

	var driver models.Driver
	if err := c.ShouldBindJSON(&driver); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат данных"})
		return
	}

	if err := driver.Create(db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, driver)
}

func UpdateDriver(c *gin.Context) {
	db, _ := database.Connect()
	defer db.Close()

	id, _ := strconv.Atoi(c.Param("id"))

	var driver models.Driver
	if err := c.ShouldBindJSON(&driver); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат данных"})
		return
	}
	driver.ID = id

	if err := driver.Update(db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, driver)
}

func DeleteDriver(c *gin.Context) {
	db, _ := database.Connect()
	defer db.Close()

	id, _ := strconv.Atoi(c.Param("id"))

	driver := models.Driver{ID: id}
	if err := driver.Delete(db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
