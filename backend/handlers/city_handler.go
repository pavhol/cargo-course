package handlers

import (
	"cargo-back/database"
	"cargo-back/models"
	"cargo-back/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetCities(c *gin.Context) {
	repo, err := repositories.NewCityRepository()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	cities, err := repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cities)
}

func GetCityByID(c *gin.Context) {
	repo, err := repositories.NewCityRepository()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	city, err := repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "город не найден"})
		return
	}
	c.JSON(http.StatusOK, city)
}

func CreateCity(c *gin.Context) {
	db, _ := database.Connect()
	defer db.Close()

	var city models.City
	if err := c.ShouldBindJSON(&city); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат данных"})
		return
	}

	if err := city.Create(db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, city)
}

func UpdateCity(c *gin.Context) {
	db, _ := database.Connect()
	defer db.Close()

	id, _ := strconv.Atoi(c.Param("id"))

	var city models.City
	if err := c.ShouldBindJSON(&city); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат данных"})
		return
	}
	city.ID = id

	if err := city.Update(db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, city)
}

func DeleteCity(c *gin.Context) {
	db, _ := database.Connect()
	defer db.Close()

	id, _ := strconv.Atoi(c.Param("id"))

	city := models.City{ID: id}
	if err := city.Delete(db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
