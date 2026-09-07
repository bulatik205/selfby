package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"
)

type ProjectListItem struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	CreatedAt   time.Time `json:"created_at"`
	WorksCount  int       `json:"works_count"`
	TotalViews  int       `json:"total_views"`
	TotalLikes  int       `json:"total_likes"`
}

func GetProjectsWithStats(db *sql.DB) http.HandlerFunc {
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

		sortBy := r.URL.Query().Get("sort")
		if sortBy == "" {
			sortBy = "created"
		}

		var orderBy string
		switch sortBy {
		case "views":
			orderBy = "total_views DESC, p.created_at DESC"
		case "likes":
			orderBy = "total_likes DESC, p.created_at DESC"
		case "works":
			orderBy = "works_count DESC, p.created_at DESC"
		case "name":
			orderBy = "p.name ASC"
		default:
			orderBy = "p.created_at DESC"
		}

		query := fmt.Sprintf(`
			SELECT 
				p.id,
				p.name,
				p.description,
				p.type,
				p.created_at,
				COUNT(DISTINCT w.id) as works_count,
				COALESCE(SUM(v.views_count), 0) as total_views,
				COALESCE(SUM(l.likes_count), 0) as total_likes
			FROM projects p
			LEFT JOIN works w ON w.at_project = p.id
			LEFT JOIN (
				SELECT element_id, COUNT(*) as views_count
				FROM views
				WHERE element_type = 'work'
				GROUP BY element_id
			) v ON v.element_id = w.id
			LEFT JOIN (
				SELECT element_id, COUNT(*) as likes_count
				FROM likes
				WHERE element_type = 'work'
				GROUP BY element_id
			) l ON l.element_id = w.id
			WHERE p.owner_id = ?
			GROUP BY p.id, p.name, p.description, p.type, p.created_at
			ORDER BY %s
		`, orderBy)

		rows, err := db.Query(query, userID)
		if err != nil {
			log.Println("Ошибка получения проектов:", err)
			respondWithError(w, http.StatusInternalServerError, "Ошибка сервера")
			return
		}
		defer rows.Close()

		projects := []ProjectListItem{}
		for rows.Next() {
			var project ProjectListItem
			err := rows.Scan(
				&project.ID,
				&project.Name,
				&project.Description,
				&project.Type,
				&project.CreatedAt,
				&project.WorksCount,
				&project.TotalViews,
				&project.TotalLikes,
			)
			if err != nil {
				log.Println("Ошибка сканирования проекта:", err)
				continue
			}
			projects = append(projects, project)
		}

		respondWithJSON(w, http.StatusOK, projects)
	}
}
