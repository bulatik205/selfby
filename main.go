package main

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/renderer/html"
)

var md = goldmark.New(
	goldmark.WithRendererOptions(
		html.WithUnsafe(),
	),
)

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/api/save", saveHandler)

	fmt.Println("Selfby running on :8080")
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
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
	if err := md.Convert(content, &buf); err != nil {
		http.Error(w, "Parse error", 500)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := `<!DOCTYPE html>
<html lang="ru">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>` + path + `</title>
    <style>
        :root {
            --bg: #0d1117;
            --text: #c9d1d9;
            --border: #30363d;
            --accent: #58a6ff;
            --code-bg: #161b22;
        }
        
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            background: var(--bg);
            color: var(--text);
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Helvetica, Arial, sans-serif;
            line-height: 1.6;
            max-width: 900px;
            margin: 0 auto;
            padding: 40px 20px;
        }
        
        h1, h2, h3, h4, h5, h6 {
            margin-top: 24px;
            margin-bottom: 16px;
            font-weight: 600;
            line-height: 1.25;
        }
        
        h1 {
            font-size: 2em;
            border-bottom: 1px solid var(--border);
            padding-bottom: 0.3em;
        }
        
        h2 {
            font-size: 1.5em;
            border-bottom: 1px solid var(--border);
            padding-bottom: 0.3em;
        }
        
        h3 {
            font-size: 1.25em;
        }
        
        p {
            margin-bottom: 16px;
        }
        
        a {
            color: var(--accent);
            text-decoration: none;
        }
        
        a:hover {
            text-decoration: underline;
        }
        
        code {
            background: var(--code-bg);
            padding: 0.2em 0.4em;
            border-radius: 6px;
            font-family: "SFMono-Regular", Consolas, "Liberation Mono", monospace;
            font-size: 85%;
        }
        
        pre {
            background: var(--code-bg);
            padding: 16px;
            border-radius: 6px;
            overflow-x: auto;
            margin-bottom: 16px;
        }
        
        pre code {
            background: none;
            padding: 0;
            font-size: 100%;
        }
        
        blockquote {
            border-left: 4px solid var(--border);
            padding-left: 16px;
            color: #8b949e;
            margin-bottom: 16px;
        }
        
        ul, ol {
            padding-left: 2em;
            margin-bottom: 16px;
        }
        
        li {
            margin-bottom: 4px;
        }
        
        img {
            max-width: 100%;
            height: auto;
        }
        
        hr {
            border: none;
            border-top: 1px solid var(--border);
            margin: 24px 0;
        }
        
        table {
            border-collapse: collapse;
            margin-bottom: 16px;
            width: 100%;
        }
        
        th, td {
            border: 1px solid var(--border);
            padding: 8px 12px;
            text-align: left;
        }
        
        th {
            background: var(--code-bg);
            font-weight: 600;
        }
        
        div[align="center"] {
            text-align: center;
        }
    </style>
</head>
<body>
` + buf.String() + `
</body>
</html>`

	w.Write([]byte(html))
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST", 405)
		return
	}

	path := r.FormValue("path")
	content := r.FormValue("content")

	cleanPath := filepath.Clean(path)
	if strings.HasPrefix(cleanPath, "..") {
		http.Error(w, "Nice try", 403)
		return
	}

	err := os.WriteFile(cleanPath, []byte(content), 0644)
	if err != nil {
		http.Error(w, "Save error", 500)
		return
	}

	w.Write([]byte("OK"))
}
