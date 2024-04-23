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
	// read request body and save filter params in some structure
	if r.Method != http.MethodPost {
		slog.Debug("Invalid request")
		ErrorController(w, r)
	}

	tmp := template.Must(template.ParseFiles(GetTmplFilepath("index.html")))

	artists, err := webapi.New().GetAllArtists()
	if err != nil {
		slog.Error(err.Error())
	}

	fltData, err := filter.PrepareFilterData(artists)
	if err != nil {
		slog.Error(err.Error())
	}

	mdate := entity.MainData{
		Artists:     artists,
		FiltersData: *fltData,
	}

	fmt.Printf("\n\n1 . after recieve data from front-end, filtersData: %+v\n\n", fltData)
	// readValidateAndSaveFilterData вызываем функцию
	readValidateAndSaveFilterData(w, r, fltData)

	// filter

	fmt.Printf("\n\n2. after recieve data from front-end, filtersData: %+v\n\n", fltData)

	// filteredArtists = Filter

	// filteredArtists > tmp > execute

	// http.Redirect(w, r, "/")

	if err := tmp.Execute(w, mdate); err != nil {
		slog.Error(err.Error())
	}
}

func readValidateAndSaveFilterData(w http.ResponseWriter, r *http.Request, flt *entity.Filters) {
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
	res, err := strconv.Atoi(LastCreationDate)
	if err != nil {
		slog.Error(err.Error())
	}
	flt.LastCreationDate = res

	flt.HighestFirstAlbum = r.FormValue("firstAlbumDate")
	Locations := r.FormValue("locations")
	flt.Locations = append(flt.Locations, Locations)

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
	timeFormat := "21-01-1995"

	//time.Parse
	// highestFirstAlbum = hAlbum
	hAlbum, err := time.Parse(timeFormat, fltData.HighestFirstAlbum)
	if err != nil {
		slog.Error(err.Error())
	}

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
		if num >= minMember && num <= maxMember { // bug
			memBuf = append(memBuf, num) // Buf !
		}
	}
    flt.Members = memBuf // 

	// "city-country" O(n)
	// IV. Location || location != "" &&
	// location => city and country "los_angeles-usa" => sliceStr []string{"city", "country"} =  strings.Split(str, "-")
	// if len(sliceStr) == 2 { add }

	if fltData.Locations[0] != "" {
		strSlice := strings.Split(fltData.Locations[0], "-")
		if len(strSlice) == 2 {
			flt.Locations = append(flt.Locations, fltData.Locations[0])
		}
	}
	// 3. save to filters
	// ftl.FirstCreationDate = creationDate

	// 4. return
	// return flt
}

// Можно написать функцию валидации отдельно
