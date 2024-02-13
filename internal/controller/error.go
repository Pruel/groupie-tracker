package controller

import (
	"html/template"
	"log/slog"
	"net/http"
)

// ErrorController
func ErrorController(w http.ResponseWriter, r *http.Request) {
	// template parsing
	tmpl := template.Must(template.ParseFiles(GetTmplFilepath("error.html")))

	// error handling
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)

	// exeting the error template with some args
	data := "404 NOT FOUND"
	if err := tmpl.Execute(w, data); err != nil {
		slog.Error(err.Error())
	}
}
