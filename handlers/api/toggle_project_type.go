package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type ToggleProjectTypeRequest struct {
	ProjectName string `json:"project_name"`
}

type ToggleProjectTypeResponse struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func ToggleProjectType(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		userID, err := getUserIDFromSession(db, r)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Необходима авторизация")
			return
		}

		var req ToggleProjectTypeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Некорректные данные")
			return
		}

		if req.ProjectName == "" {
			respondWithError(w, http.StatusBadRequest, "Не указано имя проекта")
			return
		}

		var currentType string
		err = db.QueryRow(
			"SELECT type FROM projects WHERE owner_id = ? AND name = ?",
			userID, req.ProjectName,
		).Scan(&currentType)

		if err != nil {
			if err == sql.ErrNoRows {
				respondWithError(w, http.StatusNotFound, "Проект не найден")
			} else {
				log.Println("Ошибка получения проекта:", err)
				respondWithError(w, http.StatusInternalServerError, "Ошибка сервера")
			}
			return
		}

		newType := "public"
		if currentType == "public" {
			newType = "private"
		}

		_, err = db.Exec(
			"UPDATE projects SET type = ? WHERE owner_id = ? AND name = ?",
			newType, userID, req.ProjectName,
		)
		if err != nil {
			log.Println("Ошибка обновления типа:", err)
			respondWithError(w, http.StatusInternalServerError, "Ошибка сервера")
			return
		}

		message := "Проект теперь открытый"
		if newType == "private" {
			message = "Проект теперь закрытый"
		}

		log.Printf("Изменен тип проекта: %s -> %s пользователем ID: %d",
			req.ProjectName, newType, userID)

		respondWithJSON(w, http.StatusOK, ToggleProjectTypeResponse{
			Type:    newType,
			Message: message,
		})
	}
}
