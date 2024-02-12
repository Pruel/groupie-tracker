package webapi

import (
	"io"
	"net/http"
	"time"

	"groupie-tracker/pkg/config"
)

// API URLs
const (
	apiURL       = "https://groupietrackers.herokuapp.com/api/"
	artistsURL   = "https://groupietrackers.herokuapp.com/api/artists"
	locationsURL = "https://groupietrackers.herokuapp.com/api/locations"
	datesURL     = "https://groupietrackers.herokuapp.com/api/dates"
	relationURL  = "https://groupietrackers.herokuapp.com/api/relation"
)

// Client ...
type Clinet struct {
	httpClient *http.Client
}

// New ...
func New(cfg *config.Config) *Clinet {
	return &Clinet{
		httpClient: &http.Client{
			Timeout: cfg.HTTPClient.Timeout * time.Second,
		},
	}
}

// GetData ...
func (c *Clinet) GetData(url string) error {
	// check url
	if len(url) >= len(apiURL) {
		
	}
	// detected url

	res, err := c.httpClient.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)

	return nil
}
