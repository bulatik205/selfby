package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
)

type ViewRequest struct {
	ElementID   int64  `json:"element_id"`
	ElementType string `json:"element_type"`
}

type ViewResponse struct {
	Viewed bool `json:"viewed"`
	Count  int  `json:"count"`
}

func AddView(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		userID, _ := getUserIDFromSession(db, r)

		var req ViewRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Некорректные данные")
			return
		}

		if req.ElementID == 0 {
			respondWithError(w, http.StatusBadRequest, "Не указан element_id")
			return
		}

		if req.ElementType == "" {
			req.ElementType = "work"
		}

		var count int
		err := db.QueryRow(
			"SELECT COUNT(*) FROM views WHERE element_id = ? AND element_type = ?",
			req.ElementID, req.ElementType,
		).Scan(&count)
		if err != nil {
			log.Println("Ошибка подсчета просмотров:", err)
			count = 0
		}

		var lastView sql.NullTime
		err = db.QueryRow(
			"SELECT created_at FROM views WHERE element_id = ? AND element_type = ? AND owner_id = ? ORDER BY created_at DESC LIMIT 1",
			req.ElementID, req.ElementType, userID,
		).Scan(&lastView)

		if err == nil && lastView.Valid {
			if time.Since(lastView.Time) < time.Hour {
				respondWithJSON(w, http.StatusOK, ViewResponse{
					Viewed: false,
					Count:  count,
				})
				return
			}
		}

		_, err = db.Exec(
			"INSERT INTO views (element_id, element_type, owner_id) VALUES (?, ?, ?)",
			req.ElementID, req.ElementType, userID,
		)
		if err != nil {
			if isDuplicateError(err) {
				respondWithJSON(w, http.StatusOK, ViewResponse{
					Viewed: false,
					Count:  count,
				})
				return
			}
			log.Println("Ошибка добавления просмотра:", err)
			respondWithError(w, http.StatusInternalServerError, "Ошибка сервера")
			return
		}

		count++

		respondWithJSON(w, http.StatusOK, ViewResponse{
			Viewed: true,
			Count:  count,
		})
	}
}

func GetViews(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		elementID := r.URL.Query().Get("element_id")
		elementType := r.URL.Query().Get("element_type")

		if elementType == "" {
			elementType = "work"
		}

		if elementID == "" {
			respondWithError(w, http.StatusBadRequest, "Не указан element_id")
			return
		}

		var count int
		err := db.QueryRow(
			"SELECT COUNT(*) FROM views WHERE element_id = ? AND element_type = ?",
			elementID, elementType,
		).Scan(&count)
		if err != nil {
			log.Println("Ошибка подсчета просмотров:", err)
			count = 0
		}

		respondWithJSON(w, http.StatusOK, ViewResponse{
			Viewed: false,
			Count:  count,
		})
	}
}

func isDuplicateError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "Error 1062") || strings.Contains(err.Error(), "Duplicate entry")
}
