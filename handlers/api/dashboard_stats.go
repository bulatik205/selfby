package api

import (
	"database/sql"
	"log"
	"net/http"
	"time"
)

type RecentLike struct {
	ID          int64     `json:"id"`
	ElementType string    `json:"element_type"`
	WorkTitle   string    `json:"work_title"`
	WorkSlug    string    `json:"work_slug"`
	ProjectName string    `json:"project_name"`
	LikerName   string    `json:"liker_name"`
	CreatedAt   time.Time `json:"created_at"`
}

type ProjectStatsItem struct {
	ID         int64  `json:"id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	WorksCount int    `json:"works_count"`
	TotalViews int    `json:"total_views"`
	TotalLikes int    `json:"total_likes"`
}

type DashboardStats struct {
	RecentLikes []RecentLike       `json:"recent_likes"`
	TopProjects []ProjectStatsItem `json:"top_projects"`
}

func GetDashboardStats(db *sql.DB) http.HandlerFunc {
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

		stats := DashboardStats{
			RecentLikes: []RecentLike{},
			TopProjects: []ProjectStatsItem{},
		}

		rows, err := db.Query(`
			SELECT 
				l.id,
				l.element_type,
				COALESCE(w.title, '') as work_title,
				COALESCE(w.slug, '') as work_slug,
				p.name as project_name,
				u.username as liker_name,
				l.created_at
			FROM likes l
			JOIN users u ON l.owner_id = u.id
			LEFT JOIN works w ON l.element_type = 'work' AND l.element_id = w.id
			JOIN projects p ON (
				(l.element_type = 'work' AND w.at_project = p.id) OR
				(l.element_type = 'project' AND l.element_id = p.id)
			)
			WHERE p.owner_id = ?
			ORDER BY l.created_at DESC
			LIMIT 10
		`, userID)

		if err != nil {
			log.Println("Ошибка получения последних лайков:", err)
		} else {
			defer rows.Close()
			for rows.Next() {
				var like RecentLike
				err := rows.Scan(
					&like.ID,
					&like.ElementType,
					&like.WorkTitle,
					&like.WorkSlug,
					&like.ProjectName,
					&like.LikerName,
					&like.CreatedAt,
				)
				if err != nil {
					log.Println("Ошибка сканирования лайка:", err)
					continue
				}
				stats.RecentLikes = append(stats.RecentLikes, like)
			}
		}

		rows2, err := db.Query(`
			SELECT 
				p.id,
				p.name,
				p.type,
				COUNT(DISTINCT w.id) as works_count,
				COALESCE((
					SELECT COUNT(*) FROM views 
					WHERE element_type = 'work' AND element_id IN (SELECT id FROM works WHERE at_project = p.id)
				), 0) + COALESCE((
					SELECT COUNT(*) FROM views 
					WHERE element_type = 'project' AND element_id = p.id
				), 0) as total_views,
				COALESCE((
					SELECT COUNT(*) FROM likes 
					WHERE element_type = 'work' AND element_id IN (SELECT id FROM works WHERE at_project = p.id)
				), 0) + COALESCE((
					SELECT COUNT(*) FROM likes 
					WHERE element_type = 'project' AND element_id = p.id
				), 0) as total_likes
			FROM projects p
			LEFT JOIN works w ON w.at_project = p.id
			WHERE p.owner_id = ?
			GROUP BY p.id, p.name, p.type
			ORDER BY total_views DESC, total_likes DESC
			LIMIT 5
		`, userID)

		if err != nil {
			log.Println("Ошибка получения топ проектов:", err)
		} else {
			defer rows2.Close()
			for rows2.Next() {
				var proj ProjectStatsItem
				err := rows2.Scan(
					&proj.ID,
					&proj.Name,
					&proj.Type,
					&proj.WorksCount,
					&proj.TotalViews,
					&proj.TotalLikes,
				)
				if err != nil {
					log.Println("Ошибка сканирования проекта:", err)
					continue
				}
				stats.TopProjects = append(stats.TopProjects, proj)
			}
		}

		respondWithJSON(w, http.StatusOK, stats)
	}
}
