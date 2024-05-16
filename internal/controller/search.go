package controller

import (
	"fmt"
	"groupie-tracker/internal/entity"
	"groupie-tracker/internal/webapi"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	NameType         = " - artist/band"
	LocationType     = " - location"
	MembersType      = " - member"
	CreationDateType = " - creation date"
	PubType          = " - first album date"
)

func GetSearchValue(r *http.Request) (searchValue string, searchType string) {
	searchValue = r.FormValue("search")

	if len(searchValue) != 0 {
		if !strings.Contains(searchValue, " - ") {
			slog.Error("Invalid search query")
		}
	}

	slice := strings.Split(searchValue, " - ")

	return slice[0], slice[1]
}

func Search(searchValue string, searchType string, artists []entity.Artist) (foundGroups []entity.Artist, err error) {
	searchType = " - " + searchType

	switch searchType {
	case NameType:
		foundGroups = searchByName(searchValue, artists)
	case MembersType:
		foundGroups = searchByMemebers(searchValue, artists)
	case LocationType:
		foundGroups = searchByLocations(searchValue, artists)
	case CreationDateType:
		foundGroups = searchByCreateGroup(searchValue, artists)
	case PubType:
		foundGroups = searchByPublicAlbum(searchValue, artists)
	}

	return foundGroups, nil
}

func searchByName(searchValue string, artists []entity.Artist) (foundGroups []entity.Artist) {
	if artists == nil {
		slog.Error("Error, empty data")
		return nil
	}

	fmt.Println("calling searchByName")

	for _, elem := range artists {
		if searchValue == elem.Name {
			foundGroups = append(foundGroups, elem)
		}
	}

	return foundGroups
}

func searchByMemebers(searchValue string, artists []entity.Artist) (foundGroups []entity.Artist) {
	if artists == nil {
		slog.Error("Error, empty data")
		return nil
	}

	for _, group := range artists {
		for _, mem := range group.Members {
			if searchValue == mem {
				foundGroups = append(foundGroups, group)
			}
		}
	}

	return foundGroups
}

func searchByLocations(searchValue string, artist []entity.Artist) (foundGroups []entity.Artist) {
	if artist == nil {
		slog.Error("Error, empty data")
		return nil
	}

	for _, elem := range artist {
		location, _ := webapi.New().GetLocationsByURL(elem.Locations)
		for _, loc := range location.Locations {
			loc = webapi.ParseAndFormatLocations(loc)

			if loc == searchValue {
				foundGroups = append(foundGroups, elem)
			}
		}
	}

	return foundGroups
}

func searchByCreateGroup(searchValue string, artist []entity.Artist) (foundGroups []entity.Artist) {
	if artist == nil {
		slog.Error("Error, empty data")
		return nil
	}

	for _, elem := range artist {
		if strconv.Itoa(elem.CreationDate) == searchValue {
			foundGroups = append(foundGroups, elem)
		}
	}

	return foundGroups
}

func searchByPublicAlbum(searchValue string, artist []entity.Artist) (foundGroups []entity.Artist) {
	if artist == nil {
		slog.Error("Error, empty data")
		return nil
	}
	timeFormat := "02-01-2006"

	for _, elem := range artist {
		firstAlbum, err := time.Parse(timeFormat, elem.FirstAlbum)
		if err != nil {
			slog.Error(err.Error())
			return nil
		}

		if strconv.Itoa(firstAlbum.Year()) == searchValue {
			foundGroups = append(foundGroups, elem)
		}
	}

	return foundGroups
}

func GetAllUniqueSuggestions(dataArtist []entity.Artist) map[string]string {
	alocSize := len(dataArtist)
	udata := make(map[string]string, alocSize)

	for _, artist := range dataArtist {
		nameKey := artist.Name + NameType

		if _, ok := udata[artist.Name]; !ok {
			udata[artist.Name] = nameKey
		}

		for _, mem := range artist.Members {
			if _, ok := udata[mem]; !ok {
				udata[mem] = mem + MembersType
			}
		}

		uniqLoc, err := webapi.New().GetAllUniqueLocations()
		if err != nil {
			slog.Error(err.Error())
		}
		for _, loc := range uniqLoc {
			udata[loc] = loc + LocationType
		}

		year := strconv.Itoa(artist.CreationDate)
		if _, ok := udata[year]; !ok {
			udata[year] = year + CreationDateType
		}

		timeFormat := "02-01-2006"

		falb, err := time.Parse(timeFormat, artist.FirstAlbum)
		if err != nil {
			slog.Error(err.Error())
		}
		pubAlb := strconv.Itoa(falb.Year())
		if _, ok := udata[pubAlb]; !ok {
			udata[pubAlb] = pubAlb + PubType
		}

	}
	return udata
}
