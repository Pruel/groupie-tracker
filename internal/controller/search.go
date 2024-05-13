package controller

import (
	"fmt"
	"groupie-tracker/internal/entity"
	"groupie-tracker/internal/filter"
	"groupie-tracker/internal/webapi"
	"html/template"
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

func SearchController(w http.ResponseWriter, r *http.Request) {
	tmp := template.Must(template.ParseFiles(GetTmplFilepath("main.html")))

	artists, err := webapi.New().GetAllArtists()
	if err != nil {
		slog.Error(err.Error())
	}

	filtersData, err := filter.PrepareFilterData(artists)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	// data <- getAllUniqueSuggestions
	udata := getAllUniqueSuggestions(artists)
	// fmt.Printf("Unique Sugesstion: %v \n\n", udata) // ++
	//

	// Получение значений поиска
	sValue, sType := getSearchValue(r) // Janika, // sValue, sType

	// Ищем по запросу
	foundGroups, err := Search(sValue, sType, artists)
	// error cheking
	fmt.Printf("\n\nResult: %v\n\n", foundGroups)
	mdata := entity.MainData{
		Artists:     foundGroups,
		FiltersData: *filtersData,
		SearchData:  udata,
	}

	// Выполняем наш шаблон
	err = tmp.Execute(w, mdata)
	if err != nil {
		slog.Error(err.Error())
	}
	// getSearchValue - call second function

	// getMainData = recieve all groups from web api, and return []entity.Artists -> webapi.New().GetAllArtists()

	// Search -> searchByGroupName

	// last call template with found groups as context data of the template

}

// II getSearchValue
func getSearchValue(r *http.Request) (searchValue string, searchType string) {

	// read request body = r.FormValue("search") - метод используется для получения значения поля ввода
	// <input type="text" name="search" - отправленное через HTML форму
	searchValue = r.FormValue("search") // "Value - Type " = "" || "Eminem - artist/band"

	// 1. validate search query, " - ", "", strings.Contains(str, " - ")
	if len(searchValue) != 0 { // "Value" -> big <- front
		if !strings.Contains(searchValue, " - ") {
			slog.Error("Invalid search query")
		}
	}

	// 2. split search query by " - ", ->
	// strSlice := strings.Split(searchValue, " - ") // strSlice := strings.Split(str, substr) -> []string{"value", "type") || strSlice[0] = value, strSlice[1] = searchType
	slice := strings.Split(searchValue, " - ")

	return slice[0], slice[1]
}

// III Search
func Search(searchValue string, searchType string, artists []entity.Artist) (foundGroups []entity.Artist, err error) {
	fmt.Printf("Equal: %v , sType: %s = NameType: %s\n", searchType == NameType, searchType, NameType)

	searchType = " - " + searchType

	fmt.Printf("Equal: %v , sType: %s = NameType: %s\n", searchType == NameType, searchType, NameType)

	switch searchType {
	case NameType:
		foundGroups = searchByName(searchValue, artists) // +++
	case MembersType:
		foundGroups = searchByMemebers(searchValue, artists) // +++
	case LocationType:
		foundGroups = searchByLocations(searchValue, artists) // +++
	case CreationDateType:
		foundGroups = searchByCreateGroup(searchValue, artists) // +++
	case PubType:
		foundGroups = searchByPublicAlbum(searchValue, artists) // +++
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
	timeFormat := "02-01-2006" // Important Golang date ;-) || time.Year(), time.Month(), time.Day, time.Minute, time.Second, time.Parse. Если у нас есть условная дата, то лучше проводить манипуляции
	// вместе с пакетом time(), ну и так семантически вернее будет.

	// str -> year

	// 1. str -> time -> time.Year // 1. time.Parse(timeFormat, str)
	// 2. str -> strings.CutPrefix = "02-01-2006"

	for _, elem := range artist {
		// first a
		// Функция time parse возвращает error. Мы можем опустить ошибку, но это анти-паттерн и по хорошему если функция возвращает ошибку, то мы её тоже пишем!!!
		firstAlbum, err := time.Parse(timeFormat, elem.FirstAlbum) // firstAlbum = type time
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

func getAllUniqueSuggestions(dataArtist []entity.Artist) map[string]string {

	// all groups or artists <- []entity.Artists
	alocSize := len(dataArtist)
	udata := make(map[string]string, alocSize)

	for _, artist := range dataArtist {
		// SearchData

		// // 1. unique name of the group
		nameKey := artist.Name + NameType // New - artist/band

		if _, ok := udata[artist.Name]; !ok {
			udata[artist.Name] = nameKey
		}

		// 2. unique member
		for _, mem := range artist.Members { // []string{Member}
			if _, ok := udata[mem]; !ok {
				udata[mem] = mem + MembersType
			}
		}

		// 3. unique location
		uniqLoc, err := webapi.New().GetAllUniqueLocations()
		if err != nil {
			slog.Error(err.Error())
		}
		for _, loc := range uniqLoc {
			udata[loc] = loc + LocationType
		}

		// 4. unique creation date // Daniil
		year := strconv.Itoa(artist.CreationDate)
		if _, ok := udata[year]; !ok {
			udata[year] = year + CreationDateType
		}

		// 5. unique first album publishing // Yanika Time.Year()
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

// full text seatch by types and params
// 0. front-end -> input with datalist

// 1. getUniqueSuggestions

// 2. recieving data search request from front end and validate this data ||search-> New-York - locations

// 3. search by any field (type) of the groups
