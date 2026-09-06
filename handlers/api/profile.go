package api

import (
	"database/sql"
	"log"
	"net/http"
	"time"
)

type UserProfile struct {
	ID              int64         `json:"id"`
	Username        string        `json:"username"`
	CreatedAt       time.Time     `json:"created_at"`
	ProjectsCount   int           `json:"projects_count"`
	PublicProjects  int           `json:"public_projects"`
	PrivateProjects int           `json:"private_projects"`
	WorksCount      int           `json:"works_count"`
	TotalLikes      int           `json:"total_likes"`
	TotalViews      int           `json:"total_views"`
	Projects        []ProjectItem `json:"projects"`
}

type ProjectItem struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	CreatedAt   time.Time `json:"created_at"`
	WorksCount  int       `json:"works_count"`
}

func GetProfile(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		username := r.URL.Query().Get("username")
		if username == "" {
			respondWithError(w, http.StatusBadRequest, "Не указано имя пользователя")
			return
		}

		// Получаем основную информацию о пользователе
		var profile UserProfile
		err := db.QueryRow(`
			SELECT id, username, created_at
			FROM users
			WHERE username = ?
		`, username).Scan(
			&profile.ID,
			&profile.Username,
			&profile.CreatedAt,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				respondWithError(w, http.StatusNotFound, "Пользователь не найден")
			} else {
				log.Println("Ошибка получения пользователя:", err)
				respondWithError(w, http.StatusInternalServerError, "Ошибка сервера")
			}
			return
		}

		// Получаем количество проектов
		err = db.QueryRow(`
			SELECT 
				COUNT(*) as total,
				COALESCE(SUM(CASE WHEN type = 'public' THEN 1 ELSE 0 END), 0) as public_count,
				COALESCE(SUM(CASE WHEN type = 'private' THEN 1 ELSE 0 END), 0) as private_count
			FROM projects
			WHERE owner_id = ?
		`, profile.ID).Scan(
			&profile.ProjectsCount,
			&profile.PublicProjects,
			&profile.PrivateProjects,
		)

		if err != nil {
			log.Println("Ошибка получения проектов:", err)
			profile.ProjectsCount = 0
			profile.PublicProjects = 0
			profile.PrivateProjects = 0
		}

		// Получаем количество работ
		err = db.QueryRow(`
			SELECT COUNT(*)
			FROM works
			WHERE owner_id = ?
		`, profile.ID).Scan(&profile.WorksCount)

		if err != nil {
			log.Println("Ошибка получения работ:", err)
			profile.WorksCount = 0
		}

		// Получаем суммарные лайки и просмотры
		err = db.QueryRow(`
			SELECT 
				COALESCE(SUM(likes), 0) as total_likes,
				COALESCE(SUM(views), 0) as total_views
			FROM works
			WHERE owner_id = ?
		`, profile.ID).Scan(
			&profile.TotalLikes,
			&profile.TotalViews,
		)

		if err != nil {
			log.Println("Ошибка получения статистики:", err)
			profile.TotalLikes = 0
			profile.TotalViews = 0
		}

		// Получаем список проектов пользователя
		rows, err := db.Query(`
			SELECT 
				p.id,
				p.name,
				p.description,
				p.type,
				p.created_at,
				COUNT(w.id) as works_count
			FROM projects p
			LEFT JOIN works w ON w.at_project = p.id
			WHERE p.owner_id = ?
			GROUP BY p.id, p.name, p.description, p.type, p.created_at
			ORDER BY p.created_at DESC
		`, profile.ID)

		if err != nil {
			log.Println("Ошибка получения списка проектов:", err)
			profile.Projects = []ProjectItem{}
		} else {
			defer rows.Close()

			profile.Projects = []ProjectItem{}
			for rows.Next() {
				var project ProjectItem
				err := rows.Scan(
					&project.ID,
					&project.Name,
					&project.Description,
					&project.Type,
					&project.CreatedAt,
					&project.WorksCount,
				)
				if err != nil {
					log.Println("Ошибка сканирования проекта:", err)
					continue
				}
				profile.Projects = append(profile.Projects, project)
			}
		}

		respondWithJSON(w, http.StatusOK, profile)
	}
}
