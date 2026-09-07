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
	Likes     int       `json:"likes"`
	Views     int       `json:"views"`
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
			SELECT 
				w.id,
				w.title,
				w.slug,
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
				&work.CreatedAt,
				&work.Likes,
				&work.Views,
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
