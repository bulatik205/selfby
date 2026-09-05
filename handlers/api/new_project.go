package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type ProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

type ProjectResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewProject(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		userID, err := getUserIDFromSession(db, r)
		if err != nil {
			log.Println("Ошибка получения userID из сессии:", err)
			respondWithError(w, http.StatusUnauthorized, "Необходима авторизация")
			return
		}

		var req ProjectRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Println("Ошибка декодирования JSON:", err)
			respondWithError(w, http.StatusBadRequest, "Некорректные данные")
			return
		}

		req.Name = strings.TrimSpace(req.Name)
		req.Description = strings.TrimSpace(req.Description)
		req.Type = strings.TrimSpace(req.Type)

		if errMsg := validateProject(req); errMsg != "" {
			respondWithError(w, http.StatusBadRequest, errMsg)
			return
		}

		exists, err := checkProjectNameExists(db, req.Name)
		if err != nil {
			log.Println("Ошибка проверки уникальности:", err)
			respondWithError(w, http.StatusInternalServerError, "Ошибка сервера")
			return
		}
		if exists {
			respondWithError(w, http.StatusConflict, "Проект с таким названием уже существует")
			return
		}

		projectID, err := saveProject(db, userID, req)
		if err != nil {
			log.Println("Ошибка сохранения проекта:", err)
			respondWithError(w, http.StatusInternalServerError, "Ошибка при создании проекта")
			return
		}

		project, err := getProjectByID(db, projectID)
		if err != nil {
			log.Println("Ошибка получения проекта:", err)
			respondWithError(w, http.StatusInternalServerError, "Ошибка при получении проекта")
			return
		}

		log.Printf("Создан новый проект: %s (ID: %d) пользователем ID: %d", project.Name, project.ID, userID)
		respondWithJSON(w, http.StatusCreated, project)
	}
}

func getUserIDFromSession(db *sql.DB, r *http.Request) (int64, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return 0, err
	}

	var userID int64
	var sessionExpires sql.NullTime

	err = db.QueryRow(
		"SELECT id, session_expires FROM users WHERE session_token = ?",
		cookie.Value,
	).Scan(&userID, &sessionExpires)

	if err != nil {
		return 0, err
	}

	if !sessionExpires.Valid || sessionExpires.Time.Before(time.Now()) {
		return 0, sql.ErrNoRows
	}

	return userID, nil
}

func validateProject(req ProjectRequest) string {
	if req.Name == "" {
		return "Название проекта обязательно"
	}

	if len(req.Name) < 3 {
		return "Название должно быть не менее 3 символов"
	}

	if len(req.Name) > 50 {
		return "Название должно быть не более 50 символов"
	}

	validNameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validNameRegex.MatchString(req.Name) {
		return "Название может содержать только английские буквы, цифры, дефисы и подчеркивания (без пробелов)"
	}

	if req.Type == "" {
		return "Тип проекта обязателен"
	}

	if req.Type != "private" && req.Type != "public" {
		return "Тип проекта может быть только 'private' или 'public'"
	}

	return ""
}

func checkProjectNameExists(db *sql.DB, name string) (bool, error) {
	var count int
	err := db.QueryRow(
		"SELECT COUNT(*) FROM projects WHERE LOWER(name) = LOWER(?)",
		name,
	).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func saveProject(db *sql.DB, userID int64, req ProjectRequest) (int64, error) {
	result, err := db.Exec(
		"INSERT INTO projects (owner_id, name, description, type) VALUES (?, ?, ?, ?)",
		userID, req.Name, req.Description, req.Type,
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

func getProjectByID(db *sql.DB, projectID int64) (*ProjectResponse, error) {
	project := &ProjectResponse{}
	err := db.QueryRow(
		"SELECT id, name, description, type, created_at FROM projects WHERE id = ?",
		projectID,
	).Scan(&project.ID, &project.Name, &project.Description, &project.Type, &project.CreatedAt)

	if err != nil {
		return nil, err
	}

	return project, nil
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Println("Ошибка маршалинга JSON:", err)
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(code)
	w.Write(response)
}
