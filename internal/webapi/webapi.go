package webapi

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"groupie-tracker/internal/model"
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
type Client struct {
	httpClient *http.Client
}

// New ...
func New() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

// cancatURLs ...
func cancatURLs(URLs ...string) (url string) {
	if len(URLs) == 0 {
		slog.Error("Invalid URLs")
		return ""
	}

	url = strings.Join(URLs, "/")

	return url
}

// GetData ...
func (c *Client) GetDataFromAPI(url string) ([]byte, error) {
	res, err := c.httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return data, err
	}

	return data, nil
}

// getAllArtists ...
func (c *Client) GetAllArtists() ([]model.Artist, error) {
	artist := []model.Artist{}

	data, err := c.GetDataFromAPI(artistsURL)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(data, &artist); err != nil {
		return nil, err
	}

	return artist, nil
}

// GetArtistInfoByID ...
func (c *Client) GetArtistInfoByID(ID string) (model.ArtistInfo, error) {
	artistInfo := model.ArtistInfo{}

	artist, err := c.GetArtistByID(ID)
	if err != nil {
		slog.Error(err.Error())
		return artistInfo, err
	}

	location, err := c.GetLocationsByID(ID)
	if err != nil {
		slog.Error(err.Error())
		return artistInfo, err
	}

	dates, err := c.GetDatesByID(ID)
	if err != nil {
		slog.Error(err.Error())
		return artistInfo, err
	}

	relations, err := c.GetRelationsByID(ID)
	if err != nil {
		slog.Error(err.Error())
		return artistInfo, err
	}

	return model.ArtistInfo{
		Artist:    artist,
		Locations: location,
		Dates:     dates,
		Relations: relations,
	}, nil
}

// GetArtistByID ...
func (c *Client) GetArtistByID(ID string) (model.Artist, error) {
	artist := model.Artist{}

	data, err := c.GetDataFromAPI(cancatURLs(artistsURL, ID))
	if err != nil {
		slog.Error(err.Error())
		return artist, err
	}

	if err := json.Unmarshal(data, &artist); err != nil {
		slog.Error(err.Error())
		return artist, err
	}

	return artist, nil
}

// GetLocationsByID ...
func (c *Client) GetLocationsByID(ID string) (model.Locations, error) {
	locations := model.Locations{}

	data, err := c.GetDataFromAPI(cancatURLs(locationsURL, ID))
	if err != nil {
		slog.Error(err.Error())
		return locations, err
	}

	if err := json.Unmarshal(data, &locations); err != nil {
		slog.Error(err.Error())
		return locations, err
	}

	return locations, nil
}

// GetDatesByID ...
func (c *Client) GetDatesByID(ID string) (model.Dates, error) {
	dates := model.Dates{}

	data, err := c.GetDataFromAPI(cancatURLs(locationsURL, ID))
	if err != nil {
		slog.Error(err.Error())
		return dates, err
	}

	if err := json.Unmarshal(data, &dates); err != nil {
		slog.Error(err.Error())
		return dates, err
	}

	return dates, nil
}

// GetRelationsByID ...
func (c *Client) GetRelationsByID(ID string) (model.Relations, error) {
	relations := model.Relations{}

	data, err := c.GetDataFromAPI(cancatURLs(relationURL, ID))
	if err != nil {
		slog.Error(err.Error())
		return relations, err
	}

	if err := json.Unmarshal(data, &relations); err != nil {
		slog.Error(err.Error())
		return relations, err
	}

	return relations, nil
}
