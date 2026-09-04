package api

import (
	"database/sql"
	"log"
	"net/http"
)

func GetProjects(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		userID, err := getUserIDFromSession(db, r)
		if err != nil {
			log.Println("Ошибка получения userID из сессии:", err)
			respondWithError(w, http.StatusUnauthorized, "Необходима авторизация")
			return
		}

		projects, err := getProjectsByUserID(db, userID)
		if err != nil {
			log.Println("Ошибка получения проектов:", err)
			respondWithError(w, http.StatusInternalServerError, "Ошибка при получении проектов")
			return
		}

		respondWithJSON(w, http.StatusOK, projects)
	}
}

func getProjectsByUserID(db *sql.DB, userID int64) ([]ProjectResponse, error) {
	rows, err := db.Query(
		"SELECT id, name, description, type, created_at FROM projects WHERE owner_id = ? ORDER BY created_at DESC",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []ProjectResponse

	for rows.Next() {
		var project ProjectResponse
		err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.Description,
			&project.Type,
			&project.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		projects = append(projects, project)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if projects == nil {
		projects = []ProjectResponse{}
	}

	return projects, nil
}
