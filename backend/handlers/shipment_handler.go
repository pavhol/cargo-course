package handlers

import (
	"cargo-back/database"
	"cargo-back/models"
	"cargo-back/repositories"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetDetailedShipments(c *gin.Context) {
	db, _ := database.Connect() // Подключение к БД (как ты его делаешь)
	defer db.Close()
	rows, err := db.Query(`
		SELECT
			s.id,
			s.start_date,
			s.end_date,
			s.bonus,
			r.name AS route_name,
			d.id,
			d.first_name,
			d.last_name,
			d.middle_name,
			d.experience,
			ds.calculated_payment
		FROM shipments s
		LEFT JOIN routes r ON s.route_id = r.id
		LEFT JOIN driver_shipments ds ON s.id = ds.shipment_id
		LEFT JOIN drivers d ON ds.driver_id = d.id
	`)
	if err != nil {
		c.JSON(500, gin.H{"error": "Ошибка при получении данных"})
		return
	}
	defer rows.Close()

	shipments := []models.ShipmentDetail{}

	for rows.Next() {
		var sh models.ShipmentDetail
		err := rows.Scan(
			&sh.ID,
			&sh.StartDate,
			&sh.EndDate,
			&sh.Bonus,
			&sh.RouteName,
			&sh.DriverID,
			&sh.FirstName,
			&sh.LastName,
			&sh.MiddleName,
			&sh.Experience,
			&sh.CalculatedPayment,
		)
		if err != nil {
			log.Println("Ошибка при сканировании:", err)
			continue
		}
		shipments = append(shipments, sh)
	}

	c.JSON(200, shipments)
}
func UpdateDetailedShipment(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		RouteID   int     `json:"route_id"`
		StartDate string  `json:"start_date"`
		EndDate   string  `json:"end_date"`
		Bonus     float64 `json:"bonus"`
		DriverID  int     `json:"driver_id"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Неверный формат запроса"})
		return
	}

	db, _ := database.Connect()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		c.JSON(500, gin.H{"error": "Не удалось начать транзакцию"})
		return
	}

	// Обновляем таблицу shipments
	_, err = tx.Exec(`
		UPDATE shipments SET 
			route_id = ?, 
			start_date = ?, 
			end_date = ?, 
			bonus = ?
		WHERE id = ?`,
		input.RouteID, input.StartDate, input.EndDate, input.Bonus, id,
	)
	if err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "Ошибка при обновлении shipments"})
		return
	}

	// Удаляем старые связи с водителями
	_, err = tx.Exec(`DELETE FROM driver_shipments WHERE shipment_id = ?`, id)
	if err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "Ошибка при удалении старых связей с водителями"})
		return
	}

	// Получаем базовую оплату из маршрута
	var basePayment float64
	err = tx.QueryRow(`SELECT base_payment FROM routes WHERE id = ?`, input.RouteID).Scan(&basePayment)
	if err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "Ошибка при получении оплаты маршрута"})
		return
	}

	// Считаем итоговую оплату
	calculatedPayment := basePayment + input.Bonus

	// Добавляем новую связь
	_, err = tx.Exec(`
		INSERT INTO driver_shipments (shipment_id, driver_id, calculated_payment)
		VALUES (?, ?, ?)`,
		id, input.DriverID, calculatedPayment,
	)
	if err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "Ошибка при создании связи shipment-driver"})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(500, gin.H{"error": "Ошибка при сохранении транзакции"})
		return
	}

	c.JSON(200, gin.H{"message": "Перевозка успешно обновлена"})
}

func CreateDetailedShipment(c *gin.Context) {
	var input struct {
		RouteID   int     `json:"route_id"`
		StartDate string  `json:"start_date"`
		EndDate   string  `json:"end_date"`
		Bonus     float64 `json:"bonus"`
		DriverID  int     `json:"driver_id"`
	}

	if err := c.BindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Неверный формат запроса"})
		return
	}

	db, err := database.Connect()
	if err != nil {
		c.JSON(500, gin.H{"error": "Ошибка подключения к базе данных"})
		return
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		c.JSON(500, gin.H{"error": "Не удалось начать транзакцию"})
		return
	}

	// Получаем base_payment маршрута
	var basePayment float64
	err = tx.QueryRow("SELECT base_payment FROM routes WHERE id = ?", input.RouteID).Scan(&basePayment)
	if err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "Ошибка при получении base_payment маршрута"})
		return
	}

	calculatedPayment := basePayment + input.Bonus

	// Вставляем запись в shipments
	res, err := tx.Exec(`
		INSERT INTO shipments (route_id, start_date, end_date, bonus)
		VALUES (?, ?, ?, ?)
	`, input.RouteID, input.StartDate, input.EndDate, input.Bonus)
	if err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "Ошибка при создании shipments"})
		return
	}

	// Получаем id созданной записи
	createdID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "Ошибка при получении ID новой перевозки"})
		return
	}

	// Вставляем запись в driver_shipments
	_, err = tx.Exec(`
		INSERT INTO driver_shipments (shipment_id, driver_id, calculated_payment)
		VALUES (?, ?, ?)
	`, createdID, input.DriverID, calculatedPayment)
	if err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "Ошибка при добавлении в driver_shipments"})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(500, gin.H{"error": "Ошибка при коммите транзакции"})
		return
	}

	c.JSON(200, gin.H{"message": "Перевозка успешно создана"})
}

func GetShipmentByID(c *gin.Context) {
	repo, _ := repositories.NewShipmentRepository()
	id, _ := strconv.Atoi(c.Param("id"))
	shipment, err := repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "перевозка не найдена"})
		return
	}
	c.JSON(http.StatusOK, shipment)
}

func CreateShipment(c *gin.Context) {
	db, _ := database.Connect()
	defer db.Close()

	var shipment models.Shipment
	if err := c.ShouldBind(&shipment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "неверный формат"})
		return
	}

	if err := shipment.Create(db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, shipment)
}

// UpdateShipment обновляет данные перевозки по id
func UpdateShipment(c *gin.Context) {
	// Подключаемся к базе данных
	db, err := database.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка подключения к базе данных"})
		return
	}
	defer db.Close()

	// Получаем id из параметров запроса
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный идентификатор"})
		return
	}

	// Читаем данные из запроса и связываем с моделью Shipment
	var shipment models.Shipment
	if err := c.ShouldBindJSON(&shipment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат данных"})
		return
	}
	// Устанавливаем id, полученный из URL, чтобы обновлять именно эту запись
	shipment.ID = id

	// Вызываем метод обновления записи.
	// Здесь предполагается, что в модели Shipment реализован метод Update для обновления записи.
	if err := shipment.Update(db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при обновлении перевозки: " + err.Error()})
		return
	}

	// Отправляем обновленную запись обратно клиенту
	c.JSON(http.StatusOK, shipment)
}
