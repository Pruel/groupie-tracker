package controller

import (
	"html/template"
	"log/slog"
	"net/http"

	"groupie-tracker/internal/webapi"
)

// ArtistController
func ArtistController(w http.ResponseWriter, r *http.Request) {
	// parsing a template
	tmp := template.Must(template.ParseFiles(GetTmplFilepath("artist.html")))
	// error handling
	if r.URL.Path != "/artist" {
		ErrorController(w, r)
		return
	}

	artistID := r.URL.Query().Get("artistID")
	// artID, err  := strconv.Atoi(artistID)
	// if err != nil {
	// 	slog.Error(err.Error())
	// }

	// calling a some webapi method
	client := webapi.New()
	artistInfo, err := client.GetArtistInfoByID(artistID)
	if err != nil {
		slog.Error(err.Error())
	}

	// writing status and some header
	// w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	// executing the templete
	if err := tmp.Execute(w, artistInfo); err != nil {
		slog.Error(err.Error())
	}
}
