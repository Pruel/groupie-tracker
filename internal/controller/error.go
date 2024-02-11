package controller

import (
	"net/http"
)

// ErrorController
func ErrorController(w http.ResponseWriter, r *http.Request) {
	// template parsing
	// tmpl := template.Must(template.ParseFiles("someDir/files"))

	// webapi-> some data

	// error handling
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusNotFound)

	/// exeting the error template with some args
	// tmpl.Execute(w, data)
}
