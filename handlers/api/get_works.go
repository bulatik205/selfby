package api

import (
	"database/sql"
	"log"
	"net/http"
	"time"
)

type WorkListItem struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	ProjectID int64     `json:"project_id"`
	Views     int       `json:"views"`
	Likes     int       `json:"likes"`
	CreatedAt time.Time `json:"created_at"`
}

func GetWorks(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		userID, err := getUserIDFromSession(db, r)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Необходима авторизация")
			return
		}

		projectName := r.URL.Query().Get("project")
		if projectName == "" {
			respondWithError(w, http.StatusBadRequest, "Не указано имя проекта")
			return
		}

		var projectID int64
		err = db.QueryRow(
			"SELECT id FROM projects WHERE owner_id = ? AND name = ?",
			userID, projectName,
		).Scan(&projectID)

		if err != nil {
			if err == sql.ErrNoRows {
				respondWithError(w, http.StatusNotFound, "Проект не найден")
			} else {
				log.Println("Ошибка получения проекта:", err)
				respondWithError(w, http.StatusInternalServerError, "Ошибка сервера")
			}
			return
		}

		rows, err := db.Query(
			"SELECT id, title, slug, at_project, views, likes, created_at FROM works WHERE at_project = ? ORDER BY created_at DESC",
			projectID,
		)
		if err != nil {
			log.Println("Ошибка получения работ:", err)
			respondWithError(w, http.StatusInternalServerError, "Ошибка при получении работ")
			return
		}
		defer rows.Close()

		var works []WorkListItem
		for rows.Next() {
			var work WorkListItem
			err := rows.Scan(
				&work.ID,
				&work.Title,
				&work.Slug,
				&work.ProjectID,
				&work.Views,
				&work.Likes,
				&work.CreatedAt,
			)
			if err != nil {
				log.Println("Ошибка сканирования работы:", err)
				continue
			}
			works = append(works, work)
		}

		if works == nil {
			works = []WorkListItem{}
		}

		respondWithJSON(w, http.StatusOK, works)
	}
}
