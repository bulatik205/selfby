package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func GetCurrentUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var user User
		var sessionExpires sql.NullTime

		err = db.QueryRow(`
			SELECT id, username, email, created_at, session_expires 
			FROM users 
			WHERE session_token = ?
		`, cookie.Value).Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.CreatedAt,
			&sessionExpires,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			log.Println("Ошибка получения пользователя:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if !sessionExpires.Valid || sessionExpires.Time.Before(time.Now()) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	}
}
