package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"main/database"
	"main/models"
	"main/utils"
)

func HandleUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetUser(w, r, db)
		case http.MethodPost:
			handleCreateUser(w, r, db)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func handleGetUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	user, err := database.GetUser(db, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func handleCreateUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.CreateUser(db, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func HandleLogin(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleLogin(w, r, db)
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request, db *sql.DB) { // Принимает db
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dbUser, err := database.GetUserByUsername(db, user.Username)
	if err != nil {
		log.Printf("Ошибка при получении пользователя: %v", err)
		http.Error(w, "Неправильный логин или пароль", http.StatusUnauthorized)
		return
	}

	if !database.CheckPassword(dbUser.Password, user.Password) {
		log.Printf("Неверный пароль для пользователя %s", user.Username)
		http.Error(w, "Неправильный логин или пароль", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateJWT(dbUser.ID, dbUser.Username) // Добавьте "" для role
	if err != nil {
		http.Error(w, "Ошибка при генерации токена", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", token)
	json.NewEncoder(w).Encode(map[string]string{"message": "Успешная авторизация"})
}
