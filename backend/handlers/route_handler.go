package handlers

import (
	"cargo-back/database"
	"cargo-back/models"
	"cargo-back/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetRoutes(c *gin.Context) {
	repo, err := repositories.NewRouteRepository()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	routes, err := repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, routes)
}

func GetRouteByID(c *gin.Context) {
	repo, _ := repositories.NewRouteRepository()
	id, _ := strconv.Atoi(c.Param("id"))

	route, err := repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "маршрут не найден"})
		return
	}
	c.JSON(http.StatusOK, route)
}

func CreateRoute(c *gin.Context) {
	db, _ := database.Connect()
	defer db.Close()

	var route models.Route
	// if err := c.ShouldBindJSON(&route); err != nil {
	if err := c.ShouldBind(&route); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат"})
		return
	}

	if err := route.Create(db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, route)
}

func UpdateRoute(c *gin.Context) {
	db, _ := database.Connect()
	defer db.Close()

	id, _ := strconv.Atoi(c.Param("id"))

	var route models.Route
	if err := c.ShouldBindJSON(&route); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат"})
		return
	}
	route.ID = id

	if err := route.Update(db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, route)
}

func DeleteRoute(c *gin.Context) {
	db, _ := database.Connect()
	defer db.Close()

	id, _ := strconv.Atoi(c.Param("id"))
	route := models.Route{ID: id}
	if err := route.Delete(db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}
