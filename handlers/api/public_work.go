package api

import (
	"database/sql"
	"log"
	"net/http"
	"time"
)

type PublicWorkDetail struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	ContentHTML string    `json:"content_html"`
	Views       int       `json:"views"`
	Likes       int       `json:"likes"`
	CreatedAt   time.Time `json:"created_at"`
	ProjectID   int64     `json:"project_id"`
	ProjectName string    `json:"project_name"`
	ProjectType string    `json:"project_type"`
	OwnerID     int64     `json:"owner_id"`
	OwnerName   string    `json:"owner_name"`
}

func GetPublicWork(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		username := r.URL.Query().Get("username")
		projectName := r.URL.Query().Get("project")
		workSlug := r.URL.Query().Get("slug")

		if username == "" || projectName == "" || workSlug == "" {
			respondWithError(w, http.StatusBadRequest, "Не указаны параметры")
			return
		}

		var work PublicWorkDetail
		err := db.QueryRow(`
			SELECT 
				w.id,
				w.title,
				w.slug,
				w.content_html,
				w.created_at,
				p.id as project_id,
				p.name as project_name,
				p.type as project_type,
				p.owner_id,
				u.username as owner_name,
				COALESCE(l.likes_count, 0) as likes_count,
				COALESCE(v.views_count, 0) as views_count
			FROM works w
			JOIN projects p ON w.at_project = p.id
			JOIN users u ON p.owner_id = u.id
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
			WHERE u.username = ? AND p.name = ? AND w.slug = ?
		`, username, projectName, workSlug).Scan(
			&work.ID,
			&work.Title,
			&work.Slug,
			&work.ContentHTML,
			&work.CreatedAt,
			&work.ProjectID,
			&work.ProjectName,
			&work.ProjectType,
			&work.OwnerID,
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

		if work.ProjectType == "private" {
			userID, err := getUserIDFromSession(db, r)
			if err != nil || userID != work.OwnerID {
				respondWithError(w, http.StatusForbidden, "Это приватный проект")
				return
			}
		}

		respondWithJSON(w, http.StatusOK, work)
	}
}
