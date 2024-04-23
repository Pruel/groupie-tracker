package controller

import (
	"fmt"
	"groupie-tracker/internal/entity"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// FilterController
func FilterController(w http.ResponseWriter, r *http.Request) {
	// read request body and save filter params in some structure
    if r.Method != http.MethodPost {
        slog.Debug("Invalid request")
        ErrorController(w, r)
    }

    // readValidateAndSaveFilterData вызываем функцию 

    // filteredArtists = Filter

    // filteredArtists > tmp > execute

	// http.Redirect(w, r, "/")
}

func readValidateAndSaveFilterData(w http.ResponseWriter, r *http.Request, flt *entity.Filters) *entity.Filters {
    // 1. read
    
    // 2. validate + save 

    // 3

    
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

func validateAndSaveData(flt *entity.Filters, fltData entity.Filters) *entity.Filters { 
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

    for _, num := range fltData.Members { // fltData.Members = []int{3, 7, 8}
        if num >= minMember && num <= maxMember { // bug
            flt.Members = append(flt.Members, num) 
        }
    }
    
   
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
    return flt
}

// Можно написать функцию валидации отдельно 