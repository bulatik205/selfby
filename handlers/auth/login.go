package auth

import (
	"database/sql"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func LoginUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		if err := r.ParseForm(); err != nil {
			LoginRedirectWithError(w, r, "Ошибка при обработке формы")
			return
		}

		email := strings.ToLower(strings.TrimSpace(r.FormValue("email")))
		password := r.FormValue("password")

		if email == "" || password == "" {
			LoginRedirectWithError(w, r, "Все поля обязательны для заполнения")
			return
		}

		user, err := getUserByEmail(db, email)
		if err != nil {
			if err == sql.ErrNoRows {
				LoginRedirectWithError(w, r, "Пользователь с таким email не найден")
				return
			}
			log.Println("Ошибка получения пользователя:", err)
			LoginRedirectWithError(w, r, "Ошибка сервера")
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
			LoginRedirectWithError(w, r, "Неверный пароль")
			return
		}

		sessionToken := generateSessionToken()
		if err := saveSessionToken(db, user.ID, sessionToken); err != nil {
			log.Println("Ошибка сохранения сессии:", err)
			LoginRedirectWithError(w, r, "Ошибка при создании сессии")
			return
		}

		setSessionCookie(w, sessionToken)

		log.Printf("Пользователь вошел: %s (ID: %d)", email, user.ID)
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}

type User struct {
	ID           int64
	Email        string
	PasswordHash string
}

func getUserByEmail(db *sql.DB, email string) (*User, error) {
	user := &User{}
	err := db.QueryRow(
		"SELECT id, email, password_hash FROM users WHERE email = ?",
		email,
	).Scan(&user.ID, &user.Email, &user.PasswordHash)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func LoginPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/login.html")
}

func LoginRedirectWithError(w http.ResponseWriter, r *http.Request, message string) {
	http.Redirect(w, r, "/login?error="+url.QueryEscape(message), http.StatusSeeOther)
}
