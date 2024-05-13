package controller

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"

	"groupie-tracker/internal/entity"
	"groupie-tracker/internal/filter"
	"groupie-tracker/internal/webapi"
)

// FilterController
func FilterController(w http.ResponseWriter, r *http.Request) {
	tmp := template.Must(template.ParseFiles(GetTmplFilepath("main.html")))

	// read request body and save filter params in some structure

	artists, err := webapi.New().GetAllArtists()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	fltData, err := filter.PrepareFilterData(artists)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	srchData := getAllUniqueSuggestions(artists)
	// call getAllUniqueSuggestions

	mdata := entity.MainData{
		Artists:     artists,
		FiltersData: *fltData,
		SearchData:  srchData,
	}

	if r.Method != http.MethodPost {
		slog.Error("Method not allowed!")
		ErrorController(w, r)
		return
	}

	// readValidateAndSaveFilterData вызываем функцию
	readValidateAndSaveFilterData(r, fltData)

	// filter
	filteredArtists, message := Filter(fltData, artists)
	mdata.Artists = filteredArtists
	mdata.Message = message

	if err := tmp.Execute(w, mdata); err != nil {
		slog.Error(err.Error())
		return
	}
}

func readValidateAndSaveFilterData(r *http.Request, flt *entity.Filters) {
	// 1. read
	filter := readRequest(r, flt) // TODO: fix bug = change pointer to some value type
	// 2. validate + save
	validateAndSaveData(flt, *filter)
	// 3

	// Type Pointer = &ValueOFRAMAddress && Ссылочный тип хранит в себе адрес ячейки памяти

	// var name string = "Daniil" // К примеру у нас есть строка которая весит 15mb
	// pName := &name // 0x1v2b345b55 Так-же для примера мы создаём  указатель на нашу строку которая весит 15mb
	// SecondName := name // И к примеру когда мы взяли и создали коппию то мы снова заняли ячейку памяти и у нас получился перерасход там где не нужно. // 30mb

	// Преимущество работы с ссылочным типом заключается в том что мы не перерасходуем нашу память. Мы можем из разных участков программы ссылаться на одну ячейку памяти

}

func readRequest(r *http.Request, flt *entity.Filters) *entity.Filters {
	// Connected with frontend, можно и без парсинга получать запрос от юзера, но базово можем использовать это
	if err := r.ParseForm(); err != nil {
		slog.Error(err.Error())
	}

	// Указываем что хотим получить, смотрим на атрибуты "for" and "name"
	FirstCreationDate := r.FormValue("CreationDate")

	if FirstCreationDate != "" {
		res, err := strconv.Atoi(FirstCreationDate)
		if err != nil {
			slog.Error(err.Error())
		}
		flt.CreationDate = res
	}

	FirstAlbum := r.FormValue("FirstRelease")

	if FirstAlbum != "" {
		date, err := strconv.Atoi(FirstAlbum)
		if err != nil {
			slog.Error(err.Error())
		}
		flt.FirstRelease = date
	}

	Location := r.FormValue("locations")
	locs := make([]string, 0, 1)

	if Location == "" {
		Location = "default, defualt" // "city, country"
	}

	locs = append(locs, Location)
	flt.Locations = locs

	// Создаём слайс так как нам нужно будет проверять каждый checkbox и это удобно будет сделать через слайс
	numMembers := make([]int, 0, 8) // len = 0
	flt.Members = numMembers

	for i := 1; i < 9; i++ { //
		memberKey := fmt.Sprintf("members%d", i) // Sprintf == Printf но не выводит в консоль а возвращает строку
		if member := r.FormValue(memberKey); member != "" {
			mNum, err := strconv.Atoi(member)
			if err != nil {
				slog.Error(err.Error()) // Наш кастомизированный логер который вернёт нам полный путь и строку откуда произошла ошибка
			}
			numMembers = append(numMembers, mNum)
		}
	}
	flt.Members = numMembers

	return flt
}

func validateAndSaveData(flt *entity.Filters, fltData entity.Filters) {
	minYear := 1962
	minMember := 1
	maxMember := 8

	if fltData.CreationDate >= minYear && fltData.CreationDate <= time.Now().Year() {
		flt.CreationDate = fltData.CreationDate
	}

	// II. firstAlbumDate >= minYear 1899 && firstAlbumDate <= time.Now().Year()
	if fltData.FirstRelease >= minYear && fltData.FirstRelease <= time.Now().Year() {
		flt.FirstRelease = fltData.FirstRelease
	}

	// III. Number of Members || numMembers >= 1 && numMembers <= 8
	memBuf := make([]int, 0, 8)
	for _, num := range fltData.Members { // fltData.Members = []int{3, 7, 8}
		if num >= minMember && num <= maxMember {
			memBuf = append(memBuf, num)
		}
	}

	if len(memBuf) != 0 {
		flt.Members = memBuf
	}

	// "city-country" O(n)
	// IV. Location || location != "" &&
	// location => city and country "los_angeles-usa" => sliceStr []string{"city", "country"} =  strings.Split(str, "-")
	// if len(sliceStr) == 2 { add }

	locs := make([]string, 0, 1)
	if len(fltData.Locations) != 0 && fltData.Locations[0] != "" {
		strSlice := strings.Split(fltData.Locations[0], ",")
		if len(strSlice) == 2 {
			locs = append(locs, fltData.Locations[0])
		}
	}

	flt.Locations = locs
}

// Filter
func Filter(flt *entity.Filters, artists []entity.Artist) (filteredArt []entity.Artist, message string) {
	// red case = early exit

	for _, group := range artists {
		matchloc := false   //location
		matchcd := false    // creation date
		matchmn := false    // members number
		matchalbum := false // first album

		// I. creation date
		if group.CreationDate >= flt.CreationDate && group.CreationDate <= flt.CreationDate {
			matchcd = true
		}

		// II. album publish
		if convStrToTime(group.FirstAlbum) >= flt.FirstRelease && convStrToTime(group.FirstAlbum) <= flt.FirstRelease {
			matchalbum = true
		}

		// III. members number
		for _, memNum := range flt.Members {
			if memNum == len(group.Members) {
				matchmn = true
				break
			}
		}

		// IV. location
		bufLocs, _ := webapi.New().GetLocationsByURL(group.Locations) // artists -> artist.Location = url -> request -> []stirng{"New York, USA", }
		for _, loc := range bufLocs.Locations {                       // Members + Locations <- defualt
			loc = webapi.ParseAndFormatLocations(loc)
			if flt.Locations[0] == loc {
				matchloc = true
				break
			}
		}

		if matchcd || matchalbum || matchmn || matchloc {
			filteredArt = append(filteredArt, group)
		}

	}

	if len(filteredArt) == 0 {
		// logic
		message = "No group matches the requested parameter"
	}

	return filteredArt, message
}

// convStrToTime, return time.Time.Year
func convStrToTime(strDate string) int {
	timeFormat := "02-01-2006" // "02 October 2005 10:10:10" time.UnixDate, time.ANSIC

	date, err := time.Parse(timeFormat, strDate) // group.FirstAlbum // 15-02-2024
	if err != nil {
		slog.Error(err.Error())
	}

	return date.Year() // date time.Time
}
