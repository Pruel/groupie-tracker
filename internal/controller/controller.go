package controller

import (
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"groupie-tracker/internal/webapi"
)

const (
	indexTmpl  = "/internal/view/index.html"
	artistTmpl = "/internal/view/artist.html"
	errorTmpl  = "/internal/view/error.html"
	viewDir    = "/internal/view/"
)

// getTmplFilepath ...
func GetTmplFilepath(tmplName string) (tmplFilepath string) {
	wd, err := os.Getwd()
	if err != nil {
		slog.Error(err.Error())
	}

	switch tmplName {
	case "index.html", "index":
		tmplFilepath = wd + indexTmpl
	case "artist.html", "artist":
		tmplFilepath = wd + artistTmpl
	case "error.html", "error":
		tmplFilepath = wd + errorTmpl
	default:
		tmplFilepath = wd + viewDir
	}

	return tmplFilepath
}

func MainController(w http.ResponseWriter, r *http.Request) {
	// parsing a template
	tmp := template.Must(template.ParseFiles(GetTmplFilepath("index.html")))
	// error handling
	if r.URL.Path != "/" {
		ErrorController(w, r)

		return
	}

	// calling a some webapi method
	client := webapi.New()
	artists, err := client.GetAllArtists()
	if err != nil {
		slog.Error(err.Error())
	}

	// writing status and some header
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	// executing the templete
	if err := tmp.Execute(w, artists); err != nil {
		slog.Error(err.Error())
	}
}
