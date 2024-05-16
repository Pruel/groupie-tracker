package controller

import (
	"html/template"
	"log/slog"
	"net/http"
	"os"

	"groupie-tracker/internal/entity"
	"groupie-tracker/internal/filter"
	"groupie-tracker/internal/webapi"
)

const (
	indexTmpl  = "/internal/view/index.html"
	mainTmpl   = "/internal/view/main.html"
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
	case "main.html", "main":
		tmplFilepath = wd + mainTmpl
	case "artist.html", "artist":
		tmplFilepath = wd + artistTmpl
	case "error.html", "error":
		tmplFilepath = wd + errorTmpl
	case "index.html", "index":
		tmplFilepath = wd + indexTmpl
	default:
		tmplFilepath = wd + viewDir
	}

	return tmplFilepath
}

func MainController(w http.ResponseWriter, r *http.Request) {
	tmp := template.Must(template.ParseFiles(GetTmplFilepath("main.html")))

	if r.URL.Path != "/" && r.URL.Path != "/filter" && r.URL.Path != "/search" {
		ErrorController(w, r)
		return
	}

	artists, err := webapi.New().GetAllArtists()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	fltData, err := filter.PrepareFilterData(artists)
	if err != nil {
		slog.Debug("Error prepare filter data")
		return
	}

	udata := GetAllUniqueSuggestions(artists)

	foundGroups := make([]entity.Artist, len(artists))
	copy(foundGroups, artists)
	var message string

	mdata := entity.MainData{
		Artists:     foundGroups,
		FiltersData: *fltData,
		SearchData:  udata,
		Message:     message,
	}

	if r.Method == http.MethodPost && r.URL.Path == "/filter" {
		slog.Info("Handling POST METHOD for filter")
		ReadValidateAndSaveFilterData(r, fltData)
		filteredArtists, filterMessage := Filter(fltData, artists)
		if filteredArtists == nil {
			slog.Debug("Error, empty filter")
			return
		}
		mdata.Artists = filteredArtists
		mdata.Message = filterMessage
	} else if r.Method == http.MethodPost && r.URL.Path == "/search" {
		slog.Info("Handling POST METHOD for search")
		sValue, sType := GetSearchValue(r)
		foundGroups, err = Search(sValue, sType, artists)
		if err != nil {
			slog.Error(err.Error())
			slog.Debug("Error, empty search")
			return
		}
		mdata.Artists = foundGroups
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	if err := tmp.Execute(w, mdata); err != nil {
		slog.Error(err.Error())
	}
}
