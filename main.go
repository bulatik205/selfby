package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"selfby/config"
	new_projects "selfby/handlers/api"
	"selfby/handlers/auth"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	cfg := config.LoadConfig()

	var err error
	db, err = sql.Open("mysql", cfg.GetDSN())
	if err != nil {
		log.Fatal("Ошибка подключения к БД:", err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("Не могу подключиться к MySQL:", err)
	}

	var authPath = "dashboard"
	var bladeDir = "templates"

	http.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("styles"))))
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, bladeDir+"/index.html")
	})

	http.HandleFunc("/reg", func(w http.ResponseWriter, r *http.Request) {
		if checkSession(w, r) {
			http.Redirect(w, r, authPath, http.StatusSeeOther)
		} else {
			http.ServeFile(w, r, bladeDir+"/reg.html")
		}
	})

	http.HandleFunc("/auth/reg", auth.RegisterUser(db))

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if checkSession(w, r) {
			http.Redirect(w, r, authPath, http.StatusSeeOther)
		} else {
			http.ServeFile(w, r, bladeDir+"/login.html")
		}
	})

	http.HandleFunc("/auth/login", auth.LoginUser(db))

	http.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		if checkSession(w, r) {
			http.ServeFile(w, r, bladeDir+"/dashboard.html")
		} else {
			http.Redirect(w, r, "/reg", http.StatusSeeOther)
		}
	})

	http.HandleFunc("/api/v1/newProject", new_projects.NewProject(db))

	fmt.Printf("Сервер запущен на http://localhost:%s\n", cfg.ServerPort)
	http.ListenAndServe(":"+cfg.ServerPort, nil)
}

func checkSession(w http.ResponseWriter, r *http.Request) bool {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return false
	}

	var email string
	var sessionExpires sql.NullTime

	err = db.QueryRow(
		"SELECT email, session_expires FROM users WHERE session_token = ?",
		cookie.Value,
	).Scan(&email, &sessionExpires)

	if err != nil || !sessionExpires.Valid || sessionExpires.Time.Before(time.Now()) {
		deleteSessionCookie(w)
		clearSessionToken(cookie.Value)
		return false
	}

	return true
}

func deleteSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
}

func clearSessionToken(token string) {
	db.Exec("UPDATE users SET session_token = NULL, session_expires = NULL WHERE session_token = ?", token)
}
