package controller

import (
	"groupie-tracker/internal/entity"
	"groupie-tracker/internal/webapi"
	"html/template"
	"log/slog"
	"net/http"
	"strings"
)

func SearchController(w http.ResponseWriter, r *http.Request) {
	// Парсим шаблончес
	tmp := template.Must(template.ParseFiles(GetTmplFilepath("index.html")))

	artists, err := webapi.New().GetAllArtists()
	if err != nil {
		slog.Error(err.Error())
	}

	// Получение значений поиска
	requestUser, err := getSearchValue(r)
	if err != nil {
		slog.Error(err.Error())
	}

	// Ищем по запросу
	foundGroups := Search(requestUser, artists)

	// Выполняем наш шаблон
	err = tmp.Execute(w, foundGroups)
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
	searchValue = r.FormValue("search")

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
