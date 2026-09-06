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
				w.views,
				w.likes,
				w.created_at,
				p.id as project_id,
				p.name as project_name,
				p.type as project_type,
				p.owner_id,
				u.username as owner_name
			FROM works w
			JOIN projects p ON w.at_project = p.id
			JOIN users u ON p.owner_id = u.id
			WHERE u.username = ? AND p.name = ? AND w.slug = ?
		`, username, projectName, workSlug).Scan(
			&work.ID,
			&work.Title,
			&work.Slug,
			&work.ContentHTML,
			&work.Views,
			&work.Likes,
			&work.CreatedAt,
			&work.ProjectID,
			&work.ProjectName,
			&work.ProjectType,
			&work.OwnerID,
			&work.OwnerName,
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
