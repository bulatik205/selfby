package handlers

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"selfby/config"
)

var tmpl *template.Template

func init() {
	var err error
	tmpl, err = template.ParseFiles("templates/page.html")
	if err != nil {
		log.Fatal("Failed to parse template:", err)
	}
}

func MarkdownHandler(w http.ResponseWriter, r *http.Request) {
	path := filepath.Clean(r.URL.Path)

	if strings.HasSuffix(path, ".md") {
		renderMarkdown(w, r, path)
		return
	}

	http.ServeFile(w, r, "."+path)
}

func renderMarkdown(w http.ResponseWriter, r *http.Request, path string) {
	content, err := os.ReadFile("." + path)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	var buf bytes.Buffer
	if err := config.MD.Convert(content, &buf); err != nil {
		http.Error(w, "Parse error", 500)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := renderHTMLTemplate(path, template.HTML(buf.String()))
	w.Write([]byte(html))
}

func renderHTMLTemplate(title string, content template.HTML) string {
	var buf bytes.Buffer
	data := struct {
		Title   string
		Content template.HTML
	}{
		Title:   title,
		Content: content,
	}

	if err := tmpl.Execute(&buf, data); err != nil {
		log.Printf("Template execution error: %v", err)
		return string(content)
	}

	return buf.String()
}
