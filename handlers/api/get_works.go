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

		rows, err := db.Query(`
			SELECT 
				w.id,
				w.title,
				w.slug,
				w.at_project,
				w.created_at,
				COALESCE(l.likes_count, 0) as likes_count,
				COALESCE(v.views_count, 0) as views_count
			FROM works w
			LEFT JOIN (
				SELECT element_id, COUNT(*) as likes_count
				FROM likes
				WHERE element_type = 'work'
				GROUP BY element_id
			) l ON l.element_id = w.id
			LEFT JOIN (
				SELECT element_id, COUNT(*) as views_count
				FROM views
				WHERE element_type = 'work'
				GROUP BY element_id
			) v ON v.element_id = w.id
			WHERE w.at_project = ?
			ORDER BY w.created_at DESC
		`, projectID)

		if err != nil {
			log.Println("Ошибка получения работ:", err)
			respondWithError(w, http.StatusInternalServerError, "Ошибка при получении работ")
			return
		}
		defer rows.Close()

		works := []WorkListItem{}
		for rows.Next() {
			var work WorkListItem
			err := rows.Scan(
				&work.ID,
				&work.Title,
				&work.Slug,
				&work.ProjectID,
				&work.CreatedAt,
				&work.Likes,
				&work.Views,
			)
			if err != nil {
				log.Println("Ошибка сканирования работы:", err)
				continue
			}
			works = append(works, work)
		}

		respondWithJSON(w, http.StatusOK, works)
	}
}
