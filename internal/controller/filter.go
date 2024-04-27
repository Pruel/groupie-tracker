package controller

import (
	"fmt"
	"groupie-tracker/internal/entity"
	"groupie-tracker/internal/filter"
	"groupie-tracker/internal/webapi"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"text/template"
	"time"
)

// FilterController
func FilterController(w http.ResponseWriter, r *http.Request) {
	tmp := template.Must(template.ParseFiles(GetTmplFilepath("index.html")))

	// read request body and save filter params in some structure

	artists, err := webapi.New().GetAllArtists()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	fArtists := artists

	fltData, err := filter.PrepareFilterData(artists)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	mdata := entity.MainData{
		Artists:     artists,
		FiltersData: *fltData,
	}

	if r.Method == http.MethodPost {
		fmt.Println("POST METHOD")
		// readValidateAndSaveFilterData вызываем функцию
		readValidateAndSaveFilterData(r, fltData)

		// filter
		filteredArtists := Filter(fltData, artists)
		mdata.Artists = filteredArtists

		fmt.Printf("\n\n2. after recieve data from front-end, filtersData: %+v\n\n", fltData)

		err := tmp.Execute(w, mdata)
		if err != nil {
			slog.Debug("request method is get")
			return
		}
		return
	}

	// filteredArtists = Filter

	// filteredArtists > tmp > execute

	// http.Redirect(w, r, "/")
	mdata.Artists = fArtists

	fmt.Println("GET METHOD")
	err = tmp.Execute(w, mdata)
	if err != nil {
		slog.Error(err.Error())
		return
	}
}

func readValidateAndSaveFilterData(r *http.Request, flt *entity.Filters) {
	// 1. read
	filter := readRequest(r)
	// 2. validate + save
	validateAndSaveData(flt, filter)
	// 3

	// Type Pointer = &ValueOFRAMAddress && Ссылочный тип хранит в себе адрес ячейки памяти

	// var name string = "Daniil" // К примеру у нас есть строка которая весит 15mb
	// pName := &name // 0x1v2b345b55 Так-же для примера мы создаём  указатель на нашу строку которая весит 15mb
	// SecondName := name // И к примеру когда мы взяли и создали коппию то мы снова заняли ячейку памяти и у нас получился перерасход там где не нужно. // 30mb

	// Преимущество работы с ссылочным типом заключается в том что мы не перерасходуем нашу память. Мы можем из разных участков программы ссылаться на одну ячейку памяти

}

func readRequest(r *http.Request) (flt entity.Filters) {
	r.ParseForm() // Connected with frontend, можно и без парсинга получать запрос от юзера, но базово можем использовать это

	// Указываем что хотим получить, смотрим на атрибуты "for" and "name"
	LastCreationDate := r.FormValue("creationDate")

	if LastCreationDate != "" {
		res, err := strconv.Atoi(LastCreationDate)
		if err != nil {
			slog.Error(err.Error()) // fix
		}
		flt.LastCreationDate = res
	}

	HighestFirstAlbum := r.FormValue("firstAlbumDate")

	if HighestFirstAlbum != "" {
		flt.HighestFirstAlbum = HighestFirstAlbum
	}

	fmt.Println("Current date:", flt.HighestFirstAlbum) // 1999-05-05

	Locations := r.FormValue("locations")
	locs := make([]string, 0, 1)

	locs = append(locs, Locations)
	flt.Locations = locs

	// Создаём слайс так как нам нужно будет проверять каждый checkbox и это удобно будет сделать через слайс
	numMembers := make([]int, 0, 8)

	for i := 1; i < 8; i++ {
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
	minYear := 1899
	minMember := 1
	maxMember := 8
	timeFormat := "2006-01-02"

	//time.Parse
	// highestFirstAlbum = hAlbum
	fmt.Println("Date from frontend:", fltData.HighestFirstAlbum)
	hAlbum, err := time.Parse(timeFormat, fltData.HighestFirstAlbum)
	if err != nil {
		slog.Error(err.Error()) // bug
	}

	fmt.Println("TIME into validation func: ", hAlbum)

	if fltData.LastCreationDate >= minYear && fltData.LastCreationDate <= time.Now().Year() {
		flt.LastCreationDate = fltData.LastCreationDate //

	}

	// II. firstAlbumDate >= minYear 1899 && firstAlbumDate <= time.Now().Year()
	if hAlbum.Year() >= minYear && hAlbum.Year() <= time.Now().Year() {
		flt.HighestFirstAlbum = fltData.HighestFirstAlbum
	}

	// III. Number of Members || numMembers >= 1 && numMembers <= 8
	memBuf := make([]int, 0, 8)
	for _, num := range fltData.Members { // fltData.Members = []int{3, 7, 8}
		if num >= minMember && num <= maxMember {
			memBuf = append(memBuf, num)
		}
	}
	flt.Members = memBuf

	// "city-country" O(n)
	// IV. Location || location != "" &&
	// location => city and country "los_angeles-usa" => sliceStr []string{"city", "country"} =  strings.Split(str, "-")
	// if len(sliceStr) == 2 { add }

	locs := make([]string, 0, 1)
	if fltData.Locations[0] != "" {
		strSlice := strings.Split(fltData.Locations[0], ",")
		if len(strSlice) == 2 {
			locs = append(locs, fltData.Locations[0])
		}
	}
	flt.Locations = locs
}

// Filter ..

// 1. match boolean flog = true

// I. creation date

// II. album publish date

// III.  members number

// IV. locations

// if match do append in a new filteredArtists

// return this new slice
func Filter(flt *entity.Filters, artists []entity.Artist) (filteredArtists []entity.Artist) {
	match := false

	// I. creation date
	for _, group := range artists {
		if group.CreationDate >= flt.FirstCreationDate && group.CreationDate <= flt.LastCreationDate {
			match = true
		}

		// II. album publish

		if convStrToTime(group.FirstAlbum, "GFA") >= convStrToTime(flt.LowestFirstAlbum, "LFA") && convStrToTime(group.FirstAlbum, "GFA") <= convStrToTime(flt.HighestFirstAlbum, "HFA") { // bug
			match = true
		}

		// III.  members number
		for _, num := range flt.Members {
			if len(group.Members) == num {
				match = true
				break
			}
			match = false
		}

		// IV. locations group []string{}

		gLocs, err := webapi.New().GetLocationsByURL(group.Locations)
		if err != nil {
			slog.Error(err.Error())
			continue
		}

		for _, loc := range gLocs.Locations { // []string{"Some location", "Location"}
			if strings.Contains(flt.Locations[0], loc) {
				match = true
				break // operator || continue, break,
			}
			match = false

		}

		// if match do append in a new filteredArtists
		if match {
			filteredArtists = append(filteredArtists, group)
		}

	}

	return filteredArtists
}

// convStrToTime, return time.Time.Year
func convStrToTime(strDate string, flag string) int {
	if flag == "" {
		return 0
	}
	// timeFormat := "2006-01-02" // "02 October 2005 10:10:10" time.UnixDate, time.ANSIC

	var date time.Time
	var err error
	// flt.LowestFirstAlbum // = 01-05-2024
	// flt.HighestFirstAlbum // = 2024-05-01
	// group.FirstAlbum

	switch flag { // GFA  = (group first album), flt.LowestFA = LFA, flt.HighestFA = HFA
	case "GFA":
		date, err = time.Parse(time.DateOnly, strDate) // group.FirstAlbum // 15-02-2024
	case "LFA":
		date, err = time.Parse(time.DateOnly, strDate) // group.FirstAlbum // 15-02-2024
	case "HFA":
		timeFormat := "2006-01-02"
		date, err = time.Parse(timeFormat, strDate) // group.FirstAlbum // 15-02-2024
	}

	if err != nil {
		slog.Error(err.Error())
	}

	return date.Year() // date time.Time
}
