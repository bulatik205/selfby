package main

import (
	"fmt"
	"net/http"

	"selfby/handlers"
)

func main() {
	http.HandleFunc("/", handlers.MarkdownHandler)

	fmt.Println("Selfby running on :8080")
	http.ListenAndServe(":8080", nil)
}
