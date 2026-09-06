package api

import (
	"database/sql"
	"log"
	"net/http"
	"time"
)

type PublicProject struct {
	ID          int64        `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Type        string       `json:"type"`
	CreatedAt   time.Time    `json:"created_at"`
	OwnerID     int64        `json:"owner_id"`
	OwnerName   string       `json:"owner_name"`
	Works       []PublicWork `json:"works"`
}

type PublicWork struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	Views     int       `json:"views"`
	Likes     int       `json:"likes"`
	CreatedAt time.Time `json:"created_at"`
}

func GetPublicProject(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		username := r.URL.Query().Get("username")
		projectName := r.URL.Query().Get("project")

		if username == "" || projectName == "" {
			respondWithError(w, http.StatusBadRequest, "Не указаны параметры")
			return
		}

		var project PublicProject
		err := db.QueryRow(`
			SELECT 
				p.id,
				p.name,
				p.description,
				p.type,
				p.created_at,
				p.owner_id,
				u.username as owner_name
			FROM projects p
			JOIN users u ON p.owner_id = u.id
			WHERE u.username = ? AND p.name = ?
		`, username, projectName).Scan(
			&project.ID,
			&project.Name,
			&project.Description,
			&project.Type,
			&project.CreatedAt,
			&project.OwnerID,
			&project.OwnerName,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				respondWithError(w, http.StatusNotFound, "Проект не найден")
			} else {
				log.Println("Ошибка получения проекта:", err)
				respondWithError(w, http.StatusInternalServerError, "Ошибка сервера")
			}
			return
		}

		if project.Type == "private" {
			userID, err := getUserIDFromSession(db, r)
			if err != nil || userID != project.OwnerID {
				respondWithError(w, http.StatusForbidden, "Это приватный проект")
				return
			}
		}

		rows, err := db.Query(`
			SELECT id, title, slug, views, likes, created_at
			FROM works
			WHERE at_project = ?
			ORDER BY created_at DESC
		`, project.ID)

		if err != nil {
			log.Println("Ошибка получения работ:", err)
			respondWithError(w, http.StatusInternalServerError, "Ошибка сервера")
			return
		}
		defer rows.Close()

		project.Works = []PublicWork{}
		for rows.Next() {
			var work PublicWork
			err := rows.Scan(
				&work.ID,
				&work.Title,
				&work.Slug,
				&work.Views,
				&work.Likes,
				&work.CreatedAt,
			)
			if err != nil {
				log.Println("Ошибка сканирования работы:", err)
				continue
			}
			project.Works = append(project.Works, work)
		}

		respondWithJSON(w, http.StatusOK, project)
	}
}
