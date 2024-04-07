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

	// 1. prepare for filtering data
	// Делаем структуру для наших данных.
	// mainData -> artists, filtersData
	// filter -> creationDate, firstAlbum, Members, Locations
	// Artists

	// 2.
	// FRONT all data in one form
	// <form>, begin

	// creationDate
	// inpute -> type range, value=1950, min=1900, max=2020, name=creationDate

	// firstAlbum
	// input type date, name=firstAlbum, value=1900

	// members
	// input type check box, value=1, name=members

	// locations
	// select name=locations
	// {{ range .locations }}
	// option value=some_location
	// {{ end }}
	// </form> end

	// 3. receive filters params from front and save this params in filters structure (Пишем структуру опять куда сохранять)

	// 4. check and validate data, and parse or converting filters data

	// 5. do sort by filters params

	// 6. create artists slice for saving after filtering

	// 7. after do filters -> for -> element equel by filter param -> true -> append filteredArtists type = model.Artists

	// 8. return filteredArtists

	// calling a some webapi method
	client := webapi.New()
	artists, err := client.GetAllArtists()
	if err != nil {
		slog.Error(err.Error())
	}

	// writing status and some header
	// w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	// executing the templete
	if err := tmp.Execute(w, artists); err != nil {
		slog.Error(err.Error())
	}
}
