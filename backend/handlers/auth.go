package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"cargo-back/database"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Структура данных для логина
type LoginData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Строка для подписи токена
var secretKey = "my_secret_key"

// Обработчик для логина
func Login(c *gin.Context) {
	var loginData LoginData

	// Получаем данные из тела запроса
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Подключаемся к базе данных
	db, err := database.Connect()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка подключения к базе данных"})
		return
	}
	defer db.Close()

	// Получаем данные пользователя из БД
	var id, role_id int
	var username, hashedPassword string
	err = db.QueryRow("SELECT id, username, password_hash, role_id FROM users WHERE username = ?", loginData.Username).
		Scan(&id, &username, &hashedPassword, &role_id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	// Сравниваем захешированный пароль из БД с введенным паролем
	if !checkPasswordHash(loginData.Password, hashedPassword) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Создаем JWT токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":   tokenString,
		"role_id": role_id})
}

// Функция для сравнения пароля с хешем (используя bcrypt)
func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Println("bcrypt error:", err)
		return false
	}
	return true
}
