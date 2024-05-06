package controller

import (
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
	creationDateType = " - creation date"
	PubType          = " - first album date"
)

// Davai imposter pishi constanты
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
	requestUser, err := getSearchValue(r)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	// Ищем по запросу
	foundGroups := Search(requestUser, artists)

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
func getSearchValue(r *http.Request) (searchValue string, err error) {

	// read request body = r.FormValue("search") - метод используется для получения значения поля ввода
	// <input type="text" name="search" - отправленное через HTML форму
	searchValue = r.FormValue("search") // "Value - Type " = ""

	// check search value, if str is empty -> return error
	if searchValue == "" {
		slog.Error(err.Error())
	}
	// if ok, return search value

	return searchValue, nil
}

// III
// Search(searchValue string, artists []entity.Artists), search by group name, and return found groups = foundGroups []entity.Artists
func Search(searchValue string, artist []entity.Artist) []entity.Artist {
	var foundGroups []entity.Artist

	if artist != nil {
		for _, group := range artist {
			if strings.Contains(group.Name, searchValue) {
				foundGroups = append(foundGroups, group)
			}
		}
	}

	return foundGroups
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
			udata[year] = year + creationDateType
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
