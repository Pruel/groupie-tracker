package controller

import (
	"encoding/json"
	"errors"
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

	glocs, err := prepareDataForGeolcs(artistInfo.Locations.Locations)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	artistInfo.Geolocations = glocs

	// writing status and some header
	// w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)

	// executing the templete
	if err := tmp.Execute(w, artistInfo); err != nil {
		slog.Error(err.Error())
	}
}

// prepareDataForGeolocs
func prepareDataForGeolcs(sslc []string) (fdata string, err error) {
	if len(sslc) == 0 {
		slog.Error("error, invalid data")
		return "", errors.New("error, invalid data")
	}

	bufLocs := make([]string, 0, len(sslc))

	for _, loc := range sslc {
		loc := webapi.ParseAndFormatLocations(loc)
		bufLocs = append(bufLocs, loc)
	}

	bdata, err := json.Marshal(&bufLocs)
	if err != nil {
		slog.Error(err.Error())
		return "", err
	}

	return string(bdata), nil
}
