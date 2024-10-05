package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"main/database"
	"main/handlers"
	"main/utils"

	"github.com/gorilla/mux"
)

func main() {
	db, err := database.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := mux.NewRouter()

	router.Use(jwtMiddleware)

	router.HandleFunc("/users", handlers.HandleUsers(db)).Methods("GET", "POST")
	router.HandleFunc("/users/login", handlers.HandleLogin(db)).Methods("POST")
	router.HandleFunc("/messages", handlers.HandleMessages(db)).Methods("POST")

	fmt.Println("Сервер запущен на порту 8080")
	log.Fatal(http.ListenAndServe("89.169.172.158:8080", router))
}

func jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Необходим токен авторизации", http.StatusUnauthorized)
			return
		}

		claims, err := utils.ParseJWT(tokenString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), utils.UserIDKey{}, claims.UserID)
		ctx = context.WithValue(ctx, utils.UsernameKey{}, claims.Username)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
