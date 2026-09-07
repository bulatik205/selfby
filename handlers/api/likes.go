package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type LikeRequest struct {
	ElementID   int64  `json:"element_id"`
	ElementType string `json:"element_type"`
}

type LikeResponse struct {
	Liked bool `json:"liked"`
	Count int  `json:"count"`
}

func ToggleLike(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		userID, err := getUserIDFromSession(db, r)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Необходима авторизация")
			return
		}

		var req LikeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Некорректные данные")
			return
		}

		if req.ElementID == 0 {
			respondWithError(w, http.StatusBadRequest, "Не указан element_id")
			return
		}

		if req.ElementType == "" {
			req.ElementType = "work"
		}

		var likeID int64
		err = db.QueryRow(
			"SELECT id FROM likes WHERE element_id = ? AND element_type = ? AND owner_id = ?",
			req.ElementID, req.ElementType, userID,
		).Scan(&likeID)

		var liked bool

		if err == sql.ErrNoRows {
			_, err = db.Exec(
				"INSERT INTO likes (element_id, element_type, owner_id) VALUES (?, ?, ?)",
				req.ElementID, req.ElementType, userID,
			)
			if err != nil {
				log.Println("Ошибка добавления лайка:", err)
				respondWithError(w, http.StatusInternalServerError, "Ошибка сервера")
				return
			}
			liked = true
		} else if err != nil {
			log.Println("Ошибка проверки лайка:", err)
			respondWithError(w, http.StatusInternalServerError, "Ошибка сервера")
			return
		} else {
			_, err = db.Exec(
				"DELETE FROM likes WHERE id = ?",
				likeID,
			)
			if err != nil {
				log.Println("Ошибка удаления лайка:", err)
				respondWithError(w, http.StatusInternalServerError, "Ошибка сервера")
				return
			}
			liked = false
		}

		var count int
		err = db.QueryRow(
			"SELECT COUNT(*) FROM likes WHERE element_id = ? AND element_type = ?",
			req.ElementID, req.ElementType,
		).Scan(&count)
		if err != nil {
			log.Println("Ошибка подсчета лайков:", err)
			count = 0
		}

		respondWithJSON(w, http.StatusOK, LikeResponse{
			Liked: liked,
			Count: count,
		})
	}
}

func CheckLike(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		elementID := r.URL.Query().Get("element_id")
		elementType := r.URL.Query().Get("element_type")

		if elementType == "" {
			elementType = "work"
		}

		if elementID == "" {
			respondWithError(w, http.StatusBadRequest, "Не указан element_id")
			return
		}

		var count int
		err := db.QueryRow(
			"SELECT COUNT(*) FROM likes WHERE element_id = ? AND element_type = ?",
			elementID, elementType,
		).Scan(&count)
		if err != nil {
			log.Println("Ошибка подсчета лайков:", err)
			count = 0
		}

		liked := false
		userID, err := getUserIDFromSession(db, r)
		if err == nil {
			var exists bool
			err = db.QueryRow(
				"SELECT EXISTS(SELECT 1 FROM likes WHERE element_id = ? AND element_type = ? AND owner_id = ?)",
				elementID, elementType, userID,
			).Scan(&exists)
			if err == nil {
				liked = exists
			}
		}

		respondWithJSON(w, http.StatusOK, LikeResponse{
			Liked: liked,
			Count: count,
		})
	}
}
