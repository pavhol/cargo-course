package main

import (
	"cargo-back/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Создаем новый роутер Gin
	r := gin.Default()

	// Настройка CORS, чтобы разрешить запросы с разных источников
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Разрешаем запросы с фронтенда
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	// Группируем API маршруты под /api
	api := r.Group("/api")

	// Водители
	api.GET("/drivers", handlers.GetDrivers)
	api.GET("/drivers/:id", handlers.GetDriverByID)
	api.POST("/drivers", handlers.CreateDriver)
	api.PUT("/drivers/:id", handlers.UpdateDriver)
	api.DELETE("/drivers/:id", handlers.DeleteDriver)

	// Города
	api.GET("/cities", handlers.GetCities)
	api.GET("/cities/:id", handlers.GetCityByID)
	api.POST("/cities", handlers.CreateCity)
	api.PUT("/cities/:id", handlers.UpdateCity)
	api.DELETE("/cities/:id", handlers.DeleteCity)

	// Маршруты
	api.GET("/routes", handlers.GetRoutes)
	api.GET("/routes/:id", handlers.GetRouteByID)
	api.POST("/routes", handlers.CreateRoute)
	api.PUT("/routes/:id", handlers.UpdateRoute)
	api.DELETE("/routes/:id", handlers.DeleteRoute)

	// Перевозки (Shipments)
	api.GET("/shipments/:id", handlers.GetShipmentByID)
	api.POST("/shipments", handlers.CreateShipment)
	api.POST("/shipments/detailed", handlers.CreateDetailedShipment)
	api.GET("/shipments/detailed", handlers.GetDetailedShipments)
	api.PUT("/shipments/:id", handlers.UpdateDetailedShipment)

	// Водитель и перевозка (DriverShipments)
	api.POST("/driver-shipments", handlers.AssignDriverToShipment)

	// Пользователи
	api.GET("/users", handlers.GetUsers)
	api.GET("/users/:id", handlers.GetUserByID)
	api.POST("/users", handlers.CreateUser)
	api.PUT("/users/:id", handlers.UpdateUser)
	api.DELETE("/users/:id", handlers.DeleteUser)

	// Авторизация
	api.POST("/login", handlers.Login)

	api.GET("/roles", handlers.GetRoles)

	// Запуск сервера на порту 8080
	r.Run(":8080")
}
