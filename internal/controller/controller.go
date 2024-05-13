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

func MainController(w http.ResponseWriter, r *http.Request) { // index.html, -> / <- index.html
	// parsing a template
	tmp := template.Must(template.ParseFiles(GetTmplFilepath("main.html")))
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

	fgroups := make([]entity.Artist, len(artists))
	copy(fgroups, artists)

	filtersData, err := filter.PrepareFilterData(artists)
	if err != nil {
		slog.Debug("Error prepare filter data")
	}

	udata := getAllUniqueSuggestions(artists) // this code line

	// prepare filter date
	mdata := entity.MainData{ // container structure
		Artists:     fgroups,
		FiltersData: *filtersData, // *Pointer, zero value > nil
		SearchData:  udata,

		// &value = *Pointer = addres of the RAM // reference || *Pointer = Ox1f3f5d56fd

		// *value // unreference
	}

	// fmt.Printf("Filters data:\n\n %+v\n\n", filtersData) //%v = value , %+v key:value
	// fmt.Printf("Artists: %+v\n\n", artists)              //%v = value , %+v key:value

	// writing status and some header
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	// executing the templete
	if err := tmp.Execute(w, mdata); err != nil {
		slog.Error(err.Error())
	}
}

// 1. prepare for filtering data +++
// Делаем структуру для наших данных.
// mainData -> artists, filtersData
// filter -> creationDate, firstAlbum, Members, Locations
// Artists

// 2. +++
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

// 3. receive filters params from front and save this params in filters structure (Пишем структуру опять куда сохранять) +++

// 4. check and validate data, and parse or converting filters data

// 5. do sort by filters params

// 6. create artists slice for saving after filtering

// 7. after do filters -> for -> element equel by filter param -> true -> append filteredArtists type = model.Artists

// 8. return filteredArtists
