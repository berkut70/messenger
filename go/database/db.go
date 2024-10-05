package database

import (
	"database/sql"
	"fmt"

	"main/models"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

func ConnectDB() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:Desk3745@tcp(localhost:3306)/messenger")
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("ошибка ping: %w", err)
	}

	return db, nil
}

func GetUser(db *sql.DB, id int) (*models.User, error) {
	var user models.User
	err := db.QueryRow("SELECT id, username, password FROM users WHERE id = ?", id).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения пользователя: %w", err)
	}
	return &user, nil
}

func GetUserByUsername(db *sql.DB, username string) (*models.User, error) {
	var user models.User
	err := db.QueryRow("SELECT id, username, password FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения пользователя: %w", err)
	}
	return &user, nil
}

func CheckPassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func CreateUser(db *sql.DB, user models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("ошибка хеширования пароля: %w", err)
	}

	_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", user.Username, string(hashedPassword))
	if err != nil {
		return fmt.Errorf("ошибка создания пользователя: %w", err)
	}
	return nil
}

func SendMessage(db *sql.DB, senderID int, receiverID int, text string) error {
	_, err := db.Exec("INSERT INTO messages (sender_id, receiver_id, text) VALUES (?, ?, ?)", senderID, receiverID, text)
	if err != nil {
		return fmt.Errorf("ошибка отправки сообщения: %w", err)
	}
	return nil
}
