package api

import (
	"database/sql"
	"log"
	"net/http"
)

func DeleteWork(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		userID, err := getUserIDFromSession(db, r)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Необходима авторизация")
			return
		}

		projectName := r.URL.Query().Get("project")
		workSlug := r.URL.Query().Get("slug")

		if projectName == "" || workSlug == "" {
			respondWithError(w, http.StatusBadRequest, "Не указаны параметры")
			return
		}

		var workID int64
		var workTitle string
		err = db.QueryRow(`
			SELECT w.id, w.title
			FROM works w
			JOIN projects p ON w.at_project = p.id
			WHERE p.owner_id = ? AND p.name = ? AND w.slug = ?
		`, userID, projectName, workSlug).Scan(&workID, &workTitle)

		if err != nil {
			if err == sql.ErrNoRows {
				respondWithError(w, http.StatusNotFound, "Работа не найдена")
			} else {
				log.Println("Ошибка получения работы:", err)
				respondWithError(w, http.StatusInternalServerError, "Ошибка сервера")
			}
			return
		}

		tx, err := db.Begin()
		if err != nil {
			log.Println("Ошибка начала транзакции:", err)
			respondWithError(w, http.StatusInternalServerError, "Ошибка сервера")
			return
		}
		defer tx.Rollback()

		_, err = tx.Exec(
			"DELETE FROM likes WHERE element_id = ? AND element_type = 'work'",
			workID,
		)
		if err != nil {
			log.Println("Ошибка удаления лайков:", err)
			respondWithError(w, http.StatusInternalServerError, "Ошибка сервера")
			return
		}

		_, err = tx.Exec(
			"DELETE FROM views WHERE element_id = ? AND element_type = 'work'",
			workID,
		)
		if err != nil {
			log.Println("Ошибка удаления просмотров:", err)
			respondWithError(w, http.StatusInternalServerError, "Ошибка сервера")
			return
		}

		_, err = tx.Exec("DELETE FROM works WHERE id = ?", workID)
		if err != nil {
			log.Println("Ошибка удаления работы:", err)
			respondWithError(w, http.StatusInternalServerError, "Ошибка сервера")
			return
		}

		if err = tx.Commit(); err != nil {
			log.Println("Ошибка коммита:", err)
			respondWithError(w, http.StatusInternalServerError, "Ошибка сервера")
			return
		}

		log.Printf("Удалена работа: %s (ID: %d) пользователем ID: %d", workTitle, workID, userID)

		respondWithJSON(w, http.StatusOK, map[string]string{
			"message": "Работа удалена",
		})
	}
}

func DeleteProject(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		userID, err := getUserIDFromSession(db, r)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Необходима авторизация")
			return
		}

		projectName := r.URL.Query().Get("name")
		if projectName == "" {
			respondWithError(w, http.StatusBadRequest, "Не указано имя проекта")
			return
		}

		var projectID int64
		err = db.QueryRow(
			"SELECT id FROM projects WHERE owner_id = ? AND name = ?",
			userID, projectName,
		).Scan(&projectID)

		if err != nil {
			if err == sql.ErrNoRows {
				respondWithError(w, http.StatusNotFound, "Проект не найден")
			} else {
				log.Println("Ошибка получения проекта:", err)
				respondWithError(w, http.StatusInternalServerError, "Ошибка сервера")
			}
			return
		}

		tx, err := db.Begin()
		if err != nil {
			log.Println("Ошибка начала транзакции:", err)
			respondWithError(w, http.StatusInternalServerError, "Ошибка сервера")
			return
		}
		defer tx.Rollback()

		rows, err := tx.Query("SELECT id FROM works WHERE at_project = ?", projectID)
		if err != nil {
			log.Println("Ошибка получения работ:", err)
			respondWithError(w, http.StatusInternalServerError, "Ошибка сервера")
			return
		}

		var workIDs []int64
		for rows.Next() {
			var id int64
			if err := rows.Scan(&id); err != nil {
				continue
			}
			workIDs = append(workIDs, id)
		}
		rows.Close()

		for _, workID := range workIDs {
			_, err = tx.Exec(
				"DELETE FROM likes WHERE element_id = ? AND element_type = 'work'",
				workID,
			)
			if err != nil {
				log.Println("Ошибка удаления лайков:", err)
			}

			_, err = tx.Exec(
				"DELETE FROM views WHERE element_id = ? AND element_type = 'work'",
				workID,
			)
			if err != nil {
				log.Println("Ошибка удаления просмотров:", err)
			}
		}

		_, err = tx.Exec(
			"DELETE FROM likes WHERE element_id = ? AND element_type = 'project'",
			projectID,
		)
		if err != nil {
			log.Println("Ошибка удаления лайков проекта:", err)
		}

		_, err = tx.Exec(
			"DELETE FROM views WHERE element_id = ? AND element_type = 'project'",
			projectID,
		)
		if err != nil {
			log.Println("Ошибка удаления просмотров проекта:", err)
		}

		_, err = tx.Exec("DELETE FROM works WHERE at_project = ?", projectID)
		if err != nil {
			log.Println("Ошибка удаления работ:", err)
			respondWithError(w, http.StatusInternalServerError, "Ошибка сервера")
			return
		}

		_, err = tx.Exec("DELETE FROM projects WHERE id = ?", projectID)
		if err != nil {
			log.Println("Ошибка удаления проекта:", err)
			respondWithError(w, http.StatusInternalServerError, "Ошибка сервера")
			return
		}

		if err = tx.Commit(); err != nil {
			log.Println("Ошибка коммита:", err)
			respondWithError(w, http.StatusInternalServerError, "Ошибка сервера")
			return
		}

		log.Printf("Удален проект: %s (ID: %d) пользователем ID: %d", projectName, projectID, userID)

		respondWithJSON(w, http.StatusOK, map[string]string{
			"message": "Проект удален",
		})
	}
}
