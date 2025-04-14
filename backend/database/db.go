package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Конфигурация (лучше вынести в отдельный config-файл)
const (
	DBUser     = "cargo_user"
	DBPassword = "secure_password"
	DBName     = "cargo_company"
)

var db *sql.DB

// Connect инициализирует подключение к БД
func Connect() (*sql.DB, error) {
	dsn := "cargo_user:secure_password@tcp(127.0.0.1:3306)/cargo_company?parseTime=true"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения: %v", err)
	}

	// Настройки пула соединений
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Повторная попытка подключения (3 попытки)
	for i := 0; i < 3; i++ {
		err = db.Ping()
		if err == nil {
			return db, nil
		}
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("не удалось подключиться после 3 попыток: %v", err)
}

// Close закрывает подключение
func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
