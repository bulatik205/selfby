package api

import (
	"database/sql"
	"log"
	"net/http"
	"time"
)

type WorkDetail struct {
	ID          int64     `json:"id"`
	OwnerID     int64     `json:"owner_id"`
	ProjectID   int64     `json:"project_id"`
	ProjectName string    `json:"project_name"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	ContentMD   string    `json:"content_md"`
	ContentHTML string    `json:"content_html"`
	Likes       int       `json:"likes"`
	Views       int       `json:"views"`
	CreatedAt   time.Time `json:"created_at"`
	OwnerName   string    `json:"owner_name"`
}

func GetWork(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		projectName := r.URL.Query().Get("project")
		workSlug := r.URL.Query().Get("slug")

		if projectName == "" || workSlug == "" {
			respondWithError(w, http.StatusBadRequest, "Не указаны параметры")
			return
		}

		var work WorkDetail
		err := db.QueryRow(`
			SELECT 
				w.id,
				w.owner_id,
				w.at_project,
				p.name as project_name,
				w.title,
				w.slug,
				w.content_md,
				w.content_html,
				w.created_at,
				u.username as owner_name,
				COALESCE(l.likes_count, 0) as likes_count,
				COALESCE(v.views_count, 0) as views_count
			FROM works w
			JOIN projects p ON w.at_project = p.id
			JOIN users u ON w.owner_id = u.id
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
			WHERE p.name = ? AND w.slug = ?
		`, projectName, workSlug).Scan(
			&work.ID,
			&work.OwnerID,
			&work.ProjectID,
			&work.ProjectName,
			&work.Title,
			&work.Slug,
			&work.ContentMD,
			&work.ContentHTML,
			&work.CreatedAt,
			&work.OwnerName,
			&work.Likes,
			&work.Views,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				respondWithError(w, http.StatusNotFound, "Работа не найдена")
			} else {
				log.Println("Ошибка получения работы:", err)
				respondWithError(w, http.StatusInternalServerError, "Ошибка сервера")
			}
			return
		}

		respondWithJSON(w, http.StatusOK, work)
	}
}
