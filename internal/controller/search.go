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

	// Получение значений поиска
	requestUser, err := getSearchValue(r) // Janika
	if err != nil {
		slog.Error(err.Error())
		return
	}

	searchValue := r.FormValue("search")

	fmt.Printf("\nSearchValue: %v \n", searchValue)

	// Ищем по запросу
	foundGroups := Search(requestUser, "text", artists) // Daniil

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

// II
// getSearchValue
func getSearchValue(r *http.Request) (searchValue string, searchType string) {

	// read request body = r.FormValue("search") - метод используется для получения значения поля ввода
	// <input type="text" name="search" - отправленное через HTML форму
	searchValue = r.FormValue("search") // "Value - Type " = "" || "Eminem - artist/band"
	
	// 1. validate search query, " - ", "", strings.Contains(str, " - ")
	if len(searchValue) != 0 {
		if !strings.Contains(searchValue, " - ") { 
			slog.Error("Invalid search query")
		}
	}
	
	// 2. split search query by " - ", ->
	// strSlice := strings.Split(searchValue, " - ") // strSlice := strings.Split(str, substr) -> []string{"value", "type") || strSlice[0] = value, strSlice[1] = searchType
	slice := strings.Split(searchValue, " - ")
	

	return slice[0], slice[1]
}

func Search(searchValue string, searchType string) ([]entity.Artist, error) {
	// Функция поиска для обработки сценариев
	
}

























// // III
// // Search(searchValue string, artists []entity.Artists), search by group name, and return found groups = foundGroups []entity.Artists
// func Search(searchValue string, searchType string, artist []entity.Artist) (foundGroups []entity.Artist) {
// 	if searchValue == "" || searchType == "" {
// 		slog.Error("")
// 		return nil
// 	}

// 	switch searchType { // searchType = " - members" ||
// 	case NameType:
// 		// searchByName, Daniil
// 	case MembersType:
// 		// searchByMemebers() here be func for our logic programm, Janika
// 	case LocationType:
// 		// searchByLocations, Daniil
// 	case PubType:
// 		// searchByReleaseAlbum, Janika
// 	default:

// 	}

// 	// NameType         = " - artist/band"
// 	// LocationType     = " - location"
// 	// MembersType      = " - member"
// 	// CreationDateType = " - creation date"
// 	// PubType

// 	return foundGroups
// }

func searchByName(searchValue string) []entity.Artist {
	// var foundGroups []entity.Artist

	// if artist != nil {
	// 	for _, group := range artist {
	// 		if strings.Contains(group.Name, searchValue) {
	// 			foundGroups = append(foundGroups, group)
	// 		}
	// 	}
	// }

	// return foundGroups
	// check artists on not nil (artists type slice, zero value of the slice is nil)
	// if ok
	// for _, group := range artists {
	//     if strings.Contains(group.Name, searchValue) -> ok = true {
	// foundGroups = append(foundGroups, group)
	//        }
	//      for _, loc := group.Locations {
	// if loc == searchValue {}
	//         }
	// }
	//  after return foundGroups

	return nil
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
