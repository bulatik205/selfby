package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"
)

type TopUser struct {
	ID            int64     `json:"id"`
	Username      string    `json:"username"`
	CreatedAt     time.Time `json:"created_at"`
	ProjectsCount int       `json:"projects_count"`
	WorksCount    int       `json:"works_count"`
	TotalLikes    int       `json:"total_likes"`
	TotalViews    int       `json:"total_views"`
}

func GetTopUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		sortBy := r.URL.Query().Get("sort")
		if sortBy == "" {
			sortBy = "views"
		}

		var orderBy string
		switch sortBy {
		case "likes":
			orderBy = "total_likes DESC, total_views DESC"
		case "works":
			orderBy = "works_count DESC, total_views DESC"
		case "projects":
			orderBy = "projects_count DESC, total_views DESC"
		default:
			orderBy = "total_views DESC, total_likes DESC"
		}

		query := fmt.Sprintf(`
			SELECT 
				u.id,
				u.username,
				u.created_at,
				COUNT(DISTINCT p.id) as projects_count,
				COUNT(DISTINCT w.id) as works_count,
				COALESCE(SUM(l.likes_count), 0) as total_likes,
				COALESCE(SUM(v.views_count), 0) as total_views
			FROM users u
			LEFT JOIN projects p ON p.owner_id = u.id AND p.type = 'public'
			LEFT JOIN works w ON w.owner_id = u.id
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
			GROUP BY u.id, u.username, u.created_at
			ORDER BY %s
			LIMIT 50
		`, orderBy)

		rows, err := db.Query(query)
		if err != nil {
			log.Println("Ошибка получения топ пользователей:", err)
			respondWithError(w, http.StatusInternalServerError, "Ошибка сервера")
			return
		}
		defer rows.Close()

		users := []TopUser{}
		for rows.Next() {
			var user TopUser
			err := rows.Scan(
				&user.ID,
				&user.Username,
				&user.CreatedAt,
				&user.ProjectsCount,
				&user.WorksCount,
				&user.TotalLikes,
				&user.TotalViews,
			)
			if err != nil {
				log.Println("Ошибка сканирования пользователя:", err)
				continue
			}
			users = append(users, user)
		}

		respondWithJSON(w, http.StatusOK, users)
	}
}
