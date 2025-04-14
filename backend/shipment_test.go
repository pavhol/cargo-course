package main

import (
	"bytes"
	"cargo-back/handlers"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	api := r.Group("/api")
	api.POST("/shipments/detailed", handlers.CreateDetailedShipment)
	api.GET("/shipments/:id", handlers.GetShipmentByID)
	return r
}

func TestCreateDetailedShipment(t *testing.T) {
	router := setupRouter()

	payload := map[string]interface{}{
		"route_id":   1,
		"start_date": time.Now().Format("2006-01-02"),
		"end_date":   time.Now().AddDate(0, 0, 2).Format("2006-01-02"),
		"bonus":      120.5,
		"driver_id":  1,
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/api/shipments/detailed", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Ожидался статус 200, получен %d. Ответ: %s", w.Code, w.Body.String())
	}
}

func TestGetShipmentByID(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/shipments/1", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}
