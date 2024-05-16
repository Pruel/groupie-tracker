package controller

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"groupie-tracker/internal/entity"
	"groupie-tracker/internal/webapi"
)

func ReadValidateAndSaveFilterData(r *http.Request, flt *entity.Filters) {
	filter := readRequest(r, flt)
	validateAndSaveData(flt, *filter)
}

func readRequest(r *http.Request, flt *entity.Filters) *entity.Filters {
	if err := r.ParseForm(); err != nil {
		slog.Error(err.Error())
	}

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
		Location = "default, defualt"
	}

	locs = append(locs, Location)
	flt.Locations = locs

	numMembers := make([]int, 0, 8)
	flt.Members = numMembers

	for i := 1; i < 9; i++ {
		memberKey := fmt.Sprintf("members%d", i)
		if member := r.FormValue(memberKey); member != "" {
			mNum, err := strconv.Atoi(member)
			if err != nil {
				slog.Error(err.Error())
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

	if fltData.FirstRelease >= minYear && fltData.FirstRelease <= time.Now().Year() {
		flt.FirstRelease = fltData.FirstRelease
	}

	memBuf := make([]int, 0, 8)
	for _, num := range fltData.Members {
		if num >= minMember && num <= maxMember {
			memBuf = append(memBuf, num)
		}
	}

	if len(memBuf) != 0 {
		flt.Members = memBuf
	}

	locs := make([]string, 0, 1)
	if len(fltData.Locations) != 0 && fltData.Locations[0] != "" {
		strSlice := strings.Split(fltData.Locations[0], ",")
		if len(strSlice) == 2 {
			locs = append(locs, fltData.Locations[0])
		}
	}

	flt.Locations = locs
}

func Filter(flt *entity.Filters, artists []entity.Artist) (filteredArt []entity.Artist, message string) {
	for _, group := range artists {
		matchloc := false
		matchcd := false
		matchmn := false
		matchalbum := false

		if group.CreationDate >= flt.CreationDate && group.CreationDate <= flt.CreationDate {
			matchcd = true
		}

		if convStrToTime(group.FirstAlbum) >= flt.FirstRelease && convStrToTime(group.FirstAlbum) <= flt.FirstRelease {
			matchalbum = true
		}

		for _, memNum := range flt.Members {
			if memNum == len(group.Members) {
				matchmn = true
				break
			}
		}

		bufLocs, _ := webapi.New().GetLocationsByURL(group.Locations)
		for _, loc := range bufLocs.Locations {
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
		message = "No group matches the requested parameter"
	}

	return filteredArt, message
}

func convStrToTime(strDate string) int {
	timeFormat := "02-01-2006"

	date, err := time.Parse(timeFormat, strDate)
	if err != nil {
		slog.Error(err.Error())
	}

	return date.Year()
}
