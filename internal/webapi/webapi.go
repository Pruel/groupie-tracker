package webapi

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"groupie-tracker/internal/entity"
)

// API URLs
const (
	apiURL       = "https://groupietrackers.herokuapp.com/api/"
	artistsURL   = "https://groupietrackers.herokuapp.com/api/artists"
	locationsURL = "https://groupietrackers.herokuapp.com/api/locations"
	datesURL     = "https://groupietrackers.herokuapp.com/api/dates"
	relationURL  = "https://groupietrackers.herokuapp.com/api/relation"
)

var LocationsURL = "https://groupietrackers.herokuapp.com/api/locations"

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
func (c *Client) GetAllArtists() ([]entity.Artist, error) {
	artist := []entity.Artist{}

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
func (c *Client) GetArtistInfoByID(ID string) (entity.ArtistInfo, error) {
	artistInfo := entity.ArtistInfo{}

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

	return entity.ArtistInfo{
		Artist:    artist,
		Locations: location,
		Dates:     dates,
		Relations: relations,
	}, nil
}

// GetArtistByID ...
func (c *Client) GetArtistByID(ID string) (entity.Artist, error) {
	artist := entity.Artist{}

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
func (c *Client) GetLocationsByID(ID string) (entity.Locations, error) {
	locations := entity.Locations{}

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

// GetLocationsByID ...
func (c *Client) GetLocationsByURL(URL string) (entity.Locations, error) {
	locations := entity.Locations{}

	data, err := c.GetDataFromAPI(URL)
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

// IV. locations group []string{}
// locURL := webapi.LocationsURL                // http://hearouku.app.locations
// locURL = strings.Join([]string{locURL}, "/") // // http://hearouku.app.locations/

// gID, ok := strings.CutPrefix(group.Locations, locURL) // http://hearouku.app.locations/1 = 1
// if !ok {
// 	continue
// }

// gLocs, err := cl.GetLocationsByID(gID)
// if err != nil {
// 	slog.Error(err.Error())
// 	continue
// }

// gLocs := webapi.GetLocationByURL(group.Locations)

// GetDatesByID ...
func (c *Client) GetDatesByID(ID string) (entity.Dates, error) {
	dates := entity.Dates{}

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
func (c *Client) GetRelationsByID(ID string) (entity.Relations, error) {
	relations := entity.Relations{}

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

// GetAllUniqueLocations

//  1. a new instance of the Location structure
//  2. get from webapi data by url = get all locations
//  3. after, unmarshall data into the locations instance
//  4. a. for range locations
//     b. location - elem =  create a new map with prealocating
//  5. if unique location, save this location into a new locations slice
//  6. parse locations value
//  7. return the new locations slice
func (c *Client) GetAllUniqueLocations() ([]string, error) {
	locations := entity.LocationsIndex{}

	data, err := c.GetDataFromAPI(locationsURL)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	if err := json.Unmarshal(data, &locations); err != nil { // Bug
		slog.Error(err.Error()) // Bug
		return nil, err
	}

	// map for unique locations
	alocSize := len(locations.Index)
	mapLoc := make(map[string]string, alocSize)
	uniqueLocs := make([]string, 0, alocSize)

	for _, locElem := range locations.Index {
		for _, loc := range locElem.Locations {
			// parseLocations
			loc = ParseAndFormatLocations(loc)
			if _, exists := mapLoc[loc]; !exists {
				mapLoc[loc] = loc
				uniqueLocs = append(uniqueLocs, loc)
			}
		}
	}

	return uniqueLocs, nil
}

// parseLocations
func ParseAndFormatLocations(loc string) string {
	loc = strings.ReplaceAll(loc, "-", ", ") // Берёт всю строку и первый аргумент наш целевой таргет что будет менять, второй аргумент на что меняем
	loc = strings.ReplaceAll(loc, "_", " ")  // some_cool_developer
	loc = strings.Title(loc)

	return loc
}
