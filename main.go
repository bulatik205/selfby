package main

import (
	"fmt"
	"net/http"
)

func main() {
	var bladeDir = "templates"

	http.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("styles"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, bladeDir+"/index.html")
	})

	fmt.Printf("Сервер запущен на http://localhost:8080\n")
	http.ListenAndServe(":8080", nil)
}
