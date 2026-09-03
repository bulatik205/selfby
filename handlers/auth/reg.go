package auth

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		if err := r.ParseForm(); err != nil {
			RedirectWithError(w, r, "Ошибка при обработке формы")
			return
		}

		email := strings.ToLower(strings.TrimSpace(r.FormValue("email")))
		password := r.FormValue("password")

		if errMsg := validateRegistration(email, password); errMsg != "" {
			RedirectWithError(w, r, errMsg)
			return
		}

		exists, err := checkUserExists(db, email)
		if err != nil {
			log.Println("Ошибка проверки пользователя:", err)
			RedirectWithError(w, r, "Ошибка сервера")
			return
		}
		if exists {
			RedirectWithError(w, r, "Пользователь с таким email уже существует")
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Ошибка хеширования пароля:", err)
			RedirectWithError(w, r, "Ошибка сервера")
			return
		}

		userID, err := saveUser(db, email, string(hashedPassword))
		if err != nil {
			log.Println("Ошибка сохранения пользователя:", err)
			RedirectWithError(w, r, "Ошибка при регистрации")
			return
		}

		sessionToken := generateSessionToken()
		if err := saveSessionToken(db, userID, sessionToken); err != nil {
			log.Println("Ошибка сохранения сессии:", err)
			RedirectWithError(w, r, "Ошибка при создании сессии")
			return
		}

		setSessionCookie(w, sessionToken)

		log.Printf("Зарегистрирован новый пользователь: %s (ID: %d)", email, userID)
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}
}

func generateSessionToken() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func saveSessionToken(db *sql.DB, userID int64, token string) error {
	expiresAt := time.Now().Add(24 * time.Hour)
	_, err := db.Exec(
		"UPDATE users SET session_token = ?, session_expires = ? WHERE id = ?",
		token, expiresAt, userID,
	)
	return err
}

func setSessionCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   86400,
	})
}

func RedirectWithError(w http.ResponseWriter, r *http.Request, message string) {
	http.Redirect(w, r, "/reg?error="+url.QueryEscape(message), http.StatusSeeOther)
}

func validateRegistration(email, password string) string {
	if email == "" || password == "" {
		return "Все поля обязательны для заполнения"
	}
	if !isValidEmail(email) {
		return "Некорректный email адрес"
	}
	if len(password) < 6 {
		return "Пароль должен быть не менее 6 символов"
	}
	return ""
}

func isValidEmail(email string) bool {
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return false
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return false
	}
	return true
}

func checkUserExists(db *sql.DB, email string) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ?", email).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func saveUser(db *sql.DB, email, passwordHash string) (int64, error) {
	result, err := db.Exec(
		"INSERT INTO users (email, password_hash) VALUES (?, ?)",
		email, passwordHash,
	)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
