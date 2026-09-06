package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

type UpdateWorkRequest struct {
	ProjectName string `json:"project_name"`
	Slug        string `json:"slug"`
	ContentMD   string `json:"content_md"`
}

func UpdateWork(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut && r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		userID, err := getUserIDFromSession(db, r)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Необходима авторизация")
			return
		}

		var req UpdateWorkRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Некорректные данные")
			return
		}

		var workID int64
		var projectID int64
		err = db.QueryRow(`
			SELECT w.id, w.at_project 
			FROM works w
			JOIN projects p ON w.at_project = p.id
			WHERE p.name = ? AND w.slug = ? AND w.owner_id = ?
		`, req.ProjectName, req.Slug, userID).Scan(&workID, &projectID)

		if err != nil {
			if err == sql.ErrNoRows {
				respondWithError(w, http.StatusNotFound, "Работа не найдена")
			} else {
				log.Println("Ошибка получения работы:", err)
				respondWithError(w, http.StatusInternalServerError, "Ошибка сервера")
			}
			return
		}

		md := goldmark.New(
			goldmark.WithExtensions(
				extension.GFM,
				extension.Footnote,
				extension.Typographer,
				extension.CJK,
			),
			goldmark.WithParserOptions(
				parser.WithAutoHeadingID(),
				parser.WithAttribute(),
			),
			goldmark.WithRendererOptions(
				html.WithHardWraps(),
				html.WithXHTML(),
			),
		)

		var buf strings.Builder
		if err := md.Convert([]byte(req.ContentMD), &buf); err != nil {
			log.Println("Ошибка конвертации Markdown:", err)
			respondWithError(w, http.StatusInternalServerError, "Ошибка обработки Markdown")
			return
		}
		contentHTML := buf.String()

		_, err = db.Exec(
			"UPDATE works SET content_md = ?, content_html = ? WHERE id = ?",
			req.ContentMD, contentHTML, workID,
		)
		if err != nil {
			log.Println("Ошибка обновления работы:", err)
			respondWithError(w, http.StatusInternalServerError, "Ошибка при обновлении")
			return
		}

		respondWithJSON(w, http.StatusOK, map[string]string{
			"message": "Работа обновлена",
		})
	}
}
