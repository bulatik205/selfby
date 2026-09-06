package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type NewWorkRequest struct {
	Title       string `json:"title"`
	ProjectName string `json:"project_name"`
}

type WorkResponse struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	ProjectID int64     `json:"project_id"`
	Views     int       `json:"views"`
	Likes     int       `json:"likes"`
	CreatedAt time.Time `json:"created_at"`
}

type SlugCheckRequest struct {
	Title       string `json:"title"`
	ProjectName string `json:"project_name"`
}

type SlugCheckResponse struct {
	Slug      string `json:"slug"`
	Available bool   `json:"available"`
	Message   string `json:"message,omitempty"`
}

func NewWork(db *sql.DB) http.HandlerFunc {
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

		var req NewWorkRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Некорректные данные")
			return
		}

		req.Title = strings.TrimSpace(req.Title)

		if errMsg := validateWorkTitle(req.Title); errMsg != "" {
			respondWithError(w, http.StatusBadRequest, errMsg)
			return
		}

		var projectID int64
		err = db.QueryRow(
			"SELECT id FROM projects WHERE owner_id = ? AND name = ?",
			userID, req.ProjectName,
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

		baseSlug := generateSlug(req.Title)
		slug, err := makeUniqueSlug(db, userID, baseSlug)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Ошибка сервера")
			return
		}

		result, err := db.Exec(
			"INSERT INTO works (owner_id, at_project, title, slug, content_md, content_html) VALUES (?, ?, ?, ?, '', '')",
			userID, projectID, req.Title, slug,
		)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Ошибка при создании работы")
			return
		}

		workID, err := result.LastInsertId()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Ошибка сервера")
			return
		}

		work := &WorkResponse{
			ID:        workID,
			Title:     req.Title,
			Slug:      slug,
			ProjectID: projectID,
			Views:     0,
			Likes:     0,
			CreatedAt: time.Now(),
		}

		log.Printf("Создана новая работа: %s (ID: %d) пользователем ID: %d в проекте ID: %d",
			work.Title, work.ID, userID, projectID)

		respondWithJSON(w, http.StatusCreated, work)
	}
}

func CheckSlug(db *sql.DB) http.HandlerFunc {
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

		var req SlugCheckRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			respondWithError(w, http.StatusBadRequest, "Некорректные данные")
			return
		}

		req.Title = strings.TrimSpace(req.Title)

		if req.Title == "" {
			respondWithJSON(w, http.StatusOK, SlugCheckResponse{
				Slug:      "",
				Available: false,
				Message:   "Название работы обязательно",
			})
			return
		}

		var projectID int64
		err = db.QueryRow(
			"SELECT id FROM projects WHERE owner_id = ? AND name = ?",
			userID, req.ProjectName,
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

		baseSlug := generateSlug(req.Title)

		if baseSlug == "" {
			respondWithJSON(w, http.StatusOK, SlugCheckResponse{
				Slug:      "",
				Available: false,
				Message:   "Не удалось сгенерировать slug",
			})
			return
		}

		exists, err := checkSlugExists(db, userID, baseSlug)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Ошибка сервера")
			return
		}

		if exists {
			respondWithJSON(w, http.StatusOK, SlugCheckResponse{
				Slug:      baseSlug,
				Available: false,
				Message:   "занято, будет добавлен суффикс",
			})
			return
		}

		respondWithJSON(w, http.StatusOK, SlugCheckResponse{
			Slug:      baseSlug,
			Available: true,
			Message:   "доступно",
		})
	}
}

func generateSlug(title string) string {
	translitMap := map[string]string{
		"а": "a", "б": "b", "в": "v", "г": "g", "д": "d",
		"е": "e", "ё": "e", "ж": "zh", "з": "z", "и": "i",
		"й": "y", "к": "k", "л": "l", "м": "m", "н": "n",
		"о": "o", "п": "p", "р": "r", "с": "s", "т": "t",
		"у": "u", "ф": "f", "х": "h", "ц": "ts", "ч": "ch",
		"ш": "sh", "щ": "sch", "ъ": "", "ы": "y", "ь": "",
		"э": "e", "ю": "yu", "я": "ya",
		"А": "a", "Б": "b", "В": "v", "Г": "g", "Д": "d",
		"Е": "e", "Ё": "e", "Ж": "zh", "З": "z", "И": "i",
		"Й": "y", "К": "k", "Л": "l", "М": "m", "Н": "n",
		"О": "o", "П": "p", "Р": "r", "С": "s", "Т": "t",
		"У": "u", "Ф": "f", "Х": "h", "Ц": "ts", "Ч": "ch",
		"Ш": "sh", "Щ": "sch", "Ъ": "", "Ы": "y", "Ь": "",
		"Э": "e", "Ю": "yu", "Я": "ya",
	}

	var result strings.Builder
	for _, char := range title {
		if transliterated, ok := translitMap[string(char)]; ok {
			result.WriteString(transliterated)
		} else {
			result.WriteRune(char)
		}
	}

	slug := strings.ToLower(result.String())

	reg := regexp.MustCompile(`[^a-z0-9\s-]`)
	slug = reg.ReplaceAllString(slug, " ")

	reg = regexp.MustCompile(`\s+`)
	slug = reg.ReplaceAllString(slug, "-")

	reg = regexp.MustCompile(`-+`)
	slug = reg.ReplaceAllString(slug, "-")

	slug = strings.Trim(slug, "-")

	if slug == "" {
		return "work"
	}

	return slug
}

func makeUniqueSlug(db *sql.DB, ownerID int64, baseSlug string) (string, error) {
	slug := baseSlug
	counter := 1

	for {
		exists, err := checkSlugExists(db, ownerID, slug)
		if err != nil {
			return "", err
		}

		if !exists {
			return slug, nil
		}

		slug = baseSlug + "-" + fmt.Sprintf("%d", counter)
		counter++
	}
}

func checkSlugExists(db *sql.DB, ownerID int64, slug string) (bool, error) {
	var count int
	err := db.QueryRow(
		"SELECT COUNT(*) FROM works WHERE owner_id = ? AND slug = ?",
		ownerID, slug,
	).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func validateWorkTitle(title string) string {
	if title == "" {
		return "Название работы обязательно"
	}

	if len(title) < 3 {
		return "Название должно быть не менее 3 символов"
	}

	if len(title) > 100 {
		return "Название должно быть не более 100 символов"
	}

	return ""
}
